package multiTenantPg

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/consts/constContextPg"
)

func TestResolveMultiRule(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("未设置 CTX_RULE 时返回 nil", func(t *testing.T) {
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		if rule := resolveMultiRule(ctx); nil != rule {
			t.Fatalf("expected nil rule, got %+v", rule)
		}
	})

	t.Run("CTX_RULE 为 *MultiRule 时原样返回", func(t *testing.T) {
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		want := &MultiRule{Tenant: true, MultipleTenant: true}
		ctx.Set(constContextPg.CTX_RULE, want)
		if rule := resolveMultiRule(ctx); rule != want {
			t.Fatalf("expected same rule pointer, got %+v", rule)
		}
	})

	t.Run("CTX_RULE 类型异常时不 panic 且回退 nil", func(t *testing.T) {
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Set(constContextPg.CTX_RULE, "not-a-multi-rule")
		if rule := resolveMultiRule(ctx); nil != rule {
			t.Fatalf("expected nil rule on type mismatch, got %+v", rule)
		}
	})

	t.Run("CTX_RULE 为 nil 时返回 nil", func(t *testing.T) {
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Set(constContextPg.CTX_RULE, nil)
		if rule := resolveMultiRule(ctx); nil != rule {
			t.Fatalf("expected nil rule, got %+v", rule)
		}
	})
}
