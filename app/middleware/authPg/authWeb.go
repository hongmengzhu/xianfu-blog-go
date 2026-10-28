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

// ResolveTenant 按访问域名解析租户编号，gin 与 echo 中间件共用。
//
//	@Description: 解析租户
//	@receiver m
//	@param host 请求域名
//	@return string 租户编号，未匹配到返回 -1
func (m *GroupWebMiddlewareSp) ResolveTenant(host string) string {
	val := "-1"
	// 获取租户
	load, b := m.domain.Domain.Load(host)
	if b {
		val = load
	}
	//本地域名
	if m.domain.IsLocalHostExist(host) {
		val = "1000"
	}
	// 指定域名
	if m.domain.IsServerHostExist(host, m.server) {
		val = "1000"
	}
	return val
}

// 权限验证 中间件
func GroupWebMiddleware(m *GroupWebMiddlewareSp) gin.HandlerFunc {
	return func(c *gin.Context) {
		val := m.ResolveTenant(c.Request.Host)
		log.Info(context.Background(), log.TagAppDef, "租户", log.Any("TenantNo", val))
		c.Set(constHeaderPg.WebTenantNo, val)
		c.Set(constHeaderPg.WebTemplatePg, m.pg)
		c.Set(constHeaderPg.WebTemplatePgServer, m.server)
		c.Next()
	}
}
