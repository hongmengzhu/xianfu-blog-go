package controller

import (
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/manage/domainRam/service"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/routerPg"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(WechatOaController)).Name("ManageWechatOaController").Export(gs.As[routerPg.RouteRegistrar]())
}

// WechatOaController 微信公众号扫码登录
// @Description:
type WechatOaController struct {
	routerPg.RouteRegistrar
	sv *service.WechatOaService `autowire:"?"`
}

// RegisterRoutes 注册路由
//
//	@Description:
//	@receiver c
//	@param e
func (c *WechatOaController) RegisterRoutes(e *gin.Engine) {
	group := e.Group("/xianfu/auth/manage/wechat")
	group.GET("/webhook", c.WebhookVerify)
	group.POST("/webhook", c.WebhookEvent)
	group.GET("/qrcode/:sourceNo", c.GetQRCode)
	group.GET("/poll/:ticket", c.PollScanStatus)
}

// WebhookVerify 微信服务器验证（GET /webhook）
//
//	@Description: 微信服务器发送 GET 请求验证 URL 有效性，需原样返回 echostr
//	@receiver c
//	@param ctx
func (c *WechatOaController) WebhookVerify(ctx *gin.Context) {
	signature := ctx.Query("signature")
	timestamp := ctx.Query("timestamp")
	nonce := ctx.Query("nonce")
	echostr := ctx.Query("echostr")

	// 尝试从 sourceNo 参数获取 token
	sourceNo := ctx.Query("sourceNo")
	wechatToken := ""
	if sourceNo != "" {
		wechatToken = c.sv.GetWechatOaToken(ctx, sourceNo)
	}

	// 如果配置了 token 则验证签名
	if wechatToken != "" {
		if !c.sv.VerifySignature(wechatToken, nonce, timestamp, signature) {
			log.Warnf(ctx, log.TagAppDef, "微信 webhook 签名验证失败: sourceNo=%s", sourceNo)
			ctx.String(http.StatusOK, "invalid signature")
			return
		}
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
func (c *WechatOaController) WebhookEvent(ctx *gin.Context) {
	bodyBytes, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Errorf(ctx, log.TagAppDef, "读取微信 webhook body 失败: %v", err)
		ctx.String(http.StatusOK, "")
		return
	}

	sourceNo := ctx.Query("sourceNo")
	result := c.sv.HandleEvent(ctx, bodyBytes, sourceNo)
	if result.ErrorIs() {
		log.Errorf(ctx, log.TagAppDef, "处理微信事件失败: %s", result.Message)
	}

	// 微信要求返回空字符串表示成功
	ctx.String(http.StatusOK, "")
}

// GetQRCode 获取扫码登录二维码
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *WechatOaController) GetQRCode(ctx *gin.Context) {
	sourceNo := ctx.Param("sourceNo")
	ctx.JSON(http.StatusOK, c.sv.GetQRCode(ctx, sourceNo))
}

// PollScanStatus 轮询扫码状态
//
//	@Description: 前端定时查询扫码状态
//	@receiver c
//	@param ctx
func (c *WechatOaController) PollScanStatus(ctx *gin.Context) {
	ticket := ctx.Param("ticket")
	result := c.sv.PollScanStatus(ctx, ticket)
	ctx.JSON(http.StatusOK, result)
}
