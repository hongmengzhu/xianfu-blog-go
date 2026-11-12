package ram

import (
	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/domainRam/service"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/routerPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/sdk/auth/holderPg"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(AuthLogoutController)).Name("SystemAuthLogoutController").Export(gs.As[routerPg.RouteRegistrar]())
}

type AuthLogoutController struct {
	routerPg.RouteRegistrar
	sv *service.AccountLogoutService `autowire:"?"`
}

func (c *AuthLogoutController) RegisterRoutes(e *gin.Engine) {
	group := e.Group("/xianfu/auth/sys")
	group.Any("/logout", c.Logout)
}

func (c *AuthLogoutController) Logout(ctx *gin.Context) {
	ctx.JSON(200, c.sv.Logout(holderPg.GetContextAccount(ctx)))
}
