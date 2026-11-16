package idp

import (
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/domainIdp/service"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/routerPg"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(AuthWechatOaController)).Name("ManageIdpWechatOaController").Export(gs.As[routerPg.RouteRegistrar]())
}

// AuthWechatOaController 微信公众号扫码登录
// @Description:
type AuthWechatOaController struct {
	routerPg.RouteRegistrar
	sv *service.WechatOaService `autowire:"?"`
}

// RegisterRoutes 注册路由
//
//	@Description:
//	@receiver c
//	@param e
func (c *AuthWechatOaController) RegisterRoutes(e *gin.Engine) {
	group := e.Group("/xianfu/auth/manage/wechat")
	group.GET("/webhook", c.WebhookVerify)
	group.POST("/webhook", c.WebhookEvent)
	group.GET("/qrcode/:sourceNo", c.GetQRCode)
	group.GET("/poll/:ticket", c.PollScanStatus)
}

func (c *AuthWechatOaController) verifyWebhook(ctx *gin.Context) bool {
	sourceNo := strings.TrimSpace(ctx.Query("sourceNo"))
	if sourceNo == "" {
		return false
	}
	wechatToken := c.sv.GetWechatOaToken(ctx, sourceNo)
	if strings.TrimSpace(wechatToken) == "" {
		return false
	}
	signature := ctx.Query("signature")
	timestamp := ctx.Query("timestamp")
	nonce := ctx.Query("nonce")
	if signature == "" || timestamp == "" || nonce == "" {
		return false
	}
	ts, err := strconv.ParseInt(timestamp, 10, 64)
	if nil != err || time.Now().Unix()-ts > 300 || ts-time.Now().Unix() > 300 {
		return false
	}
	return c.sv.VerifySignature(wechatToken, nonce, timestamp, signature)
}

// WebhookVerify 微信服务器验证（GET /webhook）
//
//	@Description: 微信服务器发送 GET 请求验证 URL 有效性，需原样返回 echostr
//	@receiver c
//	@param ctx
func (c *AuthWechatOaController) WebhookVerify(ctx *gin.Context) {
	echostr := ctx.Query("echostr")

	// 必须配置 token 且签名有效，未配置时不再放行。
	if !c.verifyWebhook(ctx) {
		log.Warn(ctx, log.TagAppDef, "微信 webhook 签名验证失败", log.Any("sourceNo", ctx.Query("sourceNo")))
		ctx.String(http.StatusForbidden, "invalid signature")
		return
	}

	// 原样返回 echostr（微信要求返回纯文本）
	echoInt, err := strconv.Atoi(echostr)
	if err == nil {
		ctx.String(http.StatusOK, strconv.Itoa(echoInt))
	} else {
		ctx.String(http.StatusOK, echostr)
	}
}

// WebhookEvent 微信事件推送（POST /webhook）
//
//	@Description: 接收微信扫描事件推送（SCAN / SUBSCRIBE）
//	@receiver c
//	@param ctx
func (c *AuthWechatOaController) WebhookEvent(ctx *gin.Context) {
	if !c.verifyWebhook(ctx) {
		log.Warn(ctx, log.TagAppDef, "微信事件签名验证失败", log.Any("sourceNo", ctx.Query("sourceNo")))
		ctx.String(http.StatusForbidden, "")
		return
	}
	bodyBytes, err := io.ReadAll(io.LimitReader(ctx.Request.Body, 1<<20))
	if err != nil {
		log.Error(ctx, log.TagAppDef, err, "读取微信 webhook body 失败")
		ctx.String(http.StatusOK, "")
		return
	}

	sourceNo := ctx.Query("sourceNo")
	result := c.sv.HandleEvent(ctx, bodyBytes, sourceNo)
	if result.ErrorIs() {
		log.Error(ctx, log.TagAppDef, nil, "处理微信事件失败:", log.Any("err", result.Message))
	}

	// 微信要求返回空字符串表示成功
	ctx.String(http.StatusOK, "")
}

// GetQRCode 获取扫码登录二维码
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *AuthWechatOaController) GetQRCode(ctx *gin.Context) {
	sourceNo := ctx.Param("sourceNo")
	ctx.JSON(http.StatusOK, c.sv.GetQRCode(ctx, sourceNo))
}

// PollScanStatus 轮询扫码状态
//
//	@Description: 前端定时查询扫码状态
//	@receiver c
//	@param ctx
func (c *AuthWechatOaController) PollScanStatus(ctx *gin.Context) {
	ticket := ctx.Param("ticket")
	result := c.sv.PollScanStatus(ctx, ticket)
	ctx.JSON(http.StatusOK, result)
}
