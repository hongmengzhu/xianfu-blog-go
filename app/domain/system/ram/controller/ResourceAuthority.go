package controller

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/system/ram/service"
	"github.com/hongmengzhu/xianfu-blog-go/app/middleware/authPg"
	"github.com/hongmengzhu/xianfu-blog-go/app/middleware/authPg/controllerPg"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/ram/modRamResourceAuthority"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/consts/constsRam/resourceTypeCategoryPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/state/enumStatePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/model"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/routerPg"
	"github.com/pangu-2/go-tools/tools/strPg"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(ResourceAuthorityController)).Name("SystemResourceAuthorityController").Export(gs.As[routerPg.RouteRegistrar]())
}

type ResourceAuthorityController struct {
	routerPg.RouteRegistrar
	controllerPg.SpSystemAuth
	sv *service.RamResourceAuthorityService `autowire:"?"`
}

func (c *ResourceAuthorityController) RegisterRoutes(e *gin.Engine) {
	group := e.Group("/xianfu/sys/ram/resource-authority", authPg.GroupSystemMiddleware(c.Sp))
	group.POST("/createByGroup", c.CreatByGroup)
	group.POST("/updateByRole", c.UpdateByRole)
	group.GET("/detail/:id", c.Detail)
	group.POST("/enable", c.Enable)
	group.POST("/disable", c.Disable)
	group.POST("/delete", c.Delete)
	group.POST("/recovery", c.Recovery)
	group.POST("/physicalDeletion", c.PhysicalDeletion)
	group.POST("/query", c.Query)
	group.POST("/queryByGroup", c.QueryByGroup)
	group.POST("/selectNodeAll", c.SelectNodePublic)
	group.POST("/selectNodeAllPublic", c.SelectNodeAllPublic)
}

func (c *ResourceAuthorityController) Create(ctx *gin.Context) {
	var ct modRamResourceAuthority.CreateCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.Create(ctx, ct))
}

func (c *ResourceAuthorityController) CreatByGroup(ctx *gin.Context) {
	var ct modRamResourceAuthority.CreatByGroupCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.CreatByGroup(ctx, ct))
}

func (c *ResourceAuthorityController) Update(ctx *gin.Context) {
	var ct modRamResourceAuthority.UpdateCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.Update(ctx, ct))
}

func (c *ResourceAuthorityController) UpdateByRole(ctx *gin.Context) {
	var ct modRamResourceAuthority.UpdateByTypeValueCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.UpdateByRole(ctx, ct))
}

func (c *ResourceAuthorityController) Delete(ctx *gin.Context) {
	var ct model.BaseIdsCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.LogicalDeletion(ctx, ct.Ids))
}

func (c *ResourceAuthorityController) Recovery(ctx *gin.Context) {
	var ct model.BaseIdsCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.LogicalRecovery(ctx, ct.Ids))
}

func (c *ResourceAuthorityController) PhysicalDeletion(ctx *gin.Context) {
	var ct model.BaseIdsCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.PhysicalDeletion(ctx, ct.Ids))
}

func (c *ResourceAuthorityController) Detail(ctx *gin.Context) {
	param := ctx.Param("id")

	fmt.Println(param)
	if "" == param {
		ctx.JSON(200, rg.Error[string]("id不能为空"))
		return
	}
	ctx.JSON(200, c.sv.Detail(ctx, strPg.ToInt64(param)))
}

func (c *ResourceAuthorityController) Enable(ctx *gin.Context) {
	var ct model.BaseIdsCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.Enable(ctx, ct))
}

func (c *ResourceAuthorityController) Disable(ctx *gin.Context) {
	var ct model.BaseIdsCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.Disable(ctx, ct))
}

func (c *ResourceAuthorityController) Query(ctx *gin.Context) {
	var ct modRamResourceAuthority.QueryCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.Query(ctx, ct))
}

func (c *ResourceAuthorityController) QueryByGroup(ctx *gin.Context) {
	var ct modRamResourceAuthority.QueryCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ct.TypeCategory = resourceTypeCategoryPg.Group.String()
	if strPg.IsBlank(ct.TypeValue) {
		ct.TypeValue = "0"
	}
	ctx.JSON(200, c.sv.Query(ctx, ct))
}

func (c *ResourceAuthorityController) SelectNodePublic(ctx *gin.Context) {
	ct := modRamResourceAuthority.QueryCt{State: enumStatePg.ENABLE.IndexPg()}
	ctx.JSON(200, c.sv.SelectNodePublic(ctx, ct))
}

func (c *ResourceAuthorityController) SelectNodeAllPublic(ctx *gin.Context) {
	ct := modRamResourceAuthority.QueryCt{State: enumStatePg.ENABLE.IndexPg()}
	ctx.JSON(200, c.sv.SelectNodeAllPublic(ctx, ct))
}

func (c *ResourceAuthorityController) SelectPublic(ctx *gin.Context) {
	ct := modRamResourceAuthority.QueryCt{State: enumStatePg.ENABLE.IndexPg()}
	ctx.JSON(200, c.sv.SelectPublic(ctx, ct))
}
