package authPg

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/tc/cacheTc"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/configPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/consts/constHeaderPg"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(&GroupWebMiddlewareSp{})
}

// 中间件 服务
type GroupWebMiddlewareSp struct {
	domain *cacheTc.TenantDomainCache `autowire:"?"`
	pg     configPg.Pg                `value:"${pg}"`
	server configPg.Server            `value:"${server}"`
}

// 权限验证 中间件
func GroupWebMiddleware(m *GroupWebMiddlewareSp) gin.HandlerFunc {
	return func(c *gin.Context) {
		val := "-1"
		// 获取租户
		load, b := m.domain.Domain.Load(c.Request.Host)
		if b {
			val = load
		}
		//本地域名
		if m.domain.IsLocalHostExist(c.Request.Host) {
			val = "1000"
		}
		// 指定域名
		if m.domain.IsServerHostExist(c.Request.Host, m.server) {
			val = "1000"
		}
		log.Infof(context.Background(), log.TagAppDef, "租户.no => %+v", val)
		c.Set(constHeaderPg.WebTenantNo, val)
		c.Set(constHeaderPg.WebTemplatePg, m.pg)
		c.Set(constHeaderPg.WebTemplatePgServer, m.server)
		c.Next()
	}
}
