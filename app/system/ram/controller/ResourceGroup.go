package controller

import (
	"fmt"

	"github.com/gin-gonic/gin"
	modRamResourceGroup2 "github.com/hongmengzhu/xianfu-blog-go/app/models/ram/modRamResourceGroup"
	modRamResourceRelation2 "github.com/hongmengzhu/xianfu-blog-go/app/models/ram/modRamResourceRelation"
	"github.com/hongmengzhu/xianfu-blog-go/app/system/ram/service"
	"github.com/hongmengzhu/xianfu-blog-go/middleware/authPg"
	"github.com/hongmengzhu/xianfu-blog-go/middleware/authPg/controllerPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/state/enumStatePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/model"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/routerPg"
	"github.com/pangu-2/go-tools/tools/strPg"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(ResourceGroupController)).Name("SystemResourceGroupController").Export(gs.As[routerPg.RouteRegistrar]())
}

type ResourceGroupController struct {
	routerPg.RouteRegistrar
	controllerPg.SpSystemAuth
	sv  *service.RamResourceGroupService              `autowire:"?"`
	sva *service.RamResourceGroupAuthorizationService `autowire:"?"`
	ra  *service.RamResourceAuthorizationService      `autowire:"?"`
	rr  *service.RamResourceRelationService           `autowire:"?"`
}

func (c *ResourceGroupController) RegisterRoutes(e *gin.Engine) {
	group := e.Group("/xianfu/sys/ram/resource-group", authPg.GroupSystemMiddleware(c.Sp))
	group.POST("/createUpdate", c.CreateUpdate)
	group.POST("/createUpdateCategory", c.CreateUpdateCategory)
	group.GET("/detail/:id", c.Detail)
	group.POST("/enable", c.Enable)
	group.POST("/disable", c.Disable)
	group.POST("/delete", c.Delete)
	group.POST("/recovery", c.Recovery)
	group.POST("/physicalDeletion", c.PhysicalDeletion)
	group.POST("/query", c.Query)
	group.POST("/selectNodeAll", c.SelectNodeAll)
	group.POST("/selectNodeAllPublic", c.SelectNodeAllPublic)
	group.POST("/updateByResourceGroup", c.UpdateByResourceGroup)
	group.POST("/resourceSelected", c.Selected)
	group.POST("/existName", c.ExistName)
	group.POST("/selectCategory", c.SelectCategory)
	group.POST("/queryAllCategory", c.QueryAllCategory)
}

func (c *ResourceGroupController) CreateUpdate(ctx *gin.Context) {
	var ct modRamResourceGroup2.CreateUpdateCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	if ct.ID.ToInt64() > 0 {
		ctx.JSON(200, c.sv.Update(ctx, ct))
	} else {
		ctx.JSON(200, c.sv.Create(ctx, ct))
	}
}
func (c *ResourceGroupController) CreateUpdateCategory(ctx *gin.Context) {
	var ct modRamResourceGroup2.CreateUpdateCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	if ct.ID.ToInt64() > 0 {
		ctx.JSON(200, c.sv.UpdateCategory(ctx, ct))
	} else {
		ctx.JSON(200, c.sv.CreateCategory(ctx, ct))
	}
}

func (c *ResourceGroupController) Delete(ctx *gin.Context) {
	var ct model.BaseIdsCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sva.LogicalDeletion(ctx, ct.Ids))
}

func (c *ResourceGroupController) Recovery(ctx *gin.Context) {
	var ct model.BaseIdsCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sva.LogicalRecovery(ctx, ct.Ids))
}

func (c *ResourceGroupController) PhysicalDeletion(ctx *gin.Context) {
	var ct model.BaseIdsCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sva.PhysicalDeletion(ctx, ct.Ids))
}

func (c *ResourceGroupController) Detail(ctx *gin.Context) {
	param := ctx.Param("id")

	fmt.Println(param)
	if "" == param {
		ctx.JSON(200, rg.Error[string]("id不能为空"))
		return
	}
	ctx.JSON(200, c.sv.Detail(ctx, strPg.ToInt64(param)))
}

func (c *ResourceGroupController) Enable(ctx *gin.Context) {
	var ct model.BaseIdsCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sva.Enable(ctx, ct))
}

func (c *ResourceGroupController) Disable(ctx *gin.Context) {
	var ct model.BaseIdsCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sva.Disable(ctx, ct))
}

func (c *ResourceGroupController) Query(ctx *gin.Context) {
	var ct modRamResourceGroup2.QueryCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.Query(ctx, ct))
}

func (c *ResourceGroupController) SelectNodeAll(ctx *gin.Context) {
	var ct modRamResourceGroup2.QueryPublicCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.SelectNodeAll(ctx, ct))
}

func (c *ResourceGroupController) SelectNodeAllPublic(ctx *gin.Context) {
	var ct modRamResourceGroup2.QueryPublicCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ct.State = enumStatePg.ENABLE.IndexPg()
	ctx.JSON(200, c.sv.SelectNodeAllPublic(ctx, ct))
}

func (c *ResourceGroupController) SelectCategory(ctx *gin.Context) {
	ct := modRamResourceGroup2.QueryCt{State: enumStatePg.ENABLE.IndexPg()}
	ctx.JSON(200, c.sv.SelectCategory(ctx, ct))
}
func (c *ResourceGroupController) QueryAllCategory(ctx *gin.Context) {
	ct := modRamResourceGroup2.QueryCt{State: enumStatePg.ENABLE.IndexPg()}
	ctx.JSON(200, c.sv.QueryAllCategory(ctx, ct))
}

func (c *ResourceGroupController) UpdateByResourceGroup(ctx *gin.Context) {
	var ct modRamResourceRelation2.UpdateByResourceGroupCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.ra.UpdateByResourceGroup(ctx, ct))
}

func (c *ResourceGroupController) Selected(ctx *gin.Context) {
	var ct modRamResourceRelation2.QuerySelectedCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.rr.Selected(ctx, ct.Code))
}

func (c *ResourceGroupController) ExistName(ctx *gin.Context) {
	var ct model.BaseExistWdCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.ExistName(ctx, ct))
}
