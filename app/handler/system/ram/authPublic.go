package ram

import (
	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/domainRam/service"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/ram/modPublic"
	"github.com/hongmengzhu/xianfu-blog-go/cmd"
	routerPg2 "github.com/hongmengzhu/xianfu-blog-go/pkg/core/routerPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/middleware/authPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/middleware/authPg/controllerPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/sdk/auth/holderPg"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(AuthPublicController)).Name("SystemPublicController").Export(gs.As[routerPg2.RouteRegistrar]())
}

type AuthPublicController struct {
	routerPg2.RouteRegistrar
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
	if !routerPg2.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.UpdatePassword(ctx, ct))
}
