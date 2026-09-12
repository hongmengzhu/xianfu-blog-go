package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/system/ram/service"
	"github.com/hongmengzhu/xianfu-blog-go/app/middleware/authPg"
	"github.com/hongmengzhu/xianfu-blog-go/app/middleware/authPg/controllerPg"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/ram/modPublic"
	"github.com/hongmengzhu/xianfu-blog-go/cmd"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/auth/holderPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/routerPg"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(AuthPublicController)).Name("SystemPublicController").Export(gs.As[routerPg.RouteRegistrar]())
}

type AuthPublicController struct {
	routerPg.RouteRegistrar
	controllerPg.SpSystemAuth
	sv *service.RamAccountPublicService `autowire:"?"`
}

func (c *AuthPublicController) RegisterRoutes(e *gin.Engine) {
	group := e.Group("/xianfu/sys/public", authPg.GroupSystemMiddleware(c.Sp))
	group.GET("/info", c.Public)
	group.GET("/infoPublic", c.InfoPublic)
	group.POST("/password", c.UpdatePassword)
	group.GET("/envInfoPublic", cmd.GetVersion)
}

func (c *AuthPublicController) Public(ctx *gin.Context) {
	ctx.JSON(200, c.sv.Public(ctx, holderPg.GetContextAccount(ctx)))
}
func (c *AuthPublicController) InfoPublic(ctx *gin.Context) {
	ctx.JSON(200, c.sv.InfoPublic(ctx, holderPg.GetContextAccount(ctx)))
}

func (c *AuthPublicController) UpdatePassword(ctx *gin.Context) {
	var ct modPublic.PasswordCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.UpdatePassword(ctx, ct))
}
