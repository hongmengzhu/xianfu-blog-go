package controller

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/manage/domainRam/service"
	"github.com/hongmengzhu/xianfu-blog-go/app/middleware/authPg/controllerPg"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/ram/modRamChannel"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/state/enumStatePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/model"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/routerPg"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"

	"github.com/pangu-2/go-tools/tools/strPg"
)

func init() {

}

// ChannelController 渠道
// @Description:
type ChannelController struct {
	controllerPg.SpManageAuth
	sv *service.RamChannelService `autowire:"?"`
}

// Create 创建
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *ChannelController) Create(ctx *gin.Context) {
	var ct modRamChannel.CreateUpdateCt
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
func (c *ChannelController) Update(ctx *gin.Context) {
	var ct modRamChannel.CreateUpdateCt
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
func (c *ChannelController) Delete(ctx *gin.Context) {
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
func (c *ChannelController) Recovery(ctx *gin.Context) {
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
func (c *ChannelController) PhysicalDeletion(ctx *gin.Context) {
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
func (c *ChannelController) Detail(ctx *gin.Context) {
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
func (c *ChannelController) Enable(ctx *gin.Context) {
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
func (c *ChannelController) Disable(ctx *gin.Context) {
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
func (c *ChannelController) State(ctx *gin.Context) {
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
func (c *ChannelController) Query(ctx *gin.Context) {
	var ct modRamChannel.QueryCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.Query(ctx, ct))
}

func (c *ChannelController) SelectNodePublic(ctx *gin.Context) {
	var ct modRamChannel.QueryPublicCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ct.State = enumStatePg.ENABLE.IndexPg()
	ctx.JSON(200, c.sv.SelectNodePublic(ctx, ct))
}

func (c *ChannelController) SelectNodeAllPublic(ctx *gin.Context) {
	var ct modRamChannel.QueryPublicCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ct.State = enumStatePg.ENABLE.IndexPg()
	ctx.JSON(200, c.sv.SelectNodeAllPublic(ctx, ct))
}

func (c *ChannelController) SelectPublic(ctx *gin.Context) {
	ct := modRamChannel.QueryCt{State: enumStatePg.ENABLE.IndexPg()}
	ctx.JSON(200, c.sv.SelectPublic(ctx, ct))
}

// ExistName 查重
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *ChannelController) ExistName(ctx *gin.Context) {
	var ct model.BaseExistWdCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.ExistName(ctx, ct))
}
