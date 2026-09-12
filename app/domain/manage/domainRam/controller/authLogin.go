package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/manage/domainRam/service"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/ram/modRamLogin"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/consts/constsRam/typeDomainPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/state/clientPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/routerPg"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(AuthLoginController)).Name("ManageAuthLoginController").Export(gs.As[routerPg.RouteRegistrar]())
}

// AuthLoginController 登录
// @Description:
type AuthLoginController struct {
	routerPg.RouteRegistrar
	sv *service.AccountLoginService `autowire:"?"`
}

// RegisterRoutes 注册路由
//
//	@Description:
//	@receiver c
//	@param e
func (c *AuthLoginController) RegisterRoutes(e *gin.Engine) {
	group := e.Group("/xianfu/auth/manage")
	group.POST("/login", c.Login)
	group.POST("/refresh", c.RefreshToken)
}

// Login 登陆
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *AuthLoginController) Login(ctx *gin.Context) {
	var ct modRamLogin.LoginManageCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.Login(ctx, ct, typeDomainPg.Manage))
}

// RefreshToken 刷新
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *AuthLoginController) RefreshToken(ctx *gin.Context) {
	var ct modRamLogin.TokenRefreshCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.RefreshToken(ctx, ct, typeDomainPg.Manage, clientPg.Browser))
}
