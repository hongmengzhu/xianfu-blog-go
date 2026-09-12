package controller

import (
	"github.com/gin-gonic/gin"
	service2 "github.com/hongmengzhu/xianfu-blog-go/app/domain/system/ram/service"
	"github.com/hongmengzhu/xianfu-blog-go/app/middleware/authPg"
	"github.com/hongmengzhu/xianfu-blog-go/app/middleware/authPg/controllerPg"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/ram/modRamAccount"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/enumCommonPg/appModulePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/model"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/routerPg"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(AccountController).SetAppModule(appModulePg.System)).Name("SystemAccountController").Export(gs.As[routerPg.RouteRegistrar]())
}

// AccountController 账户
// @Description:
type AccountController struct {
	routerPg.RouteRegistrar
	controllerPg.SpSystemAuth
	sv        *service2.RamAccountService         `autowire:"?"`
	ap        *service2.RamAccountPasswordService `autowire:"?"`
	appModule appModulePg.AppModule
}

// SetAppModule 设置模块
// @Description:
func (c *AccountController) SetAppModule(appModule appModulePg.AppModule) *AccountController {
	c.appModule = appModule
	return c
}

// RegisterRoutes 注册路由
//
//	@Description:
//	@receiver c
//	@param e
func (c *AccountController) RegisterRoutes(e *gin.Engine) {
	group := e.Group("/xianfu/sys/ram/account", authPg.GroupSystemMiddleware(c.Sp))
	group.POST("/enable", c.Enable)
	group.POST("/disable", c.Disable)
	group.GET("/detail/:id", c.Detail)
	group.POST("/delete", c.Delete)
	group.POST("/recovery", c.Recovery)
	group.POST("/physicalDeletion", c.PhysicalDeletion)
	group.POST("/query", c.Query)
	group.POST("/updatePassword", c.UpdatePassword)
	group.POST("/createUpdateSimple", c.CreateUpdateSimple)
	group.POST("/createUpdate", c.CreateUpdate)
	group.POST("/existAccount", c.ExistAccount)
	group.POST("/existPhone", c.ExistPhone)
	group.POST("/existMail", c.ExistMail)
	group.POST("/existCode", c.ExistCode)
	group.POST("/existIdentityCode", c.ExistIdentityCode)
	group.POST("/existRealName", c.ExistRealName)
}

// Detail 详情
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *AccountController) Detail(ctx *gin.Context) {
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
func (c *AccountController) Enable(ctx *gin.Context) {
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
func (c *AccountController) Disable(ctx *gin.Context) {
	var ct model.BaseIdsCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.Disable(ctx, ct, c.appModule))
}

// Delete 逻辑删除
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *AccountController) Delete(ctx *gin.Context) {
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
func (c *AccountController) Recovery(ctx *gin.Context) {
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
func (c *AccountController) PhysicalDeletion(ctx *gin.Context) {
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
func (c *AccountController) UpdatePassword(ctx *gin.Context) {
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
func (c *AccountController) Query(ctx *gin.Context) {
	var ct modRamAccount.QueryCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.Query(ctx, ct, c.appModule))
}

// CreateUpdate 添加
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *AccountController) CreateUpdate(ctx *gin.Context) {
	var ct modRamAccount.CreateUpdateCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.CreateUpdate(ctx, ct, c.appModule))
}

func (c *AccountController) CreateUpdateSimple(ctx *gin.Context) {
	var ct modRamAccount.CreateUpdateAccountCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.CreateUpdateAccountSimple(ctx, ct, c.appModule))
}

// ExistAccount 查重
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *AccountController) ExistAccount(ctx *gin.Context) {
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
func (c *AccountController) ExistPhone(ctx *gin.Context) {
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
func (c *AccountController) ExistMail(ctx *gin.Context) {
	var ct model.BaseExistWdCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.ExistMail(ctx, ct, c.appModule))
}

func (c *AccountController) ExistCode(ctx *gin.Context) {
	var ct model.BaseExistWdCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.ExistCode(ctx, ct, c.appModule))
}

// ExistIdentityCode 查重
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *AccountController) ExistIdentityCode(ctx *gin.Context) {
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
func (c *AccountController) ExistRealName(ctx *gin.Context) {
	var ct model.BaseExistWdCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.ExistRealName(ctx, ct, c.appModule))
}
