package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/system/ram/service"
	"github.com/hongmengzhu/xianfu-blog-go/app/middleware/authPg"
	"github.com/hongmengzhu/xianfu-blog-go/app/middleware/authPg/controllerPg"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/ram/modRamAccountLoginLog"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/model"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/routerPg"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(AccountLoginLogController)).Name("SystemAccountLoginLogController").Export(gs.As[routerPg.RouteRegistrar]())
}

type AccountLoginLogController struct {
	routerPg.RouteRegistrar
	controllerPg.SpSystemAuth
	sv *service.RamAccountLoginLogService `autowire:"?"`
}

func (c *AccountLoginLogController) RegisterRoutes(e *gin.Engine) {
	group := e.Group("/xianfu/sys/ram/account-login-log", authPg.GroupSystemMiddleware(c.Sp))
	group.POST("/physicalDeletion", c.PhysicalDeletion)
	group.POST("/query", c.Query)
}

func (c *AccountLoginLogController) PhysicalDeletion(ctx *gin.Context) {
	var ct model.BaseIdsCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.PhysicalDeletion(ctx, ct.Ids))
}

func (c *AccountLoginLogController) Query(ctx *gin.Context) {
	var ct modRamAccountLoginLog.QueryCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.Query(ctx, ct))
}
