package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/manage/domainRam/service"
	"github.com/hongmengzhu/xianfu-blog-go/app/middleware/authPg"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/ram/modRamResourceRelation"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/model"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/routerPg"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(ResourceRelationController)).Name("ManageResourceRelationController").Export(gs.As[routerPg.RouteRegistrar]())
}

// ResourceRelationController 资源关联
// @Description:
type ResourceRelationController struct {
	routerPg.RouteRegistrar
	Sp *authPg.GroupManageMiddlewareSp     `autowire:""`
	sv *service.RamResourceRelationService `autowire:"?"`
}

// RegisterRoutes 注册路由
//
//	@Description:
//	@receiver c
//	@param e
func (c *ResourceRelationController) RegisterRoutes(e *gin.Engine) {
	group := e.Group("/xianfu/manage/ram/resource-relation", authPg.GroupManageMiddleware(c.Sp))
	group.POST("/delete", c.Delete)
	group.POST("/physicalDeletion", c.PhysicalDeletion)
	group.POST("/query", c.Query)
}

// Delete 逻辑删除
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *ResourceRelationController) Delete(ctx *gin.Context) {
	var ct model.BaseIdsCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.LogicalDeletion(ctx, ct.Ids))
}

// PhysicalDeletion 物理删除
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *ResourceRelationController) PhysicalDeletion(ctx *gin.Context) {
	var ct model.BaseIdsCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.PhysicalDeletion(ctx, ct.Ids))
}

// Query 查询列表
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *ResourceRelationController) Query(ctx *gin.Context) {
	var ct modRamResourceRelation.QueryCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.Query(ctx, ct))
}
