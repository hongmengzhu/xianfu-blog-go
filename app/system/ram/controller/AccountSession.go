package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/ram/modRamAccountSession"
	"github.com/hongmengzhu/xianfu-blog-go/app/system/ram/service"
	"github.com/hongmengzhu/xianfu-blog-go/middleware/authPg"
	"github.com/hongmengzhu/xianfu-blog-go/middleware/authPg/controllerPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/model"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/routerPg"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(AccountSessionController)).Name("SystemAccountSessionController").Export(gs.As[routerPg.RouteRegistrar]())
}

type AccountSessionController struct {
	routerPg.RouteRegistrar
	controllerPg.SpSystemAuth
	sv *service.RamAccountSessionService `autowire:"?"`
}

func (c *AccountSessionController) RegisterRoutes(e *gin.Engine) {
	group := e.Group("/xianfu/sys/ram/account-session", authPg.GroupSystemMiddleware(c.Sp))
	group.POST("/physicalDeletion", c.PhysicalDeletion)
	group.POST("/query", c.Query)
}

func (c *AccountSessionController) PhysicalDeletion(ctx *gin.Context) {
	var ct model.BaseIdsCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.PhysicalDeletion(ctx, ct.Ids))
}

func (c *AccountSessionController) Query(ctx *gin.Context) {
	var ct modRamAccountSession.QueryCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.Query(ctx, ct))
}
