package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/system/ram/service"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/auth/holderPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/routerPg"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(LogoutController)).Name("SystemLogoutController").Export(gs.As[routerPg.RouteRegistrar]())
}

type LogoutController struct {
	routerPg.RouteRegistrar
	sv *service.AccountLogoutService `autowire:"?"`
}

func (c *LogoutController) RegisterRoutes(e *gin.Engine) {
	group := e.Group("/xianfu/auth/sys")
	group.Any("/logout", c.Logout)
}

func (c *LogoutController) Logout(ctx *gin.Context) {
	ctx.JSON(200, c.sv.Logout(holderPg.GetContextAccount(ctx)))
}
