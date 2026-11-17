package ram

import (
	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/domainRam/service"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/ram/modRamLogin"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/consts/constHeaderPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/consts/constsRam/typeDomainPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/state/clientPg"
	routerPg "github.com/hongmengzhu/xianfu-blog-go/pkg/core/routerPg"
	"github.com/pangu-2/go-tools/tools/strPg"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(AuthLoginController)).Name("SystemAuthLoginController").Export(gs.As[routerPg.RouteRegistrar]())
}

type AuthLoginController struct {
	routerPg.RouteRegistrar
	sv *service.AccountLoginService `autowire:"?"`
}

func (c *AuthLoginController) RegisterRoutes(e *gin.Engine) {
	group := e.Group("/xianfu/auth/sys")
	group.POST("/login", c.Login)
	group.POST("/refresh", c.RefreshToken)
}

func (c *AuthLoginController) Login(ctx *gin.Context) {
	var ct modRamLogin.LoginCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.LoginBySystem(ctx, ct, typeDomainPg.System, typeDomainPg.System.Code()))
}

func (c *AuthLoginController) RefreshToken(ctx *gin.Context) {
	var ct modRamLogin.TokenRefreshCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	// 系统域令牌使用 authorization-sys 头，归一化到通用头后复用统一刷新逻辑
	if strPg.IsBlank(ctx.GetHeader(constHeaderPg.HeaderAuthorization)) {
		if sysToken := ctx.GetHeader(constHeaderPg.HeaderAuthorization_system); strPg.IsNotBlank(sysToken) {
			ctx.Request.Header.Set(constHeaderPg.HeaderAuthorization, sysToken)
		}
	}
	if strPg.IsBlank(ctx.GetHeader(constHeaderPg.HeaderAuthorizationRefresh)) {
		if sysRefresh := ctx.GetHeader(constHeaderPg.HeaderAuthorizationRefresh_system); strPg.IsNotBlank(sysRefresh) {
			ctx.Request.Header.Set(constHeaderPg.HeaderAuthorizationRefresh, sysRefresh)
		}
	}
	ctx.JSON(200, c.sv.RefreshToken(ctx, ct, typeDomainPg.System, clientPg.Browser))
}
