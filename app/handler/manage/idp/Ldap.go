package idp

import (
	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/domainIdp/service"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/idp/modIdpLdap"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/routerPg"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(LdapController)).Name("ManageIdpLdapController").Export(gs.As[routerPg.RouteRegistrar]())
}

// LdapController LDAP 管理
type LdapController struct {
	routerPg.RouteRegistrar
	sv *service.LdapService `autowire:"?"`
}

// RegisterRoutes 注册路由
func (c *LdapController) RegisterRoutes(e *gin.Engine) {
	group := e.Group("/xianfu/manage/ram-idp/ldap")
	group.POST("/test", c.TestConnection)
	group.POST("/search", c.SearchUsers)
	group.POST("/sync", c.SyncUsers)
	group.POST("/login", c.Login)
}

// TestConnection 测试 LDAP 连接
func (c *LdapController) TestConnection(ctx *gin.Context) {
	var ct modIdpLdap.LdapTestCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.TestConnection(ctx, ct))
}

// SearchUsers 搜索 LDAP 用户
func (c *LdapController) SearchUsers(ctx *gin.Context) {
	var ct modIdpLdap.LdapSearchCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.SearchUsers(ctx, ct))
}

// SyncUsers 同步 LDAP 用户
func (c *LdapController) SyncUsers(ctx *gin.Context) {
	var ct modIdpLdap.LdapSyncCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.SyncUsers(ctx, ct))
}

// Login LDAP 登录
func (c *LdapController) Login(ctx *gin.Context) {
	var ct modIdpLdap.LdapLoginCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.Login(ctx, ct))
}
