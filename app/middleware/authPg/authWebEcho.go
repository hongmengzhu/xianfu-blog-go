package authPg

import (
	"context"

	"github.com/hongmengzhu/xianfu-blog-go/pkg/consts/constHeaderPg"
	"github.com/labstack/echo/v5"
	"go-spring.org/log"
)

// GroupWebMiddlewareEcho 前台租户中间件（echo 版本），
// 与 GroupWebMiddleware 保持相同的租户解析规则与上下文键。
//
//	@Description: 前台租户中间件
//	@param m
//	@return echo.MiddlewareFunc
func GroupWebMiddlewareEcho(m *GroupWebMiddlewareSp) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			val := m.ResolveTenant(c.Request().Host)
			log.Info(context.Background(), log.TagAppDef, "租户", log.Any("TenantNo", val))
			c.Set(constHeaderPg.WebTenantNo, val)
			c.Set(constHeaderPg.WebTemplatePg, m.pg)
			c.Set(constHeaderPg.WebTemplatePgServer, m.server)
			return next(c)
		}
	}
}
