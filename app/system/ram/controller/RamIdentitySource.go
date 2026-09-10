package controller

import (
	"fmt"

	"github.com/gin-gonic/gin"
	modRamIdentitySource2 "github.com/hongmengzhu/xianfu-blog-go/app/models/ram/modRamIdentitySource"
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
	gs.Provide(new(RamIdentitySourceController)).Name("SystemRamIdentitySourceController").Export(gs.As[routerPg.RouteRegistrar]())
}

type RamIdentitySourceController struct {
	routerPg.RouteRegistrar
	controllerPg.SpSystemAuth
	sv *service.RamIdentitySourceService `autowire:"?"`
}

func (c *RamIdentitySourceController) RegisterRoutes(e *gin.Engine) {
	group := e.Group("/xianfu/sys/ram/identitySource", authPg.GroupSystemMiddleware(c.Sp))
	group.POST("/createUpdate", c.CreateUpdate)
	group.GET("/detail/:id", c.Detail)
	group.POST("/enable", c.Enable)
	group.POST("/disable", c.Disable)
	group.POST("/delete", c.Delete)
	group.POST("/recovery", c.Recovery)
	group.POST("/physicalDeletion", c.PhysicalDeletion)
	group.POST("/query", c.Query)
	group.POST("/selectNodeAll", c.SelectNodeAll)
	group.POST("/selectNodeAllPublic", c.SelectNodeAllPublic)
	group.POST("/existName", c.ExistName)
}

func (c *RamIdentitySourceController) CreateUpdate(ctx *gin.Context) {
	var ct modRamIdentitySource2.CreateUpdateCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.CreateUpdate(ctx, ct))
}

func (c *RamIdentitySourceController) Delete(ctx *gin.Context) {
	var ct model.BaseIdsCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.LogicalDeletion(ctx, ct.Ids))
}

func (c *RamIdentitySourceController) Recovery(ctx *gin.Context) {
	var ct model.BaseIdsCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.LogicalRecovery(ctx, ct.Ids))
}

func (c *RamIdentitySourceController) PhysicalDeletion(ctx *gin.Context) {
	var ct model.BaseIdsCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.PhysicalDeletion(ctx, ct.Ids))
}

func (c *RamIdentitySourceController) Detail(ctx *gin.Context) {
	param := ctx.Param("id")

	fmt.Println(param)
	if "" == param {
		ctx.JSON(200, rg.Error[string]("id不能为空"))
		return
	}
	ctx.JSON(200, c.sv.Detail(ctx, strPg.ToInt64(param)))
}

func (c *RamIdentitySourceController) Enable(ctx *gin.Context) {
	var ct model.BaseIdsCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.Enable(ctx, ct))
}

func (c *RamIdentitySourceController) Disable(ctx *gin.Context) {
	var ct model.BaseIdsCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.Disable(ctx, ct))
}

func (c *RamIdentitySourceController) Query(ctx *gin.Context) {
	var ct modRamIdentitySource2.QueryCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.Query(ctx, ct))
}

func (c *RamIdentitySourceController) SelectNodeAll(ctx *gin.Context) {
	var ct modRamIdentitySource2.QueryPublicCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.SelectNodeAll(ctx, ct))
}

func (c *RamIdentitySourceController) SelectNodeAllPublic(ctx *gin.Context) {
	var ct modRamIdentitySource2.QueryPublicCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ct.State = enumStatePg.ENABLE.IndexPg()
	ctx.JSON(200, c.sv.SelectNodeAllPublic(ctx, ct))
}

func (c *RamIdentitySourceController) ExistName(ctx *gin.Context) {
	var ct model.BaseExistWdCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.ExistName(ctx, ct))
}
