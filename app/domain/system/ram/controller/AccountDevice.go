package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/system/ram/service"
	"github.com/hongmengzhu/xianfu-blog-go/app/middleware/authPg"
	"github.com/hongmengzhu/xianfu-blog-go/app/middleware/authPg/controllerPg"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/ram/modRamAccountDevice"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/model"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/routerPg"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(AccountDeviceController)).Name("SystemAccountDeviceController").Export(gs.As[routerPg.RouteRegistrar]())
}

type AccountDeviceController struct {
	routerPg.RouteRegistrar
	controllerPg.SpSystemAuth
	sv *service.RamAccountDeviceService `autowire:"?"`
}

func (c *AccountDeviceController) RegisterRoutes(e *gin.Engine) {
	group := e.Group("/xianfu/sys/ram/account-device", authPg.GroupSystemMiddleware(c.Sp))
	group.POST("/physicalDeletion", c.PhysicalDeletion)
	group.POST("/query", c.Query)
}

func (c *AccountDeviceController) PhysicalDeletion(ctx *gin.Context) {
	var ct model.BaseIdsCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.PhysicalDeletion(ctx, ct.Ids))
}

func (c *AccountDeviceController) Query(ctx *gin.Context) {
	var ct modRamAccountDevice.QueryCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.Query(ctx, ct))
}
