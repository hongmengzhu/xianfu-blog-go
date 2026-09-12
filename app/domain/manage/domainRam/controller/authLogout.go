package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/manage/domainRam/service"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/auth/holderPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/routerPg"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(AuthLogoutController)).Name("ManageAuthLogoutController").Export(gs.As[routerPg.RouteRegistrar]())
}

// AuthLogoutController 退出
// @Description:
type AuthLogoutController struct {
	routerPg.RouteRegistrar
	sv *service.AccountLogoutService `autowire:"?"`
}

// RegisterRoutes 注册路由
//
//	@Description:
//	@receiver c
//	@param e
func (c *AuthLogoutController) RegisterRoutes(e *gin.Engine) {
	group := e.Group("/xianfu/auth/manage")
	group.Any("/logout", c.Logout)
}

// Logout 退出
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *AuthLogoutController) Logout(ctx *gin.Context) {
	ctx.JSON(200, c.sv.Logout(holderPg.GetContextAccount(ctx)))
}
