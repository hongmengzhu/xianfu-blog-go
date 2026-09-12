package controller

import (
	"github.com/gin-gonic/gin"
	service2 "github.com/hongmengzhu/xianfu-blog-go/app/domain/system/tc/service"
	"github.com/hongmengzhu/xianfu-blog-go/app/middleware/authPg"
	"github.com/hongmengzhu/xianfu-blog-go/app/middleware/authPg/controllerPg"
	modRamAccount "github.com/hongmengzhu/xianfu-blog-go/app/models/tc/modTcAccount"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/enumCommonPg/appModulePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/state/enumStatePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/model"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/routerPg"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(TenantAccountController).SetAppModule(appModulePg.Manage)).
		Name("SystemTenantAccountController").
		Export(gs.As[routerPg.RouteRegistrar]())
}

// TenantAccountController 账户
// @Description:
type TenantAccountController struct {
	routerPg.RouteRegistrar
	controllerPg.SpSystemAuth
	sv        *service2.TcTenantAccountService         `autowire:"?"`
	ap        *service2.TcTenantAccountPasswordService `autowire:"?"`
	appModule appModulePg.AppModule
}

// SetAppModule 设置应用模块
//
//	@Description:
//	@receiver c
//	@param appModule
//	@return *TenantAccountController
func (c *TenantAccountController) SetAppModule(appModule appModulePg.AppModule) *TenantAccountController {
	c.appModule = appModule
	return c
}

// RegisterRoutes 注册路由
//
//	@Description:
//	@receiver c
//	@param e
func (c *TenantAccountController) RegisterRoutes(e *gin.Engine) {
	group := e.Group("/xianfu/sys/tc/tenant-account", authPg.GroupSystemMiddleware(c.Sp))
	group.POST("/enable", c.Enable)
	group.POST("/disable", c.Disable)
	group.GET("/detail/:id", c.Detail)
	group.POST("/delete", c.Delete)
	group.POST("/recovery", c.Recovery)
	group.POST("/physicalDeletion", c.PhysicalDeletion)
	group.POST("/query", c.Query)
	group.POST("/updatePassword", c.UpdatePassword)
	group.POST("/create", c.Create)
	group.POST("/createAccount", c.CreateAccount)
	group.POST("/update", c.Update)
	group.POST("/updateAccount", c.UpdateAccount)
	group.POST("/existAccount", c.ExistAccount)
	group.POST("/existPhone", c.ExistPhone)
	group.POST("/existMail", c.ExistMail)
	group.POST("/existIdentityCode", c.ExistIdentityCode)
	group.POST("/existRealName", c.ExistRealName)
}

// Detail 详情
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *TenantAccountController) Detail(ctx *gin.Context) {
	param := ctx.Param("id")
	if "" == param {
		ctx.JSON(200, rg.Error[string]("id不能为空"))
		return
	}
	ctx.JSON(200, c.sv.Detail(ctx, param, c.appModule))
}

// Enable 有效
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *TenantAccountController) Enable(ctx *gin.Context) {
	var ct model.BaseIdsCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.Enable(ctx, ct, c.appModule))
}

// Disable 停用
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *TenantAccountController) Disable(ctx *gin.Context) {
	var ct model.BaseIdsCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.Disable(ctx, ct, c.appModule))
}

// State 状态
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *TenantAccountController) State(ctx *gin.Context) {
	var ct model.BaseStateIdsCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	state, ok := enumStatePg.IsExistInt64(ct.State)
	if !ok {
		ctx.JSON(200, rg.Error[string]("类型不正确"))
		return
	}
	ctx.JSON(200, c.sv.StateEnableDisable(ctx, ct.Ids, state, c.appModule))
}

// Delete 逻辑删除
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *TenantAccountController) Delete(ctx *gin.Context) {
	var ct model.BaseIdsCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.LogicalDeletion(ctx, ct.Ids, c.appModule))
}

// Recovery 删除 恢复
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *TenantAccountController) Recovery(ctx *gin.Context) {
	var ct model.BaseIdsCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.LogicalRecovery(ctx, ct.Ids, c.appModule))
}

// PhysicalDeletion 物理删除
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *TenantAccountController) PhysicalDeletion(ctx *gin.Context) {
	var ct model.BaseIdsCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.PhysicalDeletion(ctx, ct.Ids, c.appModule))
}

// UpdatePassword 更新密码
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *TenantAccountController) UpdatePassword(ctx *gin.Context) {
	var ct modRamAccount.PasswordCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.ap.UpdatePassword(ctx, ct, c.appModule))
}

// Query 查询列表
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *TenantAccountController) Query(ctx *gin.Context) {
	var ct modRamAccount.QueryCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.Query(ctx, ct, c.appModule))
}

// Create 添加
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *TenantAccountController) Create(ctx *gin.Context) {
	var ct modRamAccount.CreateCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.Create(ctx, ct, c.appModule))
}

// CreateAccount 添加账号
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *TenantAccountController) CreateAccount(ctx *gin.Context) {
	var ct modRamAccount.CreateAccountCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.CreateAccount(ctx, ct, c.appModule))
}

// Update 更新
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *TenantAccountController) Update(ctx *gin.Context) {
	var ct modRamAccount.UpdateCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.Update(ctx, ct, c.appModule))
}

// UpdateAccount 更新
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *TenantAccountController) UpdateAccount(ctx *gin.Context) {
	var ct modRamAccount.UpdateAccountCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.UpdateAccount(ctx, ct, c.appModule))
}

// ExistAccount 查重
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *TenantAccountController) ExistAccount(ctx *gin.Context) {
	var ct model.BaseExistWdCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.ExistAccount(ctx, ct, c.appModule))
}

// ExistPhone 查重
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *TenantAccountController) ExistPhone(ctx *gin.Context) {
	var ct model.BaseExistWdCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.ExistPhone(ctx, ct, c.appModule))
}

// ExistMail 查重
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *TenantAccountController) ExistMail(ctx *gin.Context) {
	var ct model.BaseExistWdCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.ExistMail(ctx, ct, c.appModule))
}

// ExistIdentityCode 查重
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *TenantAccountController) ExistIdentityCode(ctx *gin.Context) {
	var ct model.BaseExistWdCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.ExistIdentityCode(ctx, ct, c.appModule))
}

// ExistRealName 查重
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *TenantAccountController) ExistRealName(ctx *gin.Context) {
	var ct model.BaseExistWdCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.ExistRealName(ctx, ct, c.appModule))
}
