package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	service2 "github.com/hongmengzhu/xianfu-blog-go/app/domain/manage/domainRam/service"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/ram/modRamCas"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/ram/modRamLogin"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/ram/modRamMfa"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/ram/modRamSaml"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/routerPg"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(AuthLoginIdpController)).Name("ManageAuthLoginIdpController").Export(gs.As[routerPg.RouteRegistrar]())
}

// AuthLoginIdpController IDP OAuth 登录
// @Description:
type AuthLoginIdpController struct {
	routerPg.RouteRegistrar
	sv            *service2.AccountLoginIdpService `autowire:"?"`
	samlService   *service2.SamlLoginService       `autowire:"?"`
	mfaService    *service2.AccountMfaService      `autowire:"?"`
	faceIdService *service2.FaceIdService          `autowire:"?"`
	casService    *service2.CasService             `autowire:"?"`
}

// RegisterRoutes 注册路由
//
//	@Description:
//	@receiver c
//	@param e
func (c *AuthLoginIdpController) RegisterRoutes(e *gin.Engine) {
	group := e.Group("/xianfu/auth/manage/idp")
	group.POST("/login", c.Login)
	group.POST("/refresh", c.RefreshToken)
	group.POST("/mfa/verify", c.MfaVerify)
	group.POST("/mfa/recover", c.MfaRecover)

	// SAML 端点
	group.GET("/saml/login/:sourceNo", c.SamlLogin)
	group.POST("/saml/callback", c.SamlCallback)

	// Face ID 端点
	group.POST("/faceid/begin", c.FaceIdBegin)
	group.POST("/faceid/verify", c.FaceIdVerify)

	// CAS 登录端点
	group.POST("/cas/login", c.CasLogin)
}

// Login IDP OAuth 登陆
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *AuthLoginIdpController) Login(ctx *gin.Context) {
	var ct modRamLogin.IdpLoginCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.Login(ctx, ct))
}

// RefreshToken 刷新 token
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *AuthLoginIdpController) RefreshToken(ctx *gin.Context) {
	var ct modRamLogin.TokenRefreshCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.RefreshToken(ctx, ct))
}

// MfaVerify MFA 验证（验证通过后签发 Token）
func (c *AuthLoginIdpController) MfaVerify(ctx *gin.Context) {
	var ct modRamMfa.VerifyCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.mfaService.VerifyLogin(ctx, ct))
}

// MfaRecover MFA 恢复码恢复
func (c *AuthLoginIdpController) MfaRecover(ctx *gin.Context) {
	var ct modRamMfa.RecoverCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.mfaService.Recover(ctx, ct))
}

// SamlLogin SAML 登录 - 生成 SAML 认证请求 URL 并重定向
func (c *AuthLoginIdpController) SamlLogin(ctx *gin.Context) {
	sourceNo := ctx.Param("sourceNo")
	if sourceNo == "" {
		ctx.JSON(200, rg.Error[string]("认证源编号不能为空"))
		return
	}

	result := c.samlService.GetSamlLoginUrl(ctx, sourceNo)
	if result.ErrorIs() {
		ctx.JSON(200, rg.Error[string](result.Message))
		return
	}

	vo := result.Data
	if vo.Method == "POST" && vo.PostBody != "" {
		// POST 绑定：返回自动提交的 HTML 表单
		ctx.Header("Content-Type", "text/html; charset=utf-8")
		html := fmt.Sprintf(`<!DOCTYPE html>
<html><body>
<form id="samlForm" method="POST" action="%s">
%s
<input type="hidden" name="RelayState" value="%s" />
</form>
<script>document.getElementById('samlForm').submit();</script>
</body></html>`, vo.RedirectUrl, vo.PostBody, vo.RedirectUrl)
		ctx.String(http.StatusOK, html)
	} else {
		// GET 绑定：重定向
		ctx.Redirect(http.StatusFound, vo.RedirectUrl)
	}
}

// SamlCallback SAML Response 回调端点（ACS）
func (c *AuthLoginIdpController) SamlCallback(ctx *gin.Context) {
	ct := modRamSaml.SamlCallbackCt{
		SAMLResponse: ctx.PostForm("SAMLResponse"),
		RelayState:   ctx.PostForm("RelayState"),
	}
	if ct.SAMLResponse == "" {
		ct.SAMLResponse = ctx.Query("SAMLResponse")
	}
	if ct.RelayState == "" {
		ct.RelayState = ctx.Query("RelayState")
	}

	result := c.samlService.HandleSamlCallback(ctx, ct)
	if result.ErrorIs() {
		ctx.JSON(200, rg.Error[string](result.Message))
		return
	}
	ctx.JSON(200, rg.OkData(result.Data))
}

// FaceIdBegin Face ID 开始验证
func (c *AuthLoginIdpController) FaceIdBegin(ctx *gin.Context) {
	var ct service2.FaceIdBeginCt
	if err := ctx.ShouldBind(&ct); err != nil {
		ctx.JSON(200, rg.ErrorDefault[string]())
		return
	}
	ctx.JSON(200, c.faceIdService.Begin(ctx, ct))
}

// FaceIdVerify Face ID 验证
func (c *AuthLoginIdpController) FaceIdVerify(ctx *gin.Context) {
	var ct service2.FaceIdVerifyCt
	if err := ctx.ShouldBind(&ct); err != nil {
		ctx.JSON(200, rg.ErrorDefault[string]())
		return
	}
	ctx.JSON(200, c.faceIdService.Verify(ctx, ct))
}

// CasLogin CAS 登录（生成 Service Ticket）
func (c *AuthLoginIdpController) CasLogin(ctx *gin.Context) {
	var ct modRamCas.CasLoginCt
	if err := ctx.ShouldBind(&ct); err != nil {
		ctx.JSON(200, rg.ErrorDefault[string]())
		return
	}
	ctx.JSON(200, c.casService.Login(ctx, ct))
}
