package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/manage/domainBasic/service"
	"github.com/hongmengzhu/xianfu-blog-go/app/middleware/authPg"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/basic/modBasicAttachment"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/routerPg"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(AttachmentController)).Name("ManageAttachmentController").Export(gs.As[routerPg.RouteRegistrar]())
}

// AttachmentController 附件上传
// @Description:
type AttachmentController struct {
	routerPg.RouteRegistrar
	Sp *authPg.GroupManageMiddlewareSp `autowire:""`
	sv *service.BasicAttachmentService `autowire:""`
}

// RegisterRoutes 注册路由
//
//	@Description:
//	@receiver c
//	@param e
func (c *AttachmentController) RegisterRoutes(e *gin.Engine) {
	group := e.Group("/xianfu/manage/basic/attachment", authPg.GroupManageMiddleware(c.Sp))
	group.POST("/upload", c.Upload)
	group.POST("/upload-more", c.UploadMore)
	group.POST("/upload-link", c.UploadLink)
	group.POST("/upload-list", c.Query)
	group.POST("/query", c.Query)
	group.POST("/makeFileOwnerPublic", c.MakeFileOwner)
	group.POST("/makeFileOwnerAllPublic", c.MakeFileOwnerAll)
	group.POST("/upload-listByOwnerPublic", c.ListByOwner)
	group.POST("/upload-detail", c.DetailListByOwner)
	group.POST("/upload-delByOwnerPublic", c.DelByOwner)
	group.POST("/upload-makeFileOwnerAllPublic", c.MakeFileOwnerAll)
	group.POST("/upload-updateByFileOwner", c.UpdateByFileOwner)
	group.POST("/upload-addByFileOwner", c.AddByFileOwnerCt)
	group.POST("/upload-ownerDel", c.DelByOwner)
}

// Upload
//
//	@Description: 但文件上传
//	@receiver c
//	@param ctx
func (c *AttachmentController) Upload(ctx *gin.Context) {
	ctx.JSON(200, c.sv.Upload(ctx))
}

// UploadMore
//
//	@Description:  多文件上传
//	@receiver c
//	@param ctx
func (c *AttachmentController) UploadMore(ctx *gin.Context) {
	ctx.JSON(200, c.sv.UploadMore(ctx))
}

// UploadLink
//
//	@Description:  多url文件上传
//	@receiver c
//	@param ctx
func (c *AttachmentController) UploadLink(ctx *gin.Context) {
	var ct modBasicAttachment.WebUrlCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.UploadLink(ctx, ct))
}

// ListByOwner
//
//	@Description:  根据文件拥有者查询
//	@receiver c
//	@param ctx
func (c *AttachmentController) ListByOwner(ctx *gin.Context) {
	var ct modBasicAttachment.ListFileOwnerCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.ListByOwner(ctx, ct))
}

// DetailListByOwner
//
//	@Description:  根据文件拥有者查询
//	@receiver c
//	@param ctx
func (c *AttachmentController) DetailListByOwner(ctx *gin.Context) {
	var ct modBasicAttachment.ListFileOwnerCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.DetailListByOwner(ctx, ct))
}

// DelByOwner
//
//	@Description:  根据文件拥有者查询
//	@receiver c
//	@param ctx
func (c *AttachmentController) DelByOwner(ctx *gin.Context) {
	var ct modBasicAttachment.DelFileOwnerCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.DelByOwner(ctx, ct))
}

// Query 查询列表
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *AttachmentController) Query(ctx *gin.Context) {
	var ct modBasicAttachment.QueryCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.Query(ctx, ct))
}

// MakeFileOwner
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *AttachmentController) MakeFileOwner(ctx *gin.Context) {
	var ct modBasicAttachment.MakeFileOwnerCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.MakeFileOwner(ctx, ct))
}

// MakeFileOwnerAll
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *AttachmentController) MakeFileOwnerAll(ctx *gin.Context) {
	var ct modBasicAttachment.MakeFileOwnerAllCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.MakeFileOwnerAll(ctx, ct))
}

// UpdateByFileOwner
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *AttachmentController) UpdateByFileOwner(ctx *gin.Context) {
	var ct modBasicAttachment.UpdateByFileOwner
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.UpdateByFileOwner(ctx, ct))
}

// AddByFileOwnerCt
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *AttachmentController) AddByFileOwnerCt(ctx *gin.Context) {
	var ct modBasicAttachment.AddByFileOwnerCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.AddByFileOwner(ctx, ct))
}
