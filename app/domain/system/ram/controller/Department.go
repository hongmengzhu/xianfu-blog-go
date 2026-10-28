package controller

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/system/ram/service"
	"github.com/hongmengzhu/xianfu-blog-go/app/middleware/authPg"
	"github.com/hongmengzhu/xianfu-blog-go/app/middleware/authPg/controllerPg"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/ram/modRamDepartment"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/state/enumStatePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/model"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/routerPg"
	"github.com/pangu-2/go-tools/tools/strPg"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(DepartmentController)).Name("SystemDepartmentController").Export(gs.As[routerPg.RouteRegistrar]())
}

// DepartmentController 部门
// @Description:
type DepartmentController struct {
	routerPg.RouteRegistrar
	controllerPg.SpSystemAuth
	sv *service.RamDepartmentService `autowire:"?"`
}

// RegisterRoutes 注册路由
//
//	@Description:
//	@receiver c
//	@param e
func (c *DepartmentController) RegisterRoutes(e *gin.Engine) {
	group := e.Group("/xianfu/sys/ram/department", authPg.GroupSystemMiddleware(c.Sp))
	group.POST("/createUpdate", c.CreateUpdate)
	group.GET("/detail/:id", c.Detail)
	group.POST("/enable", c.Enable)
	group.POST("/disable", c.Disable)
	group.POST("/delete", c.Delete)
	group.POST("/recovery", c.Recovery)
	group.POST("/physicalDeletion", c.PhysicalDeletion)
	group.POST("/query", c.Query)
	group.POST("/queryAll", c.QueryAll)
	group.POST("/selectNodeAll", c.SelectNodeAll)
	group.POST("/selectNodeAllPublic", c.SelectNodeAllPublic)
	group.POST("/exportExcel", c.ExportExcel)
	group.POST("/existName", c.ExistName)
}

func (c *DepartmentController) CreateUpdate(ctx *gin.Context) {
	var ct modRamDepartment.CreateUpdateCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	if ct.ID.ToInt64() > 0 {
		ctx.JSON(200, c.sv.Update(ctx, ct))
	} else {
		ctx.JSON(200, c.sv.Create(ctx, ct))
	}
}

func (c *DepartmentController) Delete(ctx *gin.Context) {
	var ct model.BaseIdsCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.LogicalDeletion(ctx, ct.Ids))
}

// Recovery 逻辑删除恢复
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *DepartmentController) Recovery(ctx *gin.Context) {
	var ct model.BaseIdsCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.LogicalRecovery(ctx, ct.Ids))
}

// PhysicalDeletion 物理删除
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *DepartmentController) PhysicalDeletion(ctx *gin.Context) {
	var ct model.BaseIdsCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.PhysicalDeletion(ctx, ct.Ids))
}

// Detail 详情
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *DepartmentController) Detail(ctx *gin.Context) {
	param := ctx.Param("id")

	fmt.Println(param)
	if "" == param {
		ctx.JSON(200, rg.Error[string]("id不能为空"))
		return
	}
	ctx.JSON(200, c.sv.Detail(ctx, strPg.ToInt64(param)))
}

// Enable 启用
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *DepartmentController) Enable(ctx *gin.Context) {
	var ct model.BaseIdsCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.Enable(ctx, ct))
}

// Disable 禁用
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *DepartmentController) Disable(ctx *gin.Context) {
	var ct model.BaseIdsCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.Disable(ctx, ct))
}

func (c *DepartmentController) Query(ctx *gin.Context) {
	var ct modRamDepartment.QueryCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.Query(ctx, ct))
}

func (c *DepartmentController) QueryAll(ctx *gin.Context) {
	var ct modRamDepartment.QueryCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.QueryAll(ctx, ct))
}

func (c *DepartmentController) SelectNodeAll(ctx *gin.Context) {
	var ct modRamDepartment.QueryPublicCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.SelectNodeAll(ctx, ct))
}

func (c *DepartmentController) SelectNodeAllPublic(ctx *gin.Context) {
	var ct modRamDepartment.QueryPublicCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ct.State = enumStatePg.ENABLE.IndexPg()
	ctx.JSON(200, c.sv.SelectNodeAllPublic(ctx, ct))
}

func (c *DepartmentController) ExportExcel(ctx *gin.Context) {
	ct := modRamDepartment.QueryCt{State: enumStatePg.ENABLE.IndexPg()}
	c.sv.ExportExcel(ctx, ct)
}

// ExistName 查重
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *DepartmentController) ExistName(ctx *gin.Context) {
	var ct model.BaseExistWdCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.ExistName(ctx, ct))
}
