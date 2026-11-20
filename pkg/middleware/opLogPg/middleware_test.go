package opLogPg

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/entityBasic"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestTruncateKeepsValidUTF8(t *testing.T) {
	// 5 会落在“文”字中间，回退后应只剩“中”。
	if got := truncate("中文字符", 5); got != "中" {
		t.Fatalf("截断结果应为 %q，实际 %q", "中", got)
	}
	if got := truncate("abc", 10); got != "abc" {
		t.Fatalf("未超长时不应截断，实际 %q", got)
	}
}

func TestSanitizeTraceID(t *testing.T) {
	cases := map[string]string{
		"  abc-123  ":            "abc-123",
		"\x00 abc \n":            "abc",
		"":                       "",
		strings.Repeat("a", 100): strings.Repeat("a", 80),
	}
	for in, want := range cases {
		if got := sanitizeTraceID(in); got != want {
			t.Fatalf("sanitizeTraceID(%q) = %q，期望 %q", in, got, want)
		}
	}
}

func TestContentTypeFilters(t *testing.T) {
	fileTypes := []string{
		"multipart/form-data; boundary=x",
		"application/octet-stream",
		"image/png",
		"application/pdf",
		"application/zip",
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	}
	for _, ct := range fileTypes {
		if !isFileStream(ct) {
			t.Fatalf("应识别为文件流: %q", ct)
		}
	}
	textTypes := []string{
		"application/json",
		"application/json; charset=utf-8",
		"application/vnd.api+json",
		"application/x-www-form-urlencoded",
		"application/xml",
		"text/plain",
	}
	for _, ct := range textTypes {
		if isFileStream(ct) || !isTextContentType(ct) {
			t.Fatalf("应识别为文本: %q", ct)
		}
	}
}

func TestCaptureRequestBodyJSONAndReusable(t *testing.T) {
	const body = `{"name":"测试"}`
	req := httptest.NewRequest(http.MethodPost, "/xianfu/manage/basic/operation-log/query", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = req

	if got := captureRequestBody(c); got != body {
		t.Fatalf("请求体记录错误: %q", got)
	}
	// 中间件读取后，handler 仍必须能读到完整请求体。
	restored, err := io.ReadAll(c.Request.Body)
	if err != nil {
		t.Fatalf("读取放回后的请求体失败: %v", err)
	}
	if string(restored) != body {
		t.Fatalf("请求体未完整放回: %q", string(restored))
	}
}

func TestCaptureRequestBodySkipsFileStream(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/upload", strings.NewReader("file-content"))
	req.Header.Set("Content-Type", "multipart/form-data; boundary=x")
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = req

	if got := captureRequestBody(c); got != "" {
		t.Fatalf("文件流不应被记录，实际 %q", got)
	}
	raw, _ := io.ReadAll(c.Request.Body)
	if string(raw) != "file-content" {
		t.Fatalf("文件流请求体被改动: %q", string(raw))
	}
}

func TestShouldSkip(t *testing.T) {
	cases := []struct {
		method, path string
		skip         bool
	}{
		{http.MethodOptions, "/xianfu/manage/basic/user/query", true},
		{http.MethodGet, "/assets/app.css", true},
		{http.MethodGet, "/attachment/a.png", true},
		{http.MethodPost, "/xianfu/manage/basic/user/query", false},
	}
	for _, tc := range cases {
		req := httptest.NewRequest(tc.method, tc.path, nil)
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = req
		if got := shouldSkip(c); got != tc.skip {
			t.Fatalf("%s %s shouldSkip = %v，期望 %v", tc.method, tc.path, got, tc.skip)
		}
	}
}

func TestDeriveRoute(t *testing.T) {
	cases := []struct {
		path                 string
		root, module, action string
	}{
		{"/xianfu/manage/basic/attachment/upload", "basic", "attachment", "upload"},
		{"/xianfu/sys/ram-idp/idpBinding/query", "ram-idp", "idpBinding", "query"},
		{"/api-app/bookmark/create", "api-app", "bookmark", "create"},
		{"/xianfu/pub/publicKey", "", "pub", "publicKey"},
		{"/", "", "", ""},
	}
	for _, tc := range cases {
		root, module, action := deriveRoute(tc.path)
		if root != tc.root || module != tc.module || action != tc.action {
			t.Fatalf("deriveRoute(%q) = (%q,%q,%q)，期望 (%q,%q,%q)", tc.path, root, module, action, tc.root, tc.module, tc.action)
		}
	}
}

func TestHeaderJSONOnlyUserAgent(t *testing.T) {
	h := http.Header{}
	h.Set("Authorization", "Bearer secret")
	h.Set("Cookie", "token=abc")
	h.Set("X-Request-Id", "rid-1")
	h.Set("User-Agent", "curl/8")
	got := headerJSON(h)
	if !strings.Contains(got, `"User-Agent":"curl/8"`) {
		t.Fatalf("User-Agent 未记录: %s", got)
	}
	if strings.Contains(got, "Authorization") || strings.Contains(got, "Cookie") || strings.Contains(got, "X-Request-Id") {
		t.Fatalf("当前只记录 User-Agent，其它请求头不应入库: %s", got)
	}
}

func TestNormalizeCode(t *testing.T) {
	cases := map[string]string{"-1": "", "": "", "  t-1  ": "t-1"}
	for in, want := range cases {
		if got := normalizeCode(in); got != want {
			t.Fatalf("normalizeCode(%q) = %q，期望 %q", in, got, want)
		}
	}
}

func TestDeriveModuleActionFromConst(t *testing.T) {
	cases := []struct {
		path                                         string
		root, module, moduleName, action, actionName string
	}{
		{"/xianfu/manage/basic/operation-log/query", "basic.operationLog", "basic.operationLog", "操作日志", "query", "查询"},
		{"/xianfu/manage/basic/data-dictionary-sub/query", "basic.dataDictionary", "basic.dictionary.sub", "数据字典子项", "query", "查询"},
		{"/xianfu/sys/ram-idp/idpBinding/query", "idp.idpBinding", "idp.idpBinding", "身份绑定", "query", "查询"},
		{"/xianfu/sys/ram/account-mfa/query", "ram.accountAuthMfa", "ram.accountAuthMfa", "账号多因素认证", "query", "查询"},
		{"/xianfu/manage/basic/filePub/upload", "basic", "filePub", "filePub", "upload", "upload"},
	}
	for _, tc := range cases {
		root, module, moduleName, action, actionName := deriveModuleAction(tc.path)
		if root != tc.root || module != tc.module || moduleName != tc.moduleName || action != tc.action || actionName != tc.actionName {
			t.Fatalf("deriveModuleAction(%q) = (%q,%q,%q,%q,%q)，期望 (%q,%q,%q,%q,%q)",
				tc.path, root, module, moduleName, action, actionName,
				tc.root, tc.module, tc.moduleName, tc.action, tc.actionName)
		}
	}
}

func TestDeriveModuleActionUsesMatchMap(t *testing.T) {
	cases := []struct {
		path                                   string
		root, module, moduleName, action, name string
	}{
		{"/xianfu/manage/basic/attachment/create", "basic.attachment", "basic.attachment", "附件", "create", "创建"},
		{"/xianfu/manage/basic/attachment/createUpdateSimple", "basic.attachment", "basic.attachment", "附件", "createUpdate", "创建更新"},
		{"/xianfu/manage/basic/operation-log/exportExcel", "basic.operationLog", "basic.operationLog", "操作日志", "export", "导出"},
		{"/xianfu/manage/basic/attachment/physicalDelete", "basic.attachment", "basic.attachment", "附件", "physicalDeletion", "删除"},
		{"/xianfu/manage/basic/attachment/upload", "basic.attachment", "basic.attachment", "附件", "upload", "upload"},
		{"/xianfu/manage/basic/filePub/create", "basic", "filePub", "filePub", "create", "创建"},
	}
	for _, tc := range cases {
		root, module, moduleName, action, name := deriveModuleAction(tc.path)
		if root != tc.root || module != tc.module || moduleName != tc.moduleName || action != tc.action || name != tc.name {
			t.Fatalf("deriveModuleAction(%q) = (%q,%q,%q,%q,%q)，期望 (%q,%q,%q,%q,%q)",
				tc.path, root, module, moduleName, action, name,
				tc.root, tc.module, tc.moduleName, tc.action, tc.name)
		}
	}
}

func TestPathMatcherExclude(t *testing.T) {
	m := newPathMatcher([]string{"/xianfu/pub/*", "/xianfu/auth/**", "/xianfu/*/health"})
	cases := []struct {
		path string
		want bool
	}{
		{"/xianfu/pub/login", true},
		{"/xianfu/pub/login/x", true},
		{"/xianfu/pub", false},
		{"/xianfu/public", false},
		{"/xianfu/auth", true},
		{"/xianfu/auth/login", true},
		{"/xianfu/auth/login/x", true},
		{"/xianfu/authx", false},
		{"/xianfu/manage/health", true},
		{"/xianfu/manage/health/x", true},
		{"/xianfu/manage/foo/health", false},
		{"/xianfu/other", false},
	}
	for _, tc := range cases {
		if got := m.Match(tc.path); got != tc.want {
			t.Fatalf("Match(%q) = %v，期望 %v", tc.path, got, tc.want)
		}
	}
	if got := newPathMatcher(nil).Match("/xianfu/pub/login"); got {
		t.Fatalf("空规则不应命中")
	}
}

func TestGlobalSkipsExcludedPath(t *testing.T) {
	w := &Writer{ch: make(chan *entityBasic.BasicOperationLogEntity, 8)}
	defaultWriter.Store(w)
	t.Cleanup(func() { defaultWriter.Store(nil) })

	r := gin.New()
	mw := &Middleware{Config: Config{Exclude: []string{"/xianfu/pub/*", "/xianfu/auth/**"}}}
	r.Use(mw.Global())
	for _, path := range []string{"/xianfu/pub/login", "/xianfu/auth/login", "/xianfu/manage/account/query"} {
		r.GET(path, func(c *gin.Context) { c.Status(http.StatusOK) })
	}
	for _, path := range []string{"/xianfu/pub/login", "/xianfu/auth/login", "/xianfu/manage/account/query"} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
	}
	if got := len(w.ch); got != 1 {
		t.Fatalf("排除路径不应记录日志，未排除路径应记录 1 条，实际 %d 条", got)
	}
}
