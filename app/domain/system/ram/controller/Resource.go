package controller

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/system/ram/service"
	"github.com/hongmengzhu/xianfu-blog-go/app/middleware/authPg"
	"github.com/hongmengzhu/xianfu-blog-go/app/middleware/authPg/controllerPg"
	modRamResource2 "github.com/hongmengzhu/xianfu-blog-go/app/models/ram/modRamResource"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/state/enumStatePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/model"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/routerPg"
	"github.com/pangu-2/go-tools/tools/strPg"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(ResourceController)).Name("SystemResourceController").Export(gs.As[routerPg.RouteRegistrar]())
}

type ResourceController struct {
	routerPg.RouteRegistrar
	controllerPg.SpSystemAuth
	sv *service.RamResourceService `autowire:"?"`
}

func (c *ResourceController) RegisterRoutes(e *gin.Engine) {
	group := e.Group("/xianfu/sys/ram/resource", authPg.GroupSystemMiddleware(c.Sp))
	group.POST("/createUpdate", c.CreateUpdate)
	group.POST("/createUpdateByCategory", c.CreateUpdateByCategory)
	group.GET("/detail/:id", c.Detail)
	group.POST("/enable", c.Enable)
	group.POST("/disable", c.Disable)
	group.POST("/delete", c.Delete)
	group.POST("/recovery", c.Recovery)
	group.POST("/physicalDeletion", c.PhysicalDeletion)
	group.POST("/query", c.Query)
	group.POST("/selectNodeAll", c.SelectNodeAll)
	group.POST("/selectNodeAllPublic", c.SelectNodeAllPublic)
	group.POST("/selectCategoryPublic", c.SelectCategoryPublic)
	group.POST("/selectCategory", c.SelectCategory)
	group.POST("/existName", c.ExistName)
}

func (c *ResourceController) CreateUpdate(ctx *gin.Context) {
	var ct modRamResource2.CreateUpdateCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	if ct.ID.ToInt64() > 0 {
		ctx.JSON(200, c.sv.Update(ctx, ct))
	} else {
		ctx.JSON(200, c.sv.Create(ctx, ct))
	}
}
func (c *ResourceController) CreateUpdateByCategory(ctx *gin.Context) {
	var ct modRamResource2.CreateUpdateCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ct.ParentNo = ""
	if ct.ID.ToInt64() > 0 {
		ctx.JSON(200, c.sv.Update(ctx, ct))
	} else {
		ctx.JSON(200, c.sv.Create(ctx, ct))
	}
}

func (c *ResourceController) Delete(ctx *gin.Context) {
	var ct model.BaseIdsCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.LogicalDeletion(ctx, ct.Ids))
}

func (c *ResourceController) Recovery(ctx *gin.Context) {
	var ct model.BaseIdsCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.LogicalRecovery(ctx, ct.Ids))
}

func (c *ResourceController) PhysicalDeletion(ctx *gin.Context) {
	var ct model.BaseIdsCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.PhysicalDeletion(ctx, ct.Ids))
}

func (c *ResourceController) Detail(ctx *gin.Context) {
	param := ctx.Param("id")
	fmt.Println(param)
	if "" == param {
		ctx.JSON(200, rg.Error[string]("id不能为空"))
		return
	}
	ctx.JSON(200, c.sv.Detail(ctx, strPg.ToInt64(param)))
}

func (c *ResourceController) Enable(ctx *gin.Context) {
	var ct model.BaseIdsCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.Enable(ctx, ct))
}

func (c *ResourceController) Disable(ctx *gin.Context) {
	var ct model.BaseIdsCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.Disable(ctx, ct))
}

func (c *ResourceController) Query(ctx *gin.Context) {
	var ct modRamResource2.QueryCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.Query(ctx, ct))
}

func (c *ResourceController) SelectNodeAll(ctx *gin.Context) {
	var ct modRamResource2.QueryPublicCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.SelectNodeAll(ctx, ct))
}

func (c *ResourceController) SelectNodeAllPublic(ctx *gin.Context) {
	var ct modRamResource2.QueryPublicCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ct.State = enumStatePg.ENABLE.IndexPg()
	ctx.JSON(200, c.sv.SelectNodeAllPublic(ctx, ct))
}

func (c *ResourceController) SelectPublic(ctx *gin.Context) {
	ct := modRamResource2.QueryCt{State: enumStatePg.ENABLE.IndexPg()}
	ctx.JSON(200, c.sv.SelectPublic(ctx, ct))
}

func (c *ResourceController) SelectCategoryPublic(ctx *gin.Context) {
	ct := modRamResource2.QueryCt{State: enumStatePg.ENABLE.IndexPg()}
	ctx.JSON(200, c.sv.SelectCategoryPublic(ctx, ct))
}
func (c *ResourceController) SelectCategory(ctx *gin.Context) {
	ct := modRamResource2.QueryCt{}
	ctx.JSON(200, c.sv.SelectCategoryPublic(ctx, ct))
}

func (c *ResourceController) ExistName(ctx *gin.Context) {
	var ct model.BaseExistWdCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.ExistName(ctx, ct))
}
