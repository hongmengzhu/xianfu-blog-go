package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/manage/domainApi/service"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/api/modApiDiplAccessKey"
	"github.com/hongmengzhu/xianfu-blog-go/middleware/authPg"
	"github.com/hongmengzhu/xianfu-blog-go/middleware/authPg/controllerPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/model"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/routerPg"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(DiplAccessKeyController)).Name("ManageDiplAccessKeyController").Export(gs.As[routerPg.RouteRegistrar]())
}

// DiplAccessKeyController 密钥
// @Description:
type DiplAccessKeyController struct {
	routerPg.RouteRegistrar
	controllerPg.SpManageAuth
	sv *service.ApiDiplAccessKeyService `autowire:"?"`
}

func (c *DiplAccessKeyController) RegisterRoutes(e *gin.Engine) {
	group := e.Group("/xianfu/manage/api/dipl-access-key", authPg.GroupManageMiddleware(c.Sp))
	group.POST("/enable", c.Enable)
	group.POST("/disable", c.Disable)
	group.POST("/state", c.State)
	group.POST("/delete", c.Delete)
	group.POST("/recovery", c.Recovery)
	group.POST("/physicalDeletion", c.PhysicalDeletion)
	group.POST("/query", c.Query)
	group.POST("/selectPublic", c.SelectPublic)
	group.POST("/makeNew", c.MakeNewRecord)
}

// MakeNewRecord 新记录
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *DiplAccessKeyController) MakeNewRecord(ctx *gin.Context) {
	var ct model.BaseIdCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.MakeNewRecord(ctx, ct))
}

// Delete 逻辑删除
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *DiplAccessKeyController) Delete(ctx *gin.Context) {
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
func (c *DiplAccessKeyController) Recovery(ctx *gin.Context) {
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
func (c *DiplAccessKeyController) PhysicalDeletion(ctx *gin.Context) {
	var ct model.BaseIdsCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.PhysicalDeletion(ctx, ct.Ids))
}

// Enable 启用
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *DiplAccessKeyController) Enable(ctx *gin.Context) {
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
func (c *DiplAccessKeyController) Disable(ctx *gin.Context) {
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
func (c *DiplAccessKeyController) State(ctx *gin.Context) {
	var ct model.BaseStateIdsCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.StateEnableDisable(ctx, ct))
}

// Query 查询列表
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *DiplAccessKeyController) Query(ctx *gin.Context) {
	var ct modApiDiplAccessKey.QueryCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.Query(ctx, ct))
}

func (c *DiplAccessKeyController) SelectPublic(ctx *gin.Context) {
	var ct modApiDiplAccessKey.QueryCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	//ct.State = enumStatePg.ENABLE.IndexPg()
	ctx.JSON(200, c.sv.SelectPublic(ctx, ct))
}
