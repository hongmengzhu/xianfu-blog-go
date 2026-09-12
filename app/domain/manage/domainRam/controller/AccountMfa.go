package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/manage/domainRam/service"
	"github.com/hongmengzhu/xianfu-blog-go/app/middleware/authPg"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/ram/modRamMfa"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/routerPg"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(AccountMfaController)).Name("ManageAccountMfaController").Export(gs.As[routerPg.RouteRegistrar]())
}

// AccountMfaController MFA 管理
type AccountMfaController struct {
	routerPg.RouteRegistrar
	Sp *authPg.GroupManageMiddlewareSp `autowire:""`
	sv *service.AccountMfaService      `autowire:"?"`
}

// RegisterRoutes 注册路由
func (c *AccountMfaController) RegisterRoutes(e *gin.Engine) {
	group := e.Group("/xianfu/manage/ram/account-mfa", authPg.GroupManageMiddleware(c.Sp))
	group.POST("/setup", c.Setup)
	group.POST("/setup-verify", c.SetupVerify)
	group.POST("/enable", c.Enable)
	group.POST("/disable", c.Disable)
	group.GET("/preferred", c.Preferred)
	group.POST("/recover", c.Recover)
}

// Setup 初始化 MFA 设置
func (c *AccountMfaController) Setup(ctx *gin.Context) {
	var ct modRamMfa.SetupCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.Setup(ctx, ct))
}

// SetupVerify 验证 MFA 设置
func (c *AccountMfaController) SetupVerify(ctx *gin.Context) {
	var ct modRamMfa.SetupVerifyCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.SetupVerify(ctx, ct))
}

// Enable 启用 MFA
func (c *AccountMfaController) Enable(ctx *gin.Context) {
	var ct modRamMfa.EnableCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.Enable(ctx, ct))
}

// Disable 禁用 MFA
func (c *AccountMfaController) Disable(ctx *gin.Context) {
	var ct modRamMfa.DisableCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.Disable(ctx, ct))
}

// Preferred 获取 MFA 状态
func (c *AccountMfaController) Preferred(ctx *gin.Context) {
	ctx.JSON(200, c.sv.Preferred(ctx))
}

// Recover 使用恢复码恢复
func (c *AccountMfaController) Recover(ctx *gin.Context) {
	var ct modRamMfa.RecoverCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.Recover(ctx, ct))
}
