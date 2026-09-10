package controller

import (
	"fmt"

	"github.com/gin-gonic/gin"
	modRamJobFunction2 "github.com/hongmengzhu/xianfu-blog-go/app/models/ram/modRamJobFunction"
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
	gs.Provide(new(JobFunctionController)).Name("SystemJobFunctionController").Export(gs.As[routerPg.RouteRegistrar]())
}

type JobFunctionController struct {
	routerPg.RouteRegistrar
	controllerPg.SpSystemAuth
	sv *service.RamJobFunctionService `autowire:"?"`
}

func (c *JobFunctionController) RegisterRoutes(e *gin.Engine) {
	group := e.Group("/xianfu/sys/ram/job-function", authPg.GroupSystemMiddleware(c.Sp))
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

func (c *JobFunctionController) CreateUpdate(ctx *gin.Context) {
	var ct modRamJobFunction2.CreateUpdateCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	if ct.ID.ToInt64() > 0 {
		ctx.JSON(200, c.sv.Update(ctx, ct))
	} else {
		ctx.JSON(200, c.sv.Create(ctx, ct))
	}
}

func (c *JobFunctionController) Delete(ctx *gin.Context) {
	var ct model.BaseIdsCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.LogicalDeletion(ctx, ct.Ids))
}

func (c *JobFunctionController) Recovery(ctx *gin.Context) {
	var ct model.BaseIdsCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.LogicalRecovery(ctx, ct.Ids))
}

func (c *JobFunctionController) PhysicalDeletion(ctx *gin.Context) {
	var ct model.BaseIdsCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.PhysicalDeletion(ctx, ct.Ids))
}

func (c *JobFunctionController) Detail(ctx *gin.Context) {
	param := ctx.Param("id")

	fmt.Println(param)
	if "" == param {
		ctx.JSON(200, rg.Error[string]("id不能为空"))
		return
	}
	ctx.JSON(200, c.sv.Detail(ctx, strPg.ToInt64(param)))
}

func (c *JobFunctionController) Enable(ctx *gin.Context) {
	var ct model.BaseIdsCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.Enable(ctx, ct))
}

func (c *JobFunctionController) Disable(ctx *gin.Context) {
	var ct model.BaseIdsCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.Disable(ctx, ct))
}

func (c *JobFunctionController) Query(ctx *gin.Context) {
	var ct modRamJobFunction2.QueryCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.Query(ctx, ct))
}

func (c *JobFunctionController) SelectNodeAll(ctx *gin.Context) {
	var ct modRamJobFunction2.QueryPublicCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.SelectNodeAll(ctx, ct))
}

func (c *JobFunctionController) SelectNodeAllPublic(ctx *gin.Context) {
	var ct modRamJobFunction2.QueryPublicCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ct.State = enumStatePg.ENABLE.IndexPg()
	ctx.JSON(200, c.sv.SelectNodeAllPublic(ctx, ct))
}

func (c *JobFunctionController) SelectPublic(ctx *gin.Context) {
	ct := modRamJobFunction2.QueryCt{State: enumStatePg.ENABLE.IndexPg()}
	ctx.JSON(200, c.sv.SelectPublic(ctx, ct))
}

func (c *JobFunctionController) ExistName(ctx *gin.Context) {
	var ct model.BaseExistWdCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.ExistName(ctx, ct))
}
