// Package idp_e2e 是 RAM 身份提供商（IdP）相关能力的端到端测试套件。
//
// 覆盖范围：
//   - 6 张 IdP 库表对应的控制器：身份提供商、认证源、回调白名单、身份绑定、
//     凭证、元数据缓存（create/update/detail/query/selectPublic/existName/
//     enable/disable/state/delete/recovery/physicalDeletion）
//   - IdP OAuth2 / OpenID 授权码登录（含首次自动建号、二次绑定登录、失败分支）
//   - CAS 1.0 / 2.0 / 3.0 票据签发与校验
//   - SAML2 登录 URL 生成与回调错误分支
//   - LDAP 连接测试 / 搜索 / 同步 / 登录
//
// 测试直接连接 data/config/app.toml 中配置的 PostgreSQL（默认 test-make-01），
// 会创建并清理带唯一前缀的测试数据；不依赖正在运行的服务进程。
//
// 运行方式（见 docs/test/idp-e2e-test.md）：
//
//	GOROOT=/Users/mac/develop/go PATH=$PATH:$GOROOT/bin \
//	GOCACHE=/tmp/dsh-go-build MOCKEY_CHECK_GCFLAGS=false \
//	IDP_TENANT_NO=1000 go test ./test/idp -run TestIDP -v -count=1
package idp_e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/handler/manage/idp"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/configPg/pg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/cacheService/cacheRam"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/consts/constsRam/typeDomainPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/state/clientPg"
	routerPg "github.com/hongmengzhu/xianfu-blog-go/pkg/core/routerPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/entityIdp"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/entityRam"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/repositoryIdp"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/repositoryRam"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/middleware/components/authTokenPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/dbManagerPg"
	"go-spring.org/spring/gs"
	"gorm.io/gorm"
)

const (
	testTenantNo = "1000"
	baseURL      = "/xianfu/manage/ram"
	authIdpURL   = "/xianfu/auth/manage/idp"
	manageCasURL = "/xianfu/auth/manage/cas"
)

// repoRoot 返回仓库根目录（本文件位于 <root>/test/idp）。
func repoRoot() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Dir(filepath.Dir(filepath.Dir(file)))
}

func TestMain(m *testing.M) {
	// go test 的工作目录是包目录，统一切到仓库根目录，
	// 使 app.toml 中的相对路径（data/logs、data/config）与生产一致。
	if err := os.Chdir(repoRoot()); err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}

func init() {
	gs.Provide(new(dbManagerPg.Manager))
}

// idpStack go-spring 注入的控制器、仓库与缓存对象。
type idpStack struct {
	Provider   *idp.IdentityProviderController       `autowire:"?"`
	Source     *idp.IdentitySourceController         `autowire:"?"`
	Callback   *idp.IdentitySourceCallbackController `autowire:"?"`
	Binding    *idp.IdpBindingController             `autowire:"?"`
	Credential *idp.IdpCredentialController          `autowire:"?"`
	Metadata   *idp.IdpMetadataCacheController       `autowire:"?"`
	Ldap       *idp.LdapController                   `autowire:"?"`
	Cas        *idp.CasController                    `autowire:"?"`
	AuthIdp    *idp.AuthLoginIdpController           `autowire:"?"`

	Manager *dbManagerPg.Manager                          `autowire:"?"`
	Session *cacheRam.CacheSessionPubPrive                `autowire:"?"`
	SourceR *repositoryIdp.IdpIdentitySourceRepository    `autowire:"?"`
	ProvR   *repositoryIdp.IdpIdentityProviderRepository  `autowire:"?"`
	BindR   *repositoryIdp.IdpIdpBindingRepository        `autowire:"?"`
	CredR   *repositoryIdp.IdpIdpCredentialRepository     `autowire:"?"`
	MetaR   *repositoryIdp.IdpIdpMetadataCacheRepository  `autowire:"?"`
	AccR    *repositoryRam.RamAccountRepository           `autowire:"?"`
	LogR    *repositoryRam.RamAccountSessionLogRepository `autowire:"?"`
}

// env 一次测试运行的上下文。
type env struct {
	t        *testing.T
	ts       *idpStack
	engine   *gin.Engine
	db       *gorm.DB
	token    string
	runID    string
	soft     bool
	cleanups []func()
}

// TestIDP 是唯一顶层用例：整个进程只启动一次 go-spring 容器，
// 所有场景通过 t.Run 子用例共享同一套真实路由与数据库连接。
func TestIDP(t *testing.T) {
	root := repoRoot()
	soft := os.Getenv("IDP_TEST_SOFT_DELETE") == "true"
	tenant := os.Getenv("IDP_TENANT_NO")
	if tenant == "" {
		tenant = testTenantNo
	}

	configDir := filepath.Join(root, "data", "config")
	if soft {
		// app.toml 中的配置会覆盖 app.Property 设置，因此逻辑删除模式需要
		// 基于 app.toml 生成一份 delete=false 的临时配置目录。
		src, err := os.ReadFile(filepath.Join(configDir, "app.toml"))
		if err != nil {
			t.Fatalf("读取 app.toml 失败: %v", err)
		}
		replaced := strings.Replace(string(src), "delete = true", "delete = false", 1)
		if replaced == string(src) {
			t.Fatalf("app.toml 中未找到 pg.data.delete 配置")
		}
		tmpDir, err := os.MkdirTemp("/tmp", "idp-soft-config-")
		if err != nil {
			t.Fatalf("创建临时配置目录失败: %v", err)
		}
		t.Cleanup(func() { _ = os.RemoveAll(tmpDir) })
		if err := os.WriteFile(filepath.Join(tmpDir, "app.toml"), []byte(replaced), 0o644); err != nil {
			t.Fatalf("写入临时 app.toml 失败: %v", err)
		}
		configDir = tmpDir
	}

	gs.Web(false).Configure(func(app gs.App) {
		app.Property("spring.app.config.dir", configDir)
	}).RunTest(t, func(ts *idpStack) {
		if ts.Provider == nil || ts.Source == nil || ts.Callback == nil || ts.Binding == nil ||
			ts.Credential == nil || ts.Metadata == nil || ts.Ldap == nil || ts.Cas == nil ||
			ts.AuthIdp == nil || ts.Manager == nil || ts.Session == nil {
			t.Fatalf("容器装配失败: %+v", ts)
		}
		db := ts.Manager.Db()
		ensureSchema(t, db)

		e := &env{
			t:     t,
			ts:    ts,
			db:    db,
			runID: time.Now().Format("0102-150405"),
			soft:  soft,
		}
		// 必须在 RunTest 回调内清理：回调结束后容器关闭，数据库连接会失效。
		defer e.cleanupAll()

		e.engine = gin.New()
		gin.SetMode(gin.TestMode)
		for _, r := range []routerPg.RouteRegistrar{
			ts.Provider, ts.Source, ts.Callback, ts.Binding,
			ts.Credential, ts.Metadata, ts.Ldap, ts.Cas, ts.AuthIdp,
		} {
			r.RegisterRoutes(e.engine)
		}
		e.mintToken(tenant)
		t.Logf("测试运行 ID=%s tenant=%s softDelete=%v dbDriver=%s",
			e.runID, tenant, soft, db.Dialector.Name())

		t.Run("routes-registered", func(t *testing.T) { runRouteCheck(t, e) })
		t.Run("schema", func(t *testing.T) { runSchemaCheck(t, e) })

		t.Run("identity-provider", func(t *testing.T) { runProviderCRUD(t, e) })
		t.Run("identity-source", func(t *testing.T) { runSourceCRUD(t, e) })
		t.Run("identity-source-callback", func(t *testing.T) { runCallbackCRUD(t, e) })
		t.Run("idp-binding", func(t *testing.T) { runBindingCRUD(t, e) })
		t.Run("idp-credential", func(t *testing.T) { runCredentialCRUD(t, e) })
		t.Run("idp-metadata-cache", func(t *testing.T) { runMetadataCRUD(t, e) })

		t.Run("idp-login", func(t *testing.T) { runIdpLogin(t, e) })
		t.Run("cas", func(t *testing.T) { runCAS(t, e) })
		t.Run("saml", func(t *testing.T) { runSAML(t, e) })
		t.Run("ldap", func(t *testing.T) { runLDAP(t, e) })

		t.Run("live-server-smoke", func(t *testing.T) { runLiveSmoke(t, e) })
	})
}

// ensureSchema 兜底建表。正常情况下 app 启动时 AInitTable 已建表；
// 测试容器不执行启动 runner，这里显式 AutoMigrate，保证用例可独立运行。
func ensureSchema(t *testing.T, db *gorm.DB) {
	t.Helper()
	entities := []any{
		&entityRam.RamAccountAuthMfaEntity{},
		&entityRam.RamAccountAuthWebauthnEntity{},
		&entityRam.RamAccountSessionLogEntity{},
		&entityIdp.IdpIdentityProviderEntity{},
		&entityIdp.IdpIdentitySourceEntity{},
		&entityIdp.IdpIdentitySourceCallbackEntity{},
		&entityIdp.IdpIdpBindingEntity{},
		&entityIdp.IdpIdpCredentialEntity{},
		&entityIdp.IdpIdpMetadataCacheEntity{},
	}
	if err := db.AutoMigrate(entities...); err != nil {
		t.Fatalf("AutoMigrate 失败: %v", err)
	}
}

// mintToken 使用数据库中的 AccessToken 密钥对签发管理端 PASETO，
// 与中间件 authPg.GroupManageMiddleware 的校验路径完全一致。
func (e *env) mintToken(tenant string) {
	key := e.ts.Session.PaseKeyAccessToken(context.Background(), typeDomainPg.Manage, clientPg.Browser, tenant, nil)
	if key.Private == "" || key.Public == "" {
		e.t.Fatalf("租户 %s 未取到 AccessToken 密钥对", tenant)
	}
	param := authTokenPg.NewJwtPg()
	param.UniqueId = "idp-e2e-" + e.runID
	param.LoginNo = tenant
	param.No = tenant
	param.TenantNo = tenant
	param.Type = string(typeDomainPg.Manage)
	param.Result = authTokenPg.Result{PrivateKey: key.Private, PublicKey: key.Public}
	res := authTokenPg.MakePaseToken(*param, pg.JwtConfig{Expire: 240, Issuer: "foxwho.com", Audience: "pcWeb"})
	if res.ErrorIs() {
		e.t.Fatalf("签发测试 token 失败: %s", res.Message)
	}
	e.token = res.Data.Token
}

// ---- HTTP 辅助 ----

type apiResp struct {
	Code        string          `json:"code"`
	Message     string          `json:"message"`
	MessageData json.RawMessage `json:"messageData"`
	Data        json.RawMessage `json:"data"`
	Status      int             `json:"-"`
	Header      http.Header     `json:"-"`
	Body        string          `json:"-"`
}

func parseResp(w *httptest.ResponseRecorder) *apiResp {
	r := &apiResp{Status: w.Code, Header: w.Header(), Body: w.Body.String()}
	_ = json.Unmarshal(w.Body.Bytes(), r)
	return r
}

func (e *env) do(method, path, body string, auth bool, headers map[string]string) *apiResp {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if auth {
		req.Header.Set("Authorization", e.token)
	}
	req.Header.Set("User-Agent", "dsh-idp-e2e/1.0")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	w := httptest.NewRecorder()
	e.engine.ServeHTTP(w, req)
	e.t.Helper()
	return parseResp(w)
}

func (e *env) post(path, body string) *apiResp { return e.do(http.MethodPost, path, body, true, nil) }
func (e *env) postNoAuth(path, body string) *apiResp {
	return e.do(http.MethodPost, path, body, false, nil)
}
func (e *env) get(path string) *apiResp     { return e.do(http.MethodGet, path, "", false, nil) }
func (e *env) getAuth(path string) *apiResp { return e.do(http.MethodGet, path, "", true, nil) }

// js 把 map/slice 序列化为 JSON 请求体。
func js(v any) string {
	b, _ := json.Marshal(v)
	return string(b)
}
func (e *env) postForm(path string, form map[string]string) *apiResp {
	vals := url.Values{}
	for k, v := range form {
		vals.Set(k, v)
	}
	req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(vals.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	e.engine.ServeHTTP(w, req)
	return parseResp(w)
}

func (e *env) ok(path, body string) *apiResp {
	e.t.Helper()
	r := e.post(path, body)
	if r.Code != "200" {
		e.t.Fatalf("POST %s 期望成功，实际 code=%s message=%s body=%s", path, r.Code, r.Message, r.Body)
	}
	return r
}

func (e *env) expectErr(path, body, contains string) *apiResp {
	e.t.Helper()
	r := e.post(path, body)
	if r.Code == "200" {
		e.t.Fatalf("POST %s 期望失败(%q)，实际成功 body=%s", path, contains, r.Body)
	}
	if contains != "" && !strings.Contains(r.Message, contains) {
		e.t.Fatalf("POST %s 期望 message 包含 %q，实际 %q (code=%s)", path, contains, r.Message, r.Code)
	}
	return r
}

func decodeData[T any](t *testing.T, r *apiResp) T {
	t.Helper()
	var v T
	if err := json.Unmarshal(r.Data, &v); err != nil {
		t.Fatalf("解析 data 失败: %v data=%s", err, string(r.Data))
	}
	return v
}

func (e *env) createID(path, body, table string) string {
	e.t.Helper()
	r := e.ok(path+"/create", body)
	id := decodeData[string](e.t, r)
	if id == "" {
		e.t.Fatalf("create 未返回 id: %s", r.Body)
	}
	e.track(table, id)
	return id
}

func (e *env) track(table, id string) {
	n, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		e.t.Logf("track 忽略非数字 id: %s", id)
		return
	}
	e.cleanups = append(e.cleanups, func() {
		if err := e.db.Exec("DELETE FROM "+table+" WHERE id=?", n).Error; err != nil {
			fmt.Printf("[idp-e2e] 清理 %s id=%d 失败: %v\n", table, n, err)
		}
	})
}

func (e *env) trackAccount(no string) {
	e.cleanups = append(e.cleanups, func() {
		e.db.Exec("DELETE FROM ram_account WHERE no=?", no)
		e.db.Exec("DELETE FROM ram_account_session_log WHERE ano=?", no)
	})
}

func (e *env) cleanupAll() {
	for i := len(e.cleanups) - 1; i >= 0; i-- {
		e.cleanups[i]()
	}
}

func (e *env) unique(prefix string) string {
	return fmt.Sprintf("E2E-%s-%s-%d", prefix, e.runID, time.Now().UnixNano()%100000)
}

// ---- 通用断言/用例 ----

func runRouteCheck(t *testing.T, e *env) {
	want := map[string]string{
		baseURL + "/identity-provider/create":           http.MethodPost,
		baseURL + "/identity-provider/update":           http.MethodPost,
		baseURL + "/identity-provider/detail/:id":       http.MethodGet,
		baseURL + "/identity-provider/enable":           http.MethodPost,
		baseURL + "/identity-provider/disable":          http.MethodPost,
		baseURL + "/identity-provider/state":            http.MethodPost,
		baseURL + "/identity-provider/delete":           http.MethodPost,
		baseURL + "/identity-provider/recovery":         http.MethodPost,
		baseURL + "/identity-provider/physicalDeletion": http.MethodPost,
		baseURL + "/identity-provider/query":            http.MethodPost,
		baseURL + "/identity-provider/selectPublic":     http.MethodPost,
		baseURL + "/identity-provider/existName":        http.MethodPost,
		baseURL + "/identity-source/create":             http.MethodPost,
		baseURL + "/identity-source/query":              http.MethodPost,
		baseURL + "/identity-source/selectPublic":       http.MethodPost,
		baseURL + "/identity-source-callback/create":    http.MethodPost,
		baseURL + "/identity-source-callback/query":     http.MethodPost,
		baseURL + "/idp-binding/create":                 http.MethodPost,
		baseURL + "/idp-binding/update":                 http.MethodPost,
		baseURL + "/idp-binding/query":                  http.MethodPost,
		baseURL + "/idp-binding/selectPublic":           http.MethodPost,
		baseURL + "/idp-credential/create":              http.MethodPost,
		baseURL + "/idp-credential/query":               http.MethodPost,
		baseURL + "/idp-metadata-cache/create":          http.MethodPost,
		baseURL + "/idp-metadata-cache/query":           http.MethodPost,
		baseURL + "/ldap/test":                          http.MethodPost,
		baseURL + "/ldap/search":                        http.MethodPost,
		baseURL + "/ldap/sync":                          http.MethodPost,
		baseURL + "/ldap/login":                         http.MethodPost,
		authIdpURL + "/cas/login":                       http.MethodPost,
		authIdpURL + "/login":                           http.MethodPost,
		authIdpURL + "/refresh":                         http.MethodPost,
		authIdpURL + "/saml/login/:sourceNo":            http.MethodGet,
		authIdpURL + "/saml/callback":                   http.MethodPost,
		manageCasURL + "/login":                         http.MethodPost,
		manageCasURL + "/validate":                      http.MethodGet,
		manageCasURL + "/serviceValidate":               http.MethodGet,
		manageCasURL + "/p3/serviceValidate":            http.MethodGet,
	}
	found := map[string]string{}
	for _, r := range e.engine.Routes() {
		found[r.Method+" "+r.Path] = r.Handler
	}
	for path, method := range want {
		if _, ok := found[method+" "+path]; !ok {
			t.Errorf("缺少路由 %s %s", method, path)
		}
	}
	t.Logf("已注册路由 %d 条，关键路由断言 %d 条", len(e.engine.Routes()), len(want))
}

func runSchemaCheck(t *testing.T, e *env) {
	type row struct {
		TableName string `gorm:"column:tablename"`
	}
	var rows []row
	if err := e.db.Raw("select tablename from pg_tables where schemaname='public' and tablename like 'ram_%' order by 1").Scan(&rows).Error; err != nil {
		t.Fatalf("查询表失败: %v", err)
	}
	tables := make([]string, 0, len(rows))
	for _, r := range rows {
		tables = append(tables, r.TableName)
	}
	for _, want := range []string{
		"ram_identity_provider", "ram_identity_source", "ram_identity_source_callback",
		"ram_idp_binding", "ram_idp_credential", "ram_idp_metadata_cache",
		"ram_account_session_log", "ram_account_auth_mfa", "ram_account_auth_webauthn",
	} {
		has := false
		for _, name := range tables {
			if name == want {
				has = true
				break
			}
		}
		if !has {
			t.Errorf("缺少表 %s（当前 ram_ 表: %v）", want, tables)
		}
	}
	t.Logf("ram_ 相关表 %d 张", len(tables))
}

// getNoByID 从仓库读取实体 No（Vo 未暴露 No，登录流程需要）。
func (e *env) providerNo(id string) string {
	n, _ := strconv.ParseInt(id, 10, 64)
	info, ok := e.ts.ProvR.FindById(context.Background(), n)
	if !ok {
		e.t.Fatalf("provider id=%s 不存在", id)
	}
	return info.No
}

func (e *env) sourceNo(id string) string {
	n, _ := strconv.ParseInt(id, 10, 64)
	info, ok := e.ts.SourceR.FindById(context.Background(), n)
	if !ok {
		e.t.Fatalf("source id=%s 不存在", id)
	}
	return info.No
}

// runLiveSmoke 可选的实机冒烟：设置 IDP_LIVE_BASE（如 http://127.0.0.1:9985）后，
// 用签发的 token 请求真实进程的查询接口，验证 echo/gin 入口与中间件链路。
func runLiveSmoke(t *testing.T, e *env) {
	base := os.Getenv("IDP_LIVE_BASE")
	if base == "" {
		t.Skip("未设置 IDP_LIVE_BASE，跳过实机冒烟")
	}
	req, err := http.NewRequest(http.MethodPost, strings.TrimRight(base, "/")+baseURL+"/identity-provider/query", strings.NewReader(`{}`))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", e.token)
	resp, err := (&http.Client{Timeout: 5 * time.Second}).Do(req)
	if err != nil {
		t.Fatalf("实机请求失败: %v", err)
	}
	defer resp.Body.Close()
	var env apiResp
	if err := json.NewDecoder(resp.Body).Decode(&env); err != nil {
		t.Fatalf("实机响应解析失败: %v", err)
	}
	if env.Code != "200" {
		t.Fatalf("实机查询失败: code=%s message=%s", env.Code, env.Message)
	}
	t.Logf("实机冒烟通过: %s", base)
}
