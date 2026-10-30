package controller

// 书签与书签分类接口的单元测试。
//
// 测试通过 go-spring 容器装配控制器、服务与仓库，把数据源指向临时 sqlite 文件库，
// 再用 httptest 驱动真实路由与鉴权中间件，覆盖 9 个接口（/api/bookmark 6 个、
// /api/bookmark-category 3 个）的成功与错误分支。
//
// 运行：go test ./app/domain/web/api/controller/ -v

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/web/api/service"
	cacheDiplPg "github.com/hongmengzhu/xianfu-blog-go/app/middleware/components/cachePg/cacheDiplPg"
	"github.com/hongmengzhu/xianfu-blog-go/app/middleware/components/cachePg/cacheDomainPg"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/tc/dtoTenantDomain"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityBlog"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/dbManagerPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/mapPg"
	"go-spring.org/spring/gs"
	"gorm.io/gorm"
)

const (
	testDiplKey    = "test-dipl-key"
	testDiplSecret = "test-dipl-secret"
	testTenantNo   = "1000"
	testAno        = "u-001"
	authHeaderName = "Authorization"

	pathBookmarkGetAll = "/api/bookmark/getAll"
	pathBookmarkGetMy  = "/api/bookmark/getMy"
	pathBookmarkGetOrg = "/api/bookmark/getOrg"
	pathBookmarkCreate = "/api/bookmark/create"
	pathBookmarkUpdate = "/api/bookmark/update"
	pathBookmarkDelete = "/api/bookmark/delete"

	pathCategoryGetAll       = "/api/bookmark-category/getAll"
	pathCategoryCreateUpdate = "/api/bookmark-category/createUpdate"
	pathCategoryDelete       = "/api/bookmark-category/delete"

	msgAuthRequired = "Authorization 参数不能为空"
	msgOK           = "操作成功"
)

// apiTestStack go-spring 注入的测试对象。
type apiTestStack struct {
	Bookmark *BookmarkApiController           `autowire:"?"`
	Category *BookmarkCategoryApiController   `autowire:"?"`
	Manager  *dbManagerPg.Manager             `autowire:"?"`
	Domain   *cacheDomainPg.TenantDomainCache `autowire:"?"`
}

// rsEnvelope 响应外层结构，data 保留原始 JSON 以便按各接口类型解码。
type rsEnvelope struct {
	Code        string          `json:"code"`
	Message     string          `json:"message"`
	MessageData json.RawMessage `json:"messageData"`
	Data        json.RawMessage `json:"data"`
}

// 以下为接口返回体的测试侧镜像，用于按契约校验字段名与类型。

type bookmarkVo struct {
	ID           string `json:"id"`
	No           string `json:"no"`
	Name         string `json:"name"`
	State        int    `json:"state"`
	CategoryNo   string `json:"categoryNo"`
	CategoryName string `json:"categoryName"`
	TeamNo       string `json:"teamNo"`
	UrlSource    string `json:"urlSource"`
}

type categoryVo struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	NameFl   string `json:"nameFl"`
	NameFull string `json:"nameFull"`
	No       string `json:"no"`
	Code     string `json:"code"`
	ParentNo string `json:"parentNo"`
}

type bookmarkVoAll struct {
	My           []bookmarkVo `json:"my"`
	MyCategory   []categoryVo `json:"MyCategory"`
	Team         []bookmarkVo `json:"team"`
	TeamCategory []categoryVo `json:"teamCategory"`
	Org          []bookmarkVo `json:"org"`
	OrgCategory  []categoryVo `json:"orgCategory"`
}

type categoryVoAll struct {
	My   []categoryVo `json:"my"`
	Team []categoryVo `json:"team"`
	Org  []categoryVo `json:"org"`
}

func init() {
	// main.go 中手工 Provide 的 Bean，测试入口需自行注册；
	// BookmarkService 全项目仅靠 autowire 注入，没有 init 注册。
	gs.Provide(new(dbManagerPg.Manager))
	gs.Provide(new(service.BookmarkService))
}

// runAPIStack 启动独立的测试应用：临时 sqlite 库 + 真实路由 + 鉴权中间件。
func runAPIStack(t *testing.T, fn func(t *testing.T, ts *apiTestStack, engine *gin.Engine, db *gorm.DB)) {
	t.Helper()
	dbFile := filepath.Join(t.TempDir(), "api-test.db")
	gs.Web(false).Configure(func(app gs.App) {
		app.Property("database.enabled", "true")
		app.Property("database.driver", "sqlite")
		app.Property("spring.gorm.sqlite.instances.primary.file", dbFile)
		// 用内存 journal，避免临时库产生 -wal/-shm 文件。
		app.Property("spring.gorm.sqlite.instances.primary.journal-mode", "memory")
		// 与 app.toml 中 [pg.data] delete = true 保持一致：删除为物理删除。
		app.Property("pg.data.delete", "true")
	}).RunTest(t, func(ts *apiTestStack) {
		if ts.Bookmark == nil || ts.Category == nil || ts.Manager == nil || ts.Domain == nil {
			t.Fatalf("容器装配失败: bookmark=%v category=%v manager=%v domain=%v",
				ts.Bookmark, ts.Category, ts.Manager, ts.Domain)
		}
		db := ts.Manager.Db()
		if err := db.AutoMigrate(&entityBlog.BlogBookmarkEntity{}, &entityBlog.BlogBookmarkCategoryEntity{}); err != nil {
			t.Fatalf("建表失败: %v", err)
		}
		// 容器注入的 TenantDomainCache 是零值实例（Domain map 未初始化），
		// 生产环境由 InitTenantDomain 落数据时初始化；测试里显式补上。
		if ts.Domain.Domain == nil {
			ts.Domain.Domain = &mapPg.SafeMap[string, string]{}
		}
		if ts.Domain.DomainData == nil {
			ts.Domain.DomainData = &mapPg.SafeMap[string, dtoTenantDomain.Dto]{}
		}
		cacheDiplPg.Set(testDiplKey, cacheDiplPg.DiplCo{
			Key: testDiplKey, Secret: testDiplSecret,
			TenantNo: testTenantNo, Ano: testAno, Name: "测试应用", No: "test-dipl",
		})

		gin.SetMode(gin.TestMode)
		engine := gin.New()
		ts.Bookmark.RegisterRoutes(engine)
		ts.Category.RegisterRoutes(engine)
		fn(t, ts, engine, db)
	})
}

// ---- 书签接口 ----

func TestBookmarkAPI(t *testing.T) {
	runAPIStack(t, func(t *testing.T, ts *apiTestStack, engine *gin.Engine, db *gorm.DB) {
		fx := seedBookmarkFixtures(t, db)

		t.Run("routes", func(t *testing.T) {
			want := map[string]bool{
				pathBookmarkGetAll: true, pathBookmarkGetMy: true, pathBookmarkGetOrg: true,
				pathBookmarkCreate: true, pathBookmarkUpdate: true, pathBookmarkDelete: true,
			}
			for _, r := range engine.Routes() {
				if _, ok := want[r.Path]; !ok {
					continue
				}
				if r.Method != http.MethodPost {
					t.Fatalf("路由 %s 方法 = %s, want POST", r.Path, r.Method)
				}
				delete(want, r.Path)
			}
			if len(want) != 0 {
				t.Fatalf("缺少路由: %v", want)
			}
		})

		t.Run("auth-required", func(t *testing.T) {
			env := postJSON(t, engine, pathBookmarkGetAll, `{}`, "")
			requireError(t, env, msgAuthRequired)
		})

		t.Run("getAll", func(t *testing.T) {
			env := postOK(t, engine, pathBookmarkGetAll, `{}`)
			all := decodeData[bookmarkVoAll](t, env)

			if len(all.My) != 1 || all.My[0].Name != "我的书签" {
				t.Fatalf("my = %+v, want 仅一条「我的书签」", all.My)
			}
			if len(all.Team) != 1 || all.Team[0].Name != "团队书签" {
				t.Fatalf("team = %+v, want 仅一条「团队书签」", all.Team)
			}
			if len(all.Org) != 1 || all.Org[0].Name != "组织书签" {
				t.Fatalf("org = %+v, want 仅一条「组织书签」", all.Org)
			}
			if len(all.MyCategory) != 1 || all.MyCategory[0].ID != idString(fx.myCat.ID) {
				t.Fatalf("MyCategory = %+v, want 仅包含我的分类 %d", all.MyCategory, fx.myCat.ID)
			}
			if len(all.TeamCategory) != 1 || all.TeamCategory[0].ID != idString(fx.teamCat.ID) {
				t.Fatalf("teamCategory = %+v", all.TeamCategory)
			}
			if len(all.OrgCategory) != 1 || all.OrgCategory[0].ID != idString(fx.orgCat.ID) {
				t.Fatalf("orgCategory = %+v", all.OrgCategory)
			}
			// 当前实现不填充 categoryName。
			if all.My[0].CategoryName != "" {
				t.Fatalf("categoryName = %q, want 空（当前实现未填充）", all.My[0].CategoryName)
			}
			// 其它租户 / 停用数据不能出现。
			if all.My[0].ID == idString(fx.otherTenant.ID) {
				t.Fatalf("其它租户数据不应出现在 my 中")
			}
		})

		t.Run("getMy", func(t *testing.T) {
			env := postOK(t, engine, pathBookmarkGetMy, `{}`)
			list := decodeData[[]bookmarkVo](t, env)
			if len(list) != 1 || list[0].Name != "我的书签" {
				t.Fatalf("getMy = %+v, want 仅一条「我的书签」", list)
			}
			if list[0].CategoryNo != fx.myCat.No {
				t.Fatalf("categoryNo = %q, want %q", list[0].CategoryNo, fx.myCat.No)
			}
		})

		t.Run("getOrg", func(t *testing.T) {
			env := postOK(t, engine, pathBookmarkGetOrg, `{}`)
			list := decodeData[[]bookmarkVo](t, env)
			if len(list) != 1 || list[0].Name != "组织书签" {
				t.Fatalf("getOrg = %+v, want 仅一条「组织书签」", list)
			}
		})

		t.Run("create-success", func(t *testing.T) {
			before := time.Now().Add(-time.Minute)
			env := postOK(t, engine, pathBookmarkCreate,
				`{"name":"新书签","urlSource":"https://example.com/new","tags":["golang","api"],"categoryNo":"c-1","orgNo":"o-1"}`)
			id := mustID(t, decodeData[string](t, env))
			e := getBookmark(t, db, id)

			assertEqual(t, "name", e.Name, "新书签")
			assertEqual(t, "urlSource", e.UrlSource, "https://example.com/new")
			assertEqual(t, "categoryNo", e.CategoryNo, "c-1")
			assertEqual(t, "orgNo", e.OrgNo, "o-1")
			assertEqual(t, "state", e.State, int8(1))
			// 创建不写入租户/操作员，落库为数据库默认值。
			assertEqual(t, "tenantNo", e.TenantNo, "")
			assertEqual(t, "ano", e.Ano, "")
			// 类型未写入，由数据库默认值兜底。
			assertEqual(t, "typeOwner", e.TypeOwner, "general")
			// tags 为 []string -> datatypes.JSONType 无法复制，当前不会持久化。
			if got := e.Tags.Data(); len(got) != 0 {
				t.Fatalf("tags = %v, want 空（当前实现不持久化 tags）", got)
			}
			if e.PublishTime == nil {
				t.Fatalf("publishTime 未自动写入")
			}
			if e.PublishTime.Before(before) {
				t.Fatalf("publishTime = %v, want 不早于 %v", e.PublishTime, before)
			}
		})

		t.Run("create-name-required", func(t *testing.T) {
			env := postAuth(t, engine, pathBookmarkCreate, `{"urlSource":"https://example.com/x"}`)
			requireError(t, env, "名称不能为空")
		})

		t.Run("create-url-required", func(t *testing.T) {
			env := postAuth(t, engine, pathBookmarkCreate, `{"name":"只有名称"}`)
			requireError(t, env, "来源地址不能为空")
		})

		t.Run("update-success", func(t *testing.T) {
			id, _ := createBookmark(t, engine, `{"name":"待更新","urlSource":"https://example.com/old"}`)
			env := postOK(t, engine, pathBookmarkUpdate, fmt.Sprintf(
				`{"id":"%d","name":"已更新","urlSource":"https://example.com/updated","publishTime":"2024-01-02T03:04:05Z"}`, id))
			assertEqual(t, "data", decodeData[string](t, env), idString(id))

			e := getBookmark(t, db, id)
			assertEqual(t, "name", e.Name, "已更新")
			assertEqual(t, "urlSource", e.UrlSource, "https://example.com/updated")
			want := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)
			if e.PublishTime == nil || !e.PublishTime.UTC().Equal(want) {
				t.Fatalf("publishTime = %v, want %v", e.PublishTime, want)
			}
		})

		t.Run("update-id-required", func(t *testing.T) {
			env := postAuth(t, engine, pathBookmarkUpdate, `{"name":"x","urlSource":"https://example.com/x"}`)
			requireError(t, env, "ID不能为空")
		})

		t.Run("update-not-found", func(t *testing.T) {
			env := postAuth(t, engine, pathBookmarkUpdate, `{"id":"999999","name":"x","urlSource":"https://example.com/x"}`)
			requireError(t, env, "数据不存在")
		})

		t.Run("update-name-required", func(t *testing.T) {
			id, _ := createBookmark(t, engine, `{"name":"占位","urlSource":"https://example.com/x"}`)
			env := postAuth(t, engine, pathBookmarkUpdate,
				fmt.Sprintf(`{"id":"%d","urlSource":"https://example.com/x"}`, id))
			requireError(t, env, "名称不能为空")
		})

		t.Run("delete-success", func(t *testing.T) {
			id, _ := createBookmark(t, engine, `{"name":"待删除","urlSource":"https://example.com/x"}`)
			env := postOK(t, engine, pathBookmarkDelete, fmt.Sprintf(`{"ids":["%d"]}`, id))
			assertEmptyData(t, env)
			if n := countBookmark(t, db, id); n != 0 {
				t.Fatalf("删除后仍存在 %d 条记录（pg.data.delete=true 应为物理删除）", n)
			}
		})

		t.Run("delete-not-found", func(t *testing.T) {
			env := postAuth(t, engine, pathBookmarkDelete, `{"ids":["999999"]}`)
			requireError(t, env, "数据不存在")
		})

		t.Run("delete-empty-ids", func(t *testing.T) {
			// 模型上的 validate tag 由 gin 默认校验器（tag=binding）识别不到，
			// 字段级校验不会触发，空 ids 由服务层兜底。
			requireError(t, postAuth(t, engine, pathBookmarkDelete, `{}`), "id错误")
			requireError(t, postAuth(t, engine, pathBookmarkDelete, `{"ids":[]}`), "id错误")
		})
	})
}

// ---- 书签分类接口 ----

func TestBookmarkCategoryAPI(t *testing.T) {
	runAPIStack(t, func(t *testing.T, ts *apiTestStack, engine *gin.Engine, db *gorm.DB) {
		fx := seedCategoryFixtures(t, db)

		t.Run("routes", func(t *testing.T) {
			want := map[string]bool{
				pathCategoryGetAll: true, pathCategoryCreateUpdate: true, pathCategoryDelete: true,
			}
			for _, r := range engine.Routes() {
				if _, ok := want[r.Path]; !ok {
					continue
				}
				if r.Method != http.MethodPost {
					t.Fatalf("路由 %s 方法 = %s, want POST", r.Path, r.Method)
				}
				delete(want, r.Path)
			}
			if len(want) != 0 {
				t.Fatalf("缺少路由: %v", want)
			}
		})

		t.Run("auth-required", func(t *testing.T) {
			env := postJSON(t, engine, pathCategoryGetAll, `{}`, "")
			requireError(t, env, msgAuthRequired)
		})

		t.Run("getAll", func(t *testing.T) {
			env := postOK(t, engine, pathCategoryGetAll, `{}`)
			all := decodeData[categoryVoAll](t, env)

			if len(all.My) != 1 || all.My[0].Name != "我的分类" {
				t.Fatalf("my = %+v, want 仅一条「我的分类」", all.My)
			}
			if len(all.Team) != 1 || all.Team[0].Name != "团队分类" {
				t.Fatalf("team = %+v", all.Team)
			}
			if len(all.Org) != 1 || all.Org[0].Name != "组织分类" {
				t.Fatalf("org = %+v", all.Org)
			}
			assertEqual(t, "my.id", all.My[0].ID, idString(fx.myCat.ID))
			assertEqual(t, "my.code", all.My[0].Code, "cat-my")
			assertEqual(t, "my.no", all.My[0].No, "catno-my")
		})

		t.Run("create-auto-code", func(t *testing.T) {
			id := createCategory(t, engine, `{"name":"新分类"}`)
			e := getCategory(t, db, id)

			assertEqual(t, "name", e.Name, "新分类")
			if e.Code == "" || e.Code == "系统自动建立" {
				t.Fatalf("code = %q, want 自动生成的非空标志", e.Code)
			}
			if e.No == "" {
				t.Fatalf("no 不应为空")
			}
			// 分类创建使用 holderPg（读 AUTH_LOGIN），API 中间件只设置 AUTH_LOGIN_API，
			// 因此租户号落到 "-1"。
			assertEqual(t, "tenantNo", e.TenantNo, "-1")
			assertEqual(t, "ano", e.Ano, "")
			assertEqual(t, "state", e.State, int8(1))
			assertEqual(t, "typeOwner", e.TypeOwner, "general")
			assertEqual(t, "parentNo", e.ParentNo, "")
			assertEqual(t, "parentId", e.ParentId, "")
			if !strings.Contains(e.IdLink, idString(id)) {
				t.Fatalf("idLink = %q, want 含自身 ID %d", e.IdLink, id)
			}
			if !strings.Contains(e.NoLink, e.No) {
				t.Fatalf("noLink = %q, want 含自身 No %q", e.NoLink, e.No)
			}
		})

		t.Run("create-empty-name", func(t *testing.T) {
			requireError(t, postAuth(t, engine, pathCategoryCreateUpdate, `{"name":""}`), "名称不能为空")
			requireError(t, postAuth(t, engine, pathCategoryCreateUpdate, `{}`), "名称不能为空")
		})

		t.Run("create-parent-not-found", func(t *testing.T) {
			requireError(t, postAuth(t, engine, pathCategoryCreateUpdate,
				`{"name":"子分类","parentNo":"no-such-parent"}`), "上级不存在")
		})

		t.Run("create-with-parent", func(t *testing.T) {
			parentID := createCategory(t, engine, `{"name":"父分类"}`)
			parent := getCategory(t, db, parentID)

			childID := createCategory(t, engine,
				fmt.Sprintf(`{"name":"子分类","parentNo":"%s"}`, parent.No))
			child := getCategory(t, db, childID)

			assertEqual(t, "parentNo", child.ParentNo, parent.No)
			assertEqual(t, "parentId", child.ParentId, idString(parent.ID))
			if !strings.Contains(child.IdLink, idString(parent.ID)) || !strings.Contains(child.IdLink, idString(childID)) {
				t.Fatalf("idLink = %q, want 同时含父 ID %d 与自身 ID %d", child.IdLink, parent.ID, childID)
			}
			if !strings.Contains(child.NoLink, parent.No) || !strings.Contains(child.NoLink, child.No) {
				t.Fatalf("noLink = %q, want 同时含父 No 与自身 No", child.NoLink)
			}
		})

		t.Run("create-when-id-zero", func(t *testing.T) {
			// 服务层以 ID > 0 判定更新，ID=0 时走创建分支。
			env := postOK(t, engine, pathCategoryCreateUpdate, `{"id":"0","name":"ID 为零也创建"}`)
			id := mustID(t, decodeData[string](t, env))
			assertEqual(t, "name", getCategory(t, db, id).Name, "ID 为零也创建")
		})

		t.Run("update-success", func(t *testing.T) {
			id := createCategory(t, engine, `{"name":"原名"}`)
			before := getCategory(t, db, id)

			env := postOK(t, engine, pathCategoryCreateUpdate, fmt.Sprintf(`{"id":"%d","name":"新名"}`, id))
			// 更新分支返回空字符串。
			assertEmptyData(t, env)

			after := getCategory(t, db, id)
			assertEqual(t, "name", after.Name, "新名")
			assertEqual(t, "code", after.Code, before.Code)
			assertEqual(t, "no", after.No, before.No)
			assertEqual(t, "tenantNo", after.TenantNo, before.TenantNo)
		})

		t.Run("update-not-found", func(t *testing.T) {
			requireError(t, postAuth(t, engine, pathCategoryCreateUpdate,
				`{"id":"999999","name":"x"}`), "数据不存在")
		})

		t.Run("update-empty-name", func(t *testing.T) {
			id := createCategory(t, engine, `{"name":"占位"}`)
			requireError(t, postAuth(t, engine, pathCategoryCreateUpdate,
				fmt.Sprintf(`{"id":"%d","name":""}`, id)), "名称不能为空")
		})

		t.Run("update-parent-self", func(t *testing.T) {
			id := createCategory(t, engine, `{"name":"自己"}`)
			e := getCategory(t, db, id)
			requireError(t, postAuth(t, engine, pathCategoryCreateUpdate,
				fmt.Sprintf(`{"id":"%d","name":"自己","parentNo":"%s"}`, id, e.No)), "上级不能等于自己")
		})

		t.Run("update-parent-cycle", func(t *testing.T) {
			parentID := createCategory(t, engine, `{"name":"循环父"}`)
			parent := getCategory(t, db, parentID)
			childID := createCategory(t, engine,
				fmt.Sprintf(`{"name":"循环子","parentNo":"%s"}`, parent.No))
			child := getCategory(t, db, childID)

			requireError(t, postAuth(t, engine, pathCategoryCreateUpdate,
				fmt.Sprintf(`{"id":"%d","name":"循环父","parentNo":"%s"}`, parentID, child.No)),
				"无法保存，不能设置为自己的子集")
		})

		t.Run("update-parent", func(t *testing.T) {
			aID := createCategory(t, engine, `{"name":"A"}`)
			bID := createCategory(t, engine, `{"name":"B"}`)
			b := getCategory(t, db, bID)

			postOK(t, engine, pathCategoryCreateUpdate,
				fmt.Sprintf(`{"id":"%d","name":"A","parentNo":"%s"}`, aID, b.No))

			a := getCategory(t, db, aID)
			assertEqual(t, "parentNo", a.ParentNo, b.No)
			assertEqual(t, "parentId", a.ParentId, idString(b.ID))
		})

		t.Run("delete-success", func(t *testing.T) {
			id := createCategory(t, engine, `{"name":"待删除分类"}`)
			env := postOK(t, engine, pathCategoryDelete, fmt.Sprintf(`{"ids":["%d"]}`, id))
			assertEmptyData(t, env)
			if n := countCategory(t, db, id); n != 0 {
				t.Fatalf("删除后仍存在 %d 条记录（pg.data.delete=true 应为物理删除）", n)
			}
		})

		t.Run("delete-not-found", func(t *testing.T) {
			requireError(t, postAuth(t, engine, pathCategoryDelete, `{"ids":["999999"]}`), "数据不存在")
		})

		t.Run("delete-empty-ids", func(t *testing.T) {
			requireError(t, postAuth(t, engine, pathCategoryDelete, `{}`), "id错误")
		})
	})
}

// ---- 数据准备 ----

type bookmarkFixtures struct {
	myCat       entityBlog.BlogBookmarkCategoryEntity
	teamCat     entityBlog.BlogBookmarkCategoryEntity
	orgCat      entityBlog.BlogBookmarkCategoryEntity
	otherTenant entityBlog.BlogBookmarkEntity
}

func seedBookmarkFixtures(t *testing.T, db *gorm.DB) bookmarkFixtures {
	t.Helper()
	fx := bookmarkFixtures{
		myCat: entityBlog.BlogBookmarkCategoryEntity{
			Name: "我的分类", Code: "cat-my", No: "catno-my",
			State: 1, TenantNo: testTenantNo, Ano: testAno, TypeOwner: "my",
		},
		teamCat: entityBlog.BlogBookmarkCategoryEntity{
			Name: "团队分类", Code: "cat-team", No: "catno-team",
			State: 1, TenantNo: testTenantNo, Ano: testAno, TypeOwner: "team",
		},
		orgCat: entityBlog.BlogBookmarkCategoryEntity{
			Name: "组织分类", Code: "cat-org", No: "catno-org",
			State: 1, TenantNo: testTenantNo, Ano: testAno, TypeOwner: "org",
		},
	}
	for _, v := range []any{&fx.myCat, &fx.teamCat, &fx.orgCat} {
		mustCreate(t, db, v)
	}

	fx.otherTenant = entityBlog.BlogBookmarkEntity{
		Name: "其它租户书签", UrlSource: "https://example.com/other",
		State: 1, TenantNo: "2000", Ano: "u-002", TypeOwner: "my",
	}
	rows := []any{
		&entityBlog.BlogBookmarkEntity{Name: "我的书签", UrlSource: "https://example.com/my",
			State: 1, TenantNo: testTenantNo, Ano: testAno, TypeOwner: "my", CategoryNo: fx.myCat.No},
		&entityBlog.BlogBookmarkEntity{Name: "团队书签", UrlSource: "https://example.com/team",
			State: 1, TenantNo: testTenantNo, Ano: testAno, TypeOwner: "team"},
		&entityBlog.BlogBookmarkEntity{Name: "组织书签", UrlSource: "https://example.com/org",
			State: 1, TenantNo: testTenantNo, Ano: testAno, TypeOwner: "org"},
		&entityBlog.BlogBookmarkEntity{Name: "同租户其它账号", UrlSource: "https://example.com/other-ano",
			State: 1, TenantNo: testTenantNo, Ano: "u-999", TypeOwner: "my"},
		&entityBlog.BlogBookmarkEntity{Name: "已停用书签", UrlSource: "https://example.com/disabled",
			State: 2, TenantNo: testTenantNo, Ano: testAno, TypeOwner: "my"},
		&fx.otherTenant,
	}
	for _, v := range rows {
		mustCreate(t, db, v)
	}
	return fx
}

type categoryFixtures struct {
	myCat entityBlog.BlogBookmarkCategoryEntity
}

func seedCategoryFixtures(t *testing.T, db *gorm.DB) categoryFixtures {
	t.Helper()
	fx := categoryFixtures{
		myCat: entityBlog.BlogBookmarkCategoryEntity{
			Name: "我的分类", Code: "cat-my", No: "catno-my",
			State: 1, TenantNo: testTenantNo, Ano: testAno, TypeOwner: "my",
		},
	}
	rows := []any{
		&fx.myCat,
		&entityBlog.BlogBookmarkCategoryEntity{Name: "团队分类", Code: "cat-team", No: "catno-team",
			State: 1, TenantNo: testTenantNo, Ano: testAno, TypeOwner: "team"},
		&entityBlog.BlogBookmarkCategoryEntity{Name: "组织分类", Code: "cat-org", No: "catno-org",
			State: 1, TenantNo: testTenantNo, Ano: testAno, TypeOwner: "org"},
		&entityBlog.BlogBookmarkCategoryEntity{Name: "其它租户分类", Code: "cat-other", No: "catno-other",
			State: 1, TenantNo: "2000", Ano: "u-002", TypeOwner: "my"},
		&entityBlog.BlogBookmarkCategoryEntity{Name: "已停用分类", Code: "cat-disabled", No: "catno-disabled",
			State: 2, TenantNo: testTenantNo, Ano: testAno, TypeOwner: "my"},
	}
	for _, v := range rows {
		mustCreate(t, db, v)
	}
	return fx
}

// ---- HTTP 请求辅助 ----

func testAuthHeader() string {
	return testDiplKey + ":" + cacheDiplPg.HashSha(testDiplKey, testDiplSecret)
}

func postJSON(t *testing.T, engine *gin.Engine, path, body, auth string) rsEnvelope {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, path, bytes.NewBufferString(body))
	req.Host = "localhost"
	req.Header.Set("Content-Type", "application/json")
	if auth != "" {
		req.Header.Set(authHeaderName, auth)
	}
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("POST %s http status = %d, want 200, body=%s", path, rec.Code, strings.TrimSpace(rec.Body.String()))
	}
	var env rsEnvelope
	if err := json.Unmarshal(rec.Body.Bytes(), &env); err != nil {
		t.Fatalf("POST %s 响应解析失败: %v (body=%s)", path, err, rec.Body.String())
	}
	return env
}

func postAuth(t *testing.T, engine *gin.Engine, path, body string) rsEnvelope {
	t.Helper()
	return postJSON(t, engine, path, body, testAuthHeader())
}

func postOK(t *testing.T, engine *gin.Engine, path, body string) rsEnvelope {
	t.Helper()
	env := postAuth(t, engine, path, body)
	requireCode(t, env, "200")
	requireMessage(t, env, msgOK)
	return env
}

func createBookmark(t *testing.T, engine *gin.Engine, body string) (int64, rsEnvelope) {
	t.Helper()
	env := postOK(t, engine, pathBookmarkCreate, body)
	return mustID(t, decodeData[string](t, env)), env
}

func createCategory(t *testing.T, engine *gin.Engine, body string) int64 {
	t.Helper()
	env := postOK(t, engine, pathCategoryCreateUpdate, body)
	return mustID(t, decodeData[string](t, env))
}

// ---- 响应断言辅助 ----

func requireCode(t *testing.T, env rsEnvelope, want string) {
	t.Helper()
	if env.Code != want {
		t.Fatalf("code = %q, want %q (message=%q, data=%s)", env.Code, want, env.Message, env.Data)
	}
}

func requireMessage(t *testing.T, env rsEnvelope, want string) {
	t.Helper()
	if env.Message != want {
		t.Fatalf("message = %q, want %q (code=%q)", env.Message, want, env.Code)
	}
}

// requireError 断言业务错误：code=500、message 为指定文案、data 为空、无 messageData。
func requireError(t *testing.T, env rsEnvelope, msg string) {
	t.Helper()
	requireCode(t, env, "500")
	requireMessage(t, env, msg)
	assertEmptyData(t, env)
	if len(env.MessageData) != 0 {
		t.Fatalf("messageData = %s, want 空", env.MessageData)
	}
}

func assertEmptyData(t *testing.T, env rsEnvelope) {
	t.Helper()
	s := strings.TrimSpace(string(env.Data))
	if s != "" && s != `""` && s != "null" {
		t.Fatalf("data = %s, want 空", s)
	}
}

func assertEqual[T comparable](t *testing.T, field string, got, want T) {
	t.Helper()
	if got != want {
		t.Fatalf("%s = %v, want %v", field, got, want)
	}
}

func decodeData[T any](t *testing.T, env rsEnvelope) T {
	t.Helper()
	var v T
	if len(env.Data) == 0 || string(env.Data) == "null" {
		return v
	}
	if err := json.Unmarshal(env.Data, &v); err != nil {
		t.Fatalf("data 解析失败: %v (raw=%s)", err, env.Data)
	}
	return v
}

// ---- 数据库辅助 ----

func mustCreate(t *testing.T, db *gorm.DB, v any) {
	t.Helper()
	if err := db.Create(v).Error; err != nil {
		t.Fatalf("插入测试数据失败: %v", err)
	}
}

func getBookmark(t *testing.T, db *gorm.DB, id int64) entityBlog.BlogBookmarkEntity {
	t.Helper()
	var e entityBlog.BlogBookmarkEntity
	if err := db.First(&e, "id = ?", id).Error; err != nil {
		t.Fatalf("查询书签 %d 失败: %v", id, err)
	}
	return e
}

func getCategory(t *testing.T, db *gorm.DB, id int64) entityBlog.BlogBookmarkCategoryEntity {
	t.Helper()
	var e entityBlog.BlogBookmarkCategoryEntity
	if err := db.First(&e, "id = ?", id).Error; err != nil {
		t.Fatalf("查询分类 %d 失败: %v", id, err)
	}
	return e
}

func countBookmark(t *testing.T, db *gorm.DB, id int64) int64 {
	t.Helper()
	var n int64
	if err := db.Model(&entityBlog.BlogBookmarkEntity{}).Where("id = ?", id).Count(&n).Error; err != nil {
		t.Fatalf("统计书签失败: %v", err)
	}
	return n
}

func countCategory(t *testing.T, db *gorm.DB, id int64) int64 {
	t.Helper()
	var n int64
	if err := db.Model(&entityBlog.BlogBookmarkCategoryEntity{}).Where("id = ?", id).Count(&n).Error; err != nil {
		t.Fatalf("统计分类失败: %v", err)
	}
	return n
}

func idString(id int64) string {
	return strconv.FormatInt(id, 10)
}

func mustID(t *testing.T, s string) int64 {
	t.Helper()
	id, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		t.Fatalf("id = %q 不是合法数字: %v", s, err)
	}
	return id
}
