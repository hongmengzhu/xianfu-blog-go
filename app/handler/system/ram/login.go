package ram

import (
	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/domainRam/service"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/ram/modRamLogin"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/consts/constsRam/typeDomainPg"
	routerPg2 "github.com/hongmengzhu/xianfu-blog-go/pkg/core/routerPg"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(LoginController)).Name("SystemLoginController").Export(gs.As[routerPg2.RouteRegistrar]())
}

type LoginController struct {
	routerPg2.RouteRegistrar
	sv *service.AccountLoginService `autowire:"?"`
}

func (c *LoginController) RegisterRoutes(e *gin.Engine) {
	group := e.Group("/xianfu/auth/sys")
	group.POST("/login", c.Login)
	group.POST("/refresh", c.RefreshToken)
}

func (c *LoginController) Login(ctx *gin.Context) {
	var ct modRamLogin.LoginCt
	if !routerPg2.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.LoginBySystem(ctx, ct, typeDomainPg.System, typeDomainPg.System.Code()))
}

func (c *LoginController) RefreshToken(ctx *gin.Context) {
	var ct modRamLogin.TokenRefreshCt
	if !routerPg2.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, rg.Ok[string](""))
}
