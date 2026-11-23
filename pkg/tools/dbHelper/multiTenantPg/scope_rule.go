package multiTenantPg

import (
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/consts/constContextPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/sdk/auth/holderPg/holderIfPg"
	"go-spring.org/log"
)

var ruleTypeWarnOnce sync.Once

// resolveMultiRule 从 ctx 中取出 CTX_RULE 并解析为 *MultiRule。
// 未设置或类型不符时返回 nil，由调用方回退到默认规则；
// 类型不符属于异常情况，只告警一次，避免正常请求路径刷 ERROR 日志。
func resolveMultiRule(ctx *gin.Context) *MultiRule {
	value, exists := ctx.Get(constContextPg.CTX_RULE)
	if !exists || nil == value {
		return nil
	}
	if holderRule, ok := value.(holderIfPg.IHolderRule); ok {
		if rule, ok := holderRule.(*MultiRule); ok {
			return rule
		}
	}
	ruleTypeWarnOnce.Do(func() {
		log.Warn(ctx, log.TagAppDef, "CTX_RULE 类型异常，已回退默认多租户规则", log.Any("value", value))
	})
	return nil
}
