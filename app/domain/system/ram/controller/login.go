package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/system/ram/service"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/ram/modRamLogin"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/consts/constsRam/typeDomainPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/routerPg"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(LoginController)).Name("SystemLoginController").Export(gs.As[routerPg.RouteRegistrar]())
}

type LoginController struct {
	routerPg.RouteRegistrar
	sv *service.AccountLoginService `autowire:"?"`
}

func (c *LoginController) RegisterRoutes(e *gin.Engine) {
	group := e.Group("/xianfu/auth/sys")
	group.POST("/login", c.Login)
	group.POST("/refresh", c.RefreshToken)
}

func (c *LoginController) Login(ctx *gin.Context) {
	var ct modRamLogin.LoginCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.Login(ctx, ct, typeDomainPg.System, typeDomainPg.System.Code()))
}

func (c *LoginController) RefreshToken(ctx *gin.Context) {
	var ct modRamLogin.TokenRefreshCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.RefreshToken(ctx, ct))
}
