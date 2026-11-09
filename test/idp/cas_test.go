package idp_e2e

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

type casLoginVO struct {
	ServiceTicket string `json:"serviceTicket"`
	Service       string `json:"service"`
	RedirectUrl   string `json:"redirectUrl"`
}

// casLoginDirect 直接调用 CasService.Login，并注入登录态 ano。
// 生产路由 /xianfu/auth/manage/cas/login 与 /xianfu/auth/manage/idp/cas/login
// 均未挂载鉴权中间件，ctx 中不会出现 ano，因此 HTTP 直连必然返回“用户未登录”；
// 这里用于验证服务层票据签发逻辑本身。
func (e *env) casLoginDirect(t *testing.T, ano, service string) *apiResp {
	t.Helper()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, manageCasURL+"/login",
		strings.NewReader(js(map[string]any{"service": service, "sourceNo": ""})))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Header.Set("User-Agent", "dsh-idp-e2e/1.0")
	if ano != "" {
		c.Set("ano", ano)
	}
	e.ts.Cas.Login(c)
	return parseResp(w)
}

func runCAS(t *testing.T, e *env) {
	e.t = t
	service := "https://app.example.org/cas/callback"

	newTicket := func(t *testing.T) string {
		t.Helper()
		r := e.casLoginDirect(t, testTenantNo, service)
		if r.Code != "200" {
			t.Fatalf("CAS 直接登录失败: code=%s msg=%s", r.Code, r.Message)
		}
		vo := decodeData[casLoginVO](t, r)
		if !strings.HasPrefix(vo.ServiceTicket, "ST-") || !strings.Contains(vo.RedirectUrl, "ticket=") {
			t.Fatalf("票据格式异常: %+v", vo)
		}
		return vo.ServiceTicket
	}

	t.Run("service-required", func(t *testing.T) {
		e.t = t
		// CasLoginCt.Service 带 binding:"required"，控制器 ShouldBind 先于服务层校验失败。
		r := e.casLoginDirect(t, testTenantNo, "")
		if r.Code == "200" {
			t.Fatalf("缺 service 期望失败: code=%s msg=%s", r.Code, r.Message)
		}
		if !strings.Contains(r.Message, "操作失败") && !strings.Contains(r.Message, "service 参数不能为空") {
			t.Fatalf("缺 service 失败提示异常: code=%s msg=%s", r.Code, r.Message)
		}
	})
	t.Run("not-logged-in", func(t *testing.T) {
		e.t = t
		r := e.casLoginDirect(t, "", service)
		if r.Code == "200" || !strings.Contains(r.Message, "用户未登录") {
			t.Fatalf("无 ano 期望未登录: code=%s msg=%s", r.Code, r.Message)
		}
	})
	t.Run("account-not-found", func(t *testing.T) {
		e.t = t
		r := e.casLoginDirect(t, "no-such-account", service)
		if r.Code == "200" || !strings.Contains(r.Message, "用户不存在") {
			t.Fatalf("不存在账号期望失败: code=%s msg=%s", r.Code, r.Message)
		}
	})
	t.Run("http-route-has-no-auth-context", func(t *testing.T) {
		e.t = t
		// 路由未挂载鉴权中间件，直接 POST 无法建立 ano 上下文。
		for _, path := range []string{manageCasURL + "/login", authIdpURL + "/cas/login"} {
			r := e.postNoAuth(path, js(map[string]any{"service": service}))
			if r.Code == "200" || !strings.Contains(r.Message, "用户未登录") {
				t.Fatalf("%s 期望用户未登录，实际 code=%s msg=%s", path, r.Code, r.Message)
			}
		}
	})
	t.Run("login-success", func(t *testing.T) {
		e.t = t
		st := newTicket(t)
		if st == "" {
			t.Fatalf("空票据")
		}
	})

	// ---- CAS 1.0 /validate ----
	t.Run("v1-valid", func(t *testing.T) {
		e.t = t
		st := newTicket(t)
		r := e.get(manageCasURL + "/validate?service=" + url.QueryEscape(service) + "&ticket=" + url.QueryEscape(st))
		if !strings.HasPrefix(r.Body, "yes\n") || !strings.Contains(r.Body, "manage") {
			t.Fatalf("CAS1.0 校验通过响应异常: %q", r.Body)
		}
	})
	t.Run("v1-ticket-reuse", func(t *testing.T) {
		e.t = t
		st := newTicket(t)
		_ = e.get(manageCasURL + "/validate?service=" + url.QueryEscape(service) + "&ticket=" + url.QueryEscape(st))
		r := e.get(manageCasURL + "/validate?service=" + url.QueryEscape(service) + "&ticket=" + url.QueryEscape(st))
		if !strings.HasPrefix(r.Body, "no\n") {
			t.Fatalf("票据复用应失败，实际 %q", r.Body)
		}
	})
	t.Run("v1-service-mismatch", func(t *testing.T) {
		e.t = t
		st := newTicket(t)
		r := e.get(manageCasURL + "/validate?service=" + url.QueryEscape("https://other.example.org/cb") + "&ticket=" + url.QueryEscape(st))
		if !strings.HasPrefix(r.Body, "no\n") {
			t.Fatalf("service 不匹配应失败，实际 %q", r.Body)
		}
	})
	t.Run("v1-missing-params", func(t *testing.T) {
		e.t = t
		r := e.get(manageCasURL + "/validate")
		if !strings.HasPrefix(r.Body, "no\n") {
			t.Fatalf("缺参应返回 no，实际 %q", r.Body)
		}
	})

	// ---- CAS 2.0 /serviceValidate ----
	t.Run("v2-valid-xml", func(t *testing.T) {
		e.t = t
		st := newTicket(t)
		r := e.get(manageCasURL + "/serviceValidate?service=" + url.QueryEscape(service) + "&ticket=" + url.QueryEscape(st))
		if !strings.Contains(r.Body, "cas:authenticationSuccess") || !strings.Contains(r.Body, "<cas:user>manage</cas:user>") {
			t.Fatalf("CAS2.0 成功响应异常: %s", r.Body)
		}
		if ct := r.Header.Get("Content-Type"); !strings.Contains(ct, "xml") {
			t.Fatalf("Content-Type=%q", ct)
		}
	})
	t.Run("v2-not-st-ticket", func(t *testing.T) {
		e.t = t
		r := e.get(manageCasURL + "/serviceValidate?service=" + url.QueryEscape(service) + "&ticket=BAD-1")
		if !strings.Contains(r.Body, "INVALID_TICKET_SPEC") {
			t.Fatalf("非 ST 票据应返回 INVALID_TICKET_SPEC: %s", r.Body)
		}
	})
	t.Run("v2-unknown-ticket", func(t *testing.T) {
		e.t = t
		r := e.get(manageCasURL + "/serviceValidate?service=" + url.QueryEscape(service) + "&ticket=ST-99999999")
		if !strings.Contains(r.Body, "INVALID_TICKET") || strings.Contains(r.Body, "INVALID_TICKET_SPEC") {
			t.Fatalf("未知 ST 应返回 INVALID_TICKET: %s", r.Body)
		}
	})
	t.Run("v2-service-mismatch", func(t *testing.T) {
		e.t = t
		st := newTicket(t)
		r := e.get(manageCasURL + "/serviceValidate?service=" + url.QueryEscape("https://other.example.org/cb") + "&ticket=" + url.QueryEscape(st))
		if !strings.Contains(r.Body, "INVALID_SERVICE") {
			t.Fatalf("service 不匹配应返回 INVALID_SERVICE: %s", r.Body)
		}
	})
	t.Run("v2-missing-params", func(t *testing.T) {
		e.t = t
		r := e.get(manageCasURL + "/serviceValidate")
		if !strings.Contains(r.Body, "INVALID_REQUEST") {
			t.Fatalf("缺参应返回 INVALID_REQUEST: %s", r.Body)
		}
	})

	// ---- CAS 3.0 /p3/serviceValidate ----
	t.Run("v3-valid-xml", func(t *testing.T) {
		e.t = t
		st := newTicket(t)
		r := e.get(manageCasURL + "/p3/serviceValidate?service=" + url.QueryEscape(service) + "&ticket=" + url.QueryEscape(st))
		if !strings.Contains(r.Body, "cas:authenticationSuccess") || !strings.Contains(r.Body, "<cas:user>manage</cas:user>") {
			t.Fatalf("CAS3.0 成功响应异常: %s", r.Body)
		}
	})
}
