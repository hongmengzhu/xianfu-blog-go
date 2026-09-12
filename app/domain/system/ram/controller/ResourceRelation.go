package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/system/ram/service"
	"github.com/hongmengzhu/xianfu-blog-go/app/middleware/authPg"
	"github.com/hongmengzhu/xianfu-blog-go/app/middleware/authPg/controllerPg"
	modRamResourceRelation2 "github.com/hongmengzhu/xianfu-blog-go/app/models/ram/modRamResourceRelation"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/state/enumStatePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/routerPg"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(ResourceRelationController)).Name("SystemResourceRelationController").Export(gs.As[routerPg.RouteRegistrar]())
}

type ResourceRelationController struct {
	routerPg.RouteRegistrar
	controllerPg.SpSystemAuth
	sv *service.RamResourceRelationService `autowire:"?"`
}

func (c *ResourceRelationController) RegisterRoutes(e *gin.Engine) {
	group := e.Group("/xianfu/sys/ram/resource-relation", authPg.GroupSystemMiddleware(c.Sp))
	group.POST("/selected", c.Selected)
}

func (c *ResourceRelationController) Query(ctx *gin.Context) {
	var ct modRamResourceRelation2.QueryCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.Query(ctx, ct))
}

func (c *ResourceRelationController) SelectNodePublic(ctx *gin.Context) {
	ct := modRamResourceRelation2.QueryCt{State: enumStatePg.ENABLE.IndexPg()}
	ctx.JSON(200, c.sv.SelectNodePublic(ctx, ct))
}

func (c *ResourceRelationController) SelectNodeAllPublic(ctx *gin.Context) {
	ct := modRamResourceRelation2.QueryCt{State: enumStatePg.ENABLE.IndexPg()}
	ctx.JSON(200, c.sv.SelectNodeAllPublic(ctx, ct))
}

func (c *ResourceRelationController) SelectPublic(ctx *gin.Context) {
	ct := modRamResourceRelation2.QueryCt{State: enumStatePg.ENABLE.IndexPg()}
	ctx.JSON(200, c.sv.SelectPublic(ctx, ct))
}

func (c *ResourceRelationController) Selected(ctx *gin.Context) {
	var ct modRamResourceRelation2.QuerySelectedCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.Selected(ctx, ct.Code))
}
