package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/system/ram/service"
	"github.com/hongmengzhu/xianfu-blog-go/app/middleware/authPg"
	"github.com/hongmengzhu/xianfu-blog-go/app/middleware/authPg/controllerPg"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/ram/modRamResourceGroupRelation"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/consts/constsRam/resourceTypeCategoryPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/routerPg"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(ResourceGroupRelationController)).Name("SystemResourceGroupRelationController").Export(gs.As[routerPg.RouteRegistrar]())
}

type ResourceGroupRelationController struct {
	routerPg.RouteRegistrar
	controllerPg.SpSystemAuth
	sv *service.RamResourceGroupRelationService `autowire:"?"`
}

func (c *ResourceGroupRelationController) RegisterRoutes(e *gin.Engine) {
	group := e.Group("/xianfu/sys/ram/resource-group-relation", authPg.GroupSystemMiddleware(c.Sp))
	group.POST("/selectedByRole", c.SelectedByRole)
}

func (c *ResourceGroupRelationController) Selected(ctx *gin.Context) {
	var ct modRamResourceGroupRelation.QueryByTypeValueCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.Selected(ctx, ct))
}

// SelectedByRole
//
//	@Description: 角色
//	@receiver c
//	@param ctx
func (c *ResourceGroupRelationController) SelectedByRole(ctx *gin.Context) {
	var ct modRamResourceGroupRelation.QueryByTypeValueCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ct.TypeCategory = resourceTypeCategoryPg.Role.Index()
	ctx.JSON(200, c.sv.Selected(ctx, ct))
}
