package authPg

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/middleware/components/cachePg/cacheDiplPg"
)

func TestGroupApiMiddlewareFailsClosed(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(GroupApiMiddleware(&GroupApiMiddlewareSp{}))
	handlerCalled := false
	router.POST("/x", func(c *gin.Context) { handlerCalled = true })

	req := httptest.NewRequest(http.MethodPost, "/x", strings.NewReader("{}"))
	req.Header.Set("Authorization", "not-exists:bad")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	if handlerCalled {
		t.Fatal("无效密钥不应进入业务 handler")
	}
	if !strings.Contains(resp.Body.String(), "密钥无效") {
		t.Fatalf("应返回密钥无效: %s", resp.Body.String())
	}
}

func TestGroupApiMiddlewareAllowsValidKey(t *testing.T) {
	gin.SetMode(gin.TestMode)
	const (
		key    = "unit-test-key"
		secret = "unit-test-secret"
		tenant = "unit-tenant"
	)
	cacheDiplPg.Set(key, cacheDiplPg.DiplCo{
		Key:      key,
		Secret:   secret,
		HashSha:  cacheDiplPg.HashSha(key, secret),
		TenantNo: tenant,
		State:    1,
	})
	router := gin.New()
	router.Use(GroupApiMiddleware(&GroupApiMiddlewareSp{}))
	handlerCalled := false
	router.POST("/x", func(c *gin.Context) { handlerCalled = true })

	req := httptest.NewRequest(http.MethodPost, "/x", strings.NewReader("{}"))
	req.Header.Set("Authorization", key+":"+cacheDiplPg.HashSha(key, secret))
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	if !handlerCalled {
		t.Fatalf("有效密钥应进入业务: %s", resp.Body.String())
	}
}
