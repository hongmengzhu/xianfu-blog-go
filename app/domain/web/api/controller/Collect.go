package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/web/api/service"
	"github.com/hongmengzhu/xianfu-blog-go/app/middleware/authPg"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/blog/modBlogCollect/modApiBlogCollect"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/routerPg"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(CollectController)).Export(gs.As[routerPg.RouteRegistrar]())
}

type CollectController struct {
	routerPg.RouteRegistrar
	Sp *authPg.GroupApiMiddlewareSp `autowire:"?"`
	sv *service.CollectService      `autowire:"?"`
}

// RegisterRoutes
//
//	@Description: 注册路由
//	@receiver c
//	@param e
func (c *CollectController) RegisterRoutes(e *gin.Engine) {
	group := e.Group("/api/collect", authPg.GroupApiMiddleware(c.Sp))
	group.POST("/push", c.Push)
	group.POST("/pushAll", c.PushAll)
}

// Push 推送
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *CollectController) Push(ctx *gin.Context) {
	var ct modApiBlogCollect.PushCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.Push(ctx, ct))
}

// PushAll 推送
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *CollectController) PushAll(ctx *gin.Context) {
	var ct modApiBlogCollect.PushAll
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.PushAll(ctx, ct))
}
