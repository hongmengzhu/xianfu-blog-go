package controller

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/manage/domainBasic/service"
	"github.com/hongmengzhu/xianfu-blog-go/app/middleware/authPg"
	"github.com/hongmengzhu/xianfu-blog-go/app/middleware/validatorPg"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/basic/modBasicModelRules"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/state/enumStatePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/model"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/routerPg"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(ModelRulesController)).Name("ManageModelRulesController").Export(gs.As[routerPg.RouteRegistrar]())
}

// ModelRulesController 国家
// @Description:
type ModelRulesController struct {
	routerPg.RouteRegistrar
	Sp *authPg.GroupManageMiddlewareSp `autowire:""`
	sv *service.BasicModelRulesService `autowire:"?"`
}

// RegisterRoutes 注册路由
//
//	@Description:
//	@receiver c
//	@param e
func (c *ModelRulesController) RegisterRoutes(e *gin.Engine) {
	group := e.Group("/xianfu/manage/basic/modelRules", authPg.GroupManageMiddleware(c.Sp))
	group.POST("/createUpdate", c.CreateUpdate)
	group.POST("/createUpdateData", c.CreateUpdateData)
	group.GET("/detail/:id", c.Detail)
	group.POST("/enable", c.Enable)
	group.POST("/disable", c.Disable)
	group.POST("/state", c.State)
	group.POST("/delete", c.Delete)
	group.POST("/recovery", c.Recovery)
	group.POST("/physicalDeletion", c.PhysicalDeletion)
	group.POST("/query", c.Query)
	group.POST("/selectNodeAllPublic", c.SelectNodeAllPublic)
	group.POST("/existName", c.ExistName)
	group.POST("/existCode", c.ExistCode)
	group.POST("/allByValueNo", c.AllByValueNo)
}

// CreateUpdateData 创建
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *ModelRulesController) CreateUpdateData(ctx *gin.Context) {
	var ct modBasicModelRules.CreateUpdateDataCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.CreateUpdateData(ctx, ct))
}

// CreateUpdate 创建
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *ModelRulesController) CreateUpdate(ctx *gin.Context) {
	var ct modBasicModelRules.CreateUpdateCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.CreateUpdate(ctx, ct))
}

// Delete 逻辑删除
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *ModelRulesController) Delete(ctx *gin.Context) {
	var ct model.BaseIdsCt[string]
	if err := ctx.ShouldBind(&ct); err != nil {
		if jsonErr, ok := err.(*json.UnmarshalTypeError); ok {
			// 处理 JSON 解析错误（类型不匹配）
			errMsg := fmt.Sprintf("参数解析错误：字段 %s 期望类型 %s，但实际传入类型 %s",
				jsonErr.Field, jsonErr.Type.String(), jsonErr.Value)
			ctx.JSON(200, rg.ErrorMessageData[string](errMsg))
			return
		}
		//对 返回 错误进行转义 成中文
		translate := validatorPg.Translate(err, &ct)
		if len(translate) > 0 {
			ctx.JSON(200, rg.ErrorMessageData[string](translate))
			return
		}
		ctx.JSON(200, rg.ErrorDefault[string]())
		return
	}
	ctx.JSON(200, c.sv.LogicalDeletion(ctx, ct))
}

// Recovery 逻辑删除恢复
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *ModelRulesController) Recovery(ctx *gin.Context) {
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
func (c *ModelRulesController) PhysicalDeletion(ctx *gin.Context) {
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
func (c *ModelRulesController) Detail(ctx *gin.Context) {
	param := ctx.Param("id")

	fmt.Println(param)
	if "" == param {
		ctx.JSON(200, rg.Error[string]("id不能为空"))
		return
	}
	//ctx.JSON(200, c.sv.Detail(ctx, strPg.ToInt64(param)))
}

// Enable 启用
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *ModelRulesController) Enable(ctx *gin.Context) {
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
func (c *ModelRulesController) Disable(ctx *gin.Context) {
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
func (c *ModelRulesController) State(ctx *gin.Context) {
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
func (c *ModelRulesController) Query(ctx *gin.Context) {
	var ct modBasicModelRules.QueryCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.Query(ctx, ct))
}

func (c *ModelRulesController) AllByValueNo(ctx *gin.Context) {
	var ct modBasicModelRules.QueryPublicCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ct.State = enumStatePg.ENABLE.IndexPg()
	ctx.JSON(200, c.sv.AllByValueNo(ctx, ct))
}

func (c *ModelRulesController) SelectNodeAllPublic(ctx *gin.Context) {
	var ct modBasicModelRules.QueryPublicCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ct.State = enumStatePg.ENABLE.IndexPg()
	ctx.JSON(200, c.sv.SelectNodeAllPublic(ctx, ct))
}

// ExistName 查重
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *ModelRulesController) ExistName(ctx *gin.Context) {
	var ct model.BaseExistWdCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.ExistName(ctx, ct))
}

// ExistCode 查重
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *ModelRulesController) ExistCode(ctx *gin.Context) {
	var ct model.BaseExistWdCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.ExistCode(ctx, ct))
}
