package authPg

import (
	"context"

	"github.com/hongmengzhu/xianfu-blog-go/app/middleware/components/cachePg/cacheDomainPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/configPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/consts/constHeaderPg"
	"github.com/labstack/echo/v5"
	"github.com/pangu-2/go-tools/tools/strPg"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(&GroupWebVhostMiddlewareSp{})
}

// 中间件 服务
type GroupWebVhostMiddlewareSp struct {
	Domain *cacheDomainPg.TenantDomainCache `autowire:"?"`
	pg     configPg.Pg                      `value:"${pg}"`
	server configPg.Server                  `value:"${server}"`
}

// ResolveTenant 按访问域名解析租户编号，gin 与 echo 中间件共用。
//
//	@Description: 解析租户
//	@receiver m
//	@param host 请求域名
//	@return string 租户编号，未匹配到返回 -1
func (m *GroupWebVhostMiddlewareSp) ResolveTenant(host string) string {
	val := "-1"
	// 获取租户
	load, b := m.Domain.Domain.Load(host)
	//getMap, b2 := m.Domain.GetMap(host)
	//log.Debug(context.Background(), log.TagAppDef, "ResolveTenant", func() []log.Field {
	//	return []log.Field{
	//		log.Any("host", host),
	//		log.Any("bool", b),
	//		log.Any("data", load),
	//		log.Any("getMap", getMap),
	//		log.Any("b2", b2),
	//	}
	//})
	if b {
		val = load
	}
	//本地域名
	if m.Domain.IsLocalHostExist(host) {
		val = "1000"
	}
	// 指定域名
	if m.Domain.IsServerHostExist(host, m.server) {
		val = "1000"
	}
	return val
}

// GroupWebVhostMiddlewareEcho 前台租户中间件（echo 版本），
// 与 GroupWebMiddleware 保持相同的租户解析规则与上下文键。
//
//	@Description: 前台租户中间件
//	@param m
//	@return echo.MiddlewareFunc
func GroupWebVhostMiddlewareEcho(m *GroupWebVhostMiddlewareSp) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			if strPg.IsNotBlank(c.Request().Host) {
				val := m.ResolveTenant(c.Request().Host)
				log.Info(context.Background(), log.TagAppDef, "租户", log.Any("host", c.Request().Host), log.Any("TenantNo", val))
				c.Set(constHeaderPg.WebTenantNo, val)
				c.Set(constHeaderPg.WebVhostTemplatePg, m.pg)
				c.Set(constHeaderPg.WebVhostTemplatePgServer, m.server)
				c.Set(constHeaderPg.WebVhostTemplatePgDomain, m.Domain)
				if get, b := m.Domain.GetMap(c.Request().Host); b {
					c.Set(constHeaderPg.WebVhostTemplatePgDomainData, get)
				}
			}
			return next(c)
		}
	}
}
