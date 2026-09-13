package controller

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/manage/domainBasic/service"
	"github.com/hongmengzhu/xianfu-blog-go/app/middleware/authPg"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/basic/modBasicAccountApplyDenyList"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/state/enumStatePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/model"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/routerPg"
	"github.com/pangu-2/go-tools/tools/strPg"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(AccountApplyDenyListController)).Name("ManageAccountApplyDenyListController").Export(gs.As[routerPg.RouteRegistrar]())
}

// AccountApplyDenyListController 账号申请黑名单
// @Description:
type AccountApplyDenyListController struct {
	routerPg.RouteRegistrar
	Sp *authPg.GroupManageMiddlewareSp           `autowire:"?"`
	sv *service.BasicAccountApplyDenyListService `autowire:"?"`
}

// RegisterRoutes 注册路由
//
//	@Description:
//	@receiver c
//	@param e
func (c *AccountApplyDenyListController) RegisterRoutes(e *gin.Engine) {
	group := e.Group("/xianfu/manage/basic/accountApplyDenyList", authPg.GroupManageMiddleware(c.Sp))
	group.GET("/detail/:id", c.Detail)
	group.POST("/query", c.Query)
	group.POST("/selectPublic", c.SelectPublic)
	group.POST("/selectNodePublic", c.SelectNodePublic)
	group.POST("/selectNodeAllPublic", c.SelectNodeAllPublic)
	group.POST("/exportExcel", c.ExportExcel)
	group.POST("/existName", c.ExistName)
	group.POST("/existExpr", c.ExistExpr)
}

// Create 创建
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *AccountApplyDenyListController) Create(ctx *gin.Context) {
	var ct modBasicAccountApplyDenyList.CreateUpdateCt
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
func (c *AccountApplyDenyListController) Update(ctx *gin.Context) {
	var ct modBasicAccountApplyDenyList.CreateUpdateCt
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
func (c *AccountApplyDenyListController) Delete(ctx *gin.Context) {
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
func (c *AccountApplyDenyListController) Recovery(ctx *gin.Context) {
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
func (c *AccountApplyDenyListController) PhysicalDeletion(ctx *gin.Context) {
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
func (c *AccountApplyDenyListController) Detail(ctx *gin.Context) {
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
func (c *AccountApplyDenyListController) Enable(ctx *gin.Context) {
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
func (c *AccountApplyDenyListController) Disable(ctx *gin.Context) {
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
func (c *AccountApplyDenyListController) State(ctx *gin.Context) {
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
func (c *AccountApplyDenyListController) Query(ctx *gin.Context) {
	var ct modBasicAccountApplyDenyList.QueryCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.Query(ctx, ct))
}

func (c *AccountApplyDenyListController) SelectPublic(ctx *gin.Context) {
	var ct modBasicAccountApplyDenyList.QueryPublicCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	if ct.State.ToInt8() <= 0 {
		ct.State = enumStatePg.ENABLE.IndexPg()
	}
	ctx.JSON(200, c.sv.SelectPublic(ctx, ct))
}

func (c *AccountApplyDenyListController) SelectNodePublic(ctx *gin.Context) {
	var ct modBasicAccountApplyDenyList.QueryPublicCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ct.State = enumStatePg.ENABLE.IndexPg()
	ctx.JSON(200, c.sv.SelectNodePublic(ctx, ct))
}
func (c *AccountApplyDenyListController) SelectNodeAllPublic(ctx *gin.Context) {
	var ct modBasicAccountApplyDenyList.QueryPublicCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ct.State = enumStatePg.ENABLE.IndexPg()
	ctx.JSON(200, c.sv.SelectNodeAllPublic(ctx, ct))
}

func (c *AccountApplyDenyListController) ExportExcel(ctx *gin.Context) {
	ct := modBasicAccountApplyDenyList.QueryCt{State: enumStatePg.ENABLE.IndexPg()}
	c.sv.ExportExcel(ctx, ct)
}

// ExistName 查重
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *AccountApplyDenyListController) ExistName(ctx *gin.Context) {
	var ct model.BaseExistWdCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.ExistName(ctx, ct))
}

// ExistExpr 查重
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *AccountApplyDenyListController) ExistExpr(ctx *gin.Context) {
	var ct model.BaseExistWdCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.ExistExpr(ctx, ct))
}
