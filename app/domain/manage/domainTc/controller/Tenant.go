package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/manage/domainTc/service"
	"github.com/hongmengzhu/xianfu-blog-go/app/middleware/authPg"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/tc/modTcTenant"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/state/enumStatePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/model"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/routerPg"
	"github.com/pangu-2/go-tools/tools/strPg"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(TenantController)).Name("ManageTenantController").Export(gs.As[routerPg.RouteRegistrar]())
}

// TenantController 租户
// @Description:
type TenantController struct {
	routerPg.RouteRegistrar
	Sp *authPg.GroupManageMiddlewareSp `autowire:""`
	sv *service.TcTenantService        `autowire:"?"`
}

// RegisterRoutes 注册路由
//
//	@Description:
//	@receiver c
//	@param e
func (c *TenantController) RegisterRoutes(e *gin.Engine) {
	group := e.Group("/xianfu/manage/tc/tenant", authPg.GroupManageMiddleware(c.Sp))
	group.GET("/detail/:id", c.Detail)
	group.POST("/query", c.Query)
	group.POST("/selectPublic", c.SelectPublic)
	group.POST("/selectNodePublic", c.SelectNodePublic)
	group.POST("/selectNodeAllPublic", c.SelectNodeAllPublic)
	group.POST("/existName", c.ExistName)
	group.POST("/existCode", c.ExistCode)
	group.POST("/existNo", c.ExistNo)
}

// Create 创建
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *TenantController) Create(ctx *gin.Context) {
	var ct modTcTenant.CreateCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.Create(ctx, ct))
}

// Update 更新
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *TenantController) Update(ctx *gin.Context) {
	var ct modTcTenant.UpdateCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.Update(ctx, ct))
}

// Delete 逻辑删除
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *TenantController) Delete(ctx *gin.Context) {
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
func (c *TenantController) Recovery(ctx *gin.Context) {
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
func (c *TenantController) PhysicalDeletion(ctx *gin.Context) {
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
func (c *TenantController) Detail(ctx *gin.Context) {
	param := ctx.Param("id")
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
func (c *TenantController) Enable(ctx *gin.Context) {
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
func (c *TenantController) Disable(ctx *gin.Context) {
	var ct model.BaseIdsCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.Disable(ctx, ct))
}

// State 状态
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *TenantController) State(ctx *gin.Context) {
	var ct model.BaseStateIdsCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	state, ok := enumStatePg.IsExistInt64(ct.State)
	if !ok {
		ctx.JSON(200, rg.Error[string]("类型不正确"))
		return
	}
	ctx.JSON(200, c.sv.StateEnableDisable(ctx, ct.Ids, state))
}

// Query 查询列表
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *TenantController) Query(ctx *gin.Context) {
	var ct modTcTenant.QueryCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.Query(ctx, ct))
}

// SelectNodePublic 选择节点公开
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *TenantController) SelectNodePublic(ctx *gin.Context) {
	ct := modTcTenant.QueryCt{State: enumStatePg.ENABLE.IndexPg()}
	ctx.JSON(200, c.sv.SelectNodePublic(ctx, ct))
}

// SelectNodeAllPublic 选择所有节点公开
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *TenantController) SelectNodeAllPublic(ctx *gin.Context) {
	ct := modTcTenant.QueryCt{State: enumStatePg.ENABLE.IndexPg()}
	ctx.JSON(200, c.sv.SelectNodeAllPublic(ctx, ct))
}

// SelectPublic 选择公开
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *TenantController) SelectPublic(ctx *gin.Context) {
	ct := modTcTenant.QueryCt{State: enumStatePg.ENABLE.IndexPg()}
	ctx.JSON(200, c.sv.SelectPublic(ctx, ct))
}

// ExistName 查重
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *TenantController) ExistName(ctx *gin.Context) {
	var ct model.BaseExistWdCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.ExistName(ctx, ct))
}

// ExistNo 查重
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *TenantController) ExistNo(ctx *gin.Context) {
	var ct model.BaseExistWdCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.ExistNo(ctx, ct))
}

// ExistCode 查重
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *TenantController) ExistCode(ctx *gin.Context) {
	var ct model.BaseExistWdCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.ExistCode(ctx, ct))
}
