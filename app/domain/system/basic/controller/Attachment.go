package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/system/basic/service"
	"github.com/hongmengzhu/xianfu-blog-go/app/middleware/authPg"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/basic/modBasicAttachment"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/routerPg"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(AttachmentController)).Name("SystemAttachmentController").Export(gs.As[routerPg.RouteRegistrar]())
}

type AttachmentController struct {
	routerPg.RouteRegistrar
	Sp *authPg.GroupSystemMiddlewareSp `autowire:"?"`
	sv *service.BasicAttachmentService `autowire:"?"`
}

func (c *AttachmentController) RegisterRoutes(e *gin.Engine) {
	group := e.Group("/xianfu/sys/basic/attachment", authPg.GroupSystemMiddleware(c.Sp))
	group.POST("/upload", c.Upload)
	group.POST("/upload-more", c.UploadMore)
	group.POST("/upload-link", c.UploadLink)
	group.POST("/upload-list", c.Query)
	group.POST("/query", c.Query)
}

func (c *AttachmentController) Upload(ctx *gin.Context) {
	ctx.JSON(200, c.sv.Upload(ctx))
}

func (c *AttachmentController) UploadMore(ctx *gin.Context) {
	ctx.JSON(200, c.sv.UploadMore(ctx))
}

func (c *AttachmentController) UploadLink(ctx *gin.Context) {
	var ct modBasicAttachment.WebUrlCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.UploadLink(ctx, ct))
}

func (c *AttachmentController) Query(ctx *gin.Context) {
	var ct modBasicAttachment.QueryCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.Query(ctx, ct))
}

func (c *AttachmentController) UpdateAddByFileOwner(ctx *gin.Context) {
	var ct modBasicAttachment.AddByFileOwnerCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.UpdateAddByFileOwner(ctx, ct))
}
