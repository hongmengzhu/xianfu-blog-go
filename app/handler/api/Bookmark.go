package api

import (
	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/domainBookmark/service"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/bm/modBmBookmark/modApiBmBookmark"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/blog/bookmarkTypeOwnerPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/model"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/routerPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/middleware/authPg"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(BookmarkApiController)).Export(gs.As[routerPg.RouteRegistrar]())
}

type BookmarkApiController struct {
	routerPg.RouteRegistrar
	Sp *authPg.GroupApiMiddlewareSp  `autowire:"?"`
	sv *service.BmBookmarkApiService `autowire:"?"`
}

// RegisterRoutes
//
//	@Description: 注册路由
//	@receiver c
//	@param e
func (c *BookmarkApiController) RegisterRoutes(e *gin.Engine) {
	group := e.Group("/api/bookmark", authPg.GroupApiMiddleware(c.Sp))
	group.POST("/getAll", c.GetAll)
	group.POST("/getMy", c.GetMy)
	group.POST("/getOrg", c.GetOrg)
	group.POST("/create", c.Create)
	group.POST("/createByOrg", c.CreateByOrg)
	group.POST("/updateByOrg", c.Create)
	group.POST("/update", c.Update)
	group.POST("/delete", c.Delete)
	group.POST("/setTypeReading", c.SetTypeReading)
}

// GetAll 获取所有
func (c *BookmarkApiController) GetAll(ctx *gin.Context) {
	ctx.JSON(200, c.sv.GetAll(ctx))
}

// GetMy 获取我的
func (c *BookmarkApiController) GetMy(ctx *gin.Context) {
	var ct modApiBmBookmark.QueryCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.GetMy(ctx, ct))
}

// GetOrg 获取组织
func (c *BookmarkApiController) GetOrg(ctx *gin.Context) {
	var ct modApiBmBookmark.QueryCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.GetOrg(ctx, ct))
}

// Create 创建
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *BookmarkApiController) Create(ctx *gin.Context) {
	var ct modApiBmBookmark.CreateUpdateCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ct.TypeOwner = bookmarkTypeOwnerPg.MY.Code()
	ctx.JSON(200, c.sv.Create(ctx, ct))
}

// CreateByOrg 创建
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *BookmarkApiController) CreateByOrg(ctx *gin.Context) {
	var ct modApiBmBookmark.CreateUpdateCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ct.TypeOwner = bookmarkTypeOwnerPg.ORG.Code()
	ctx.JSON(200, c.sv.Create(ctx, ct))
}

// Update 更新
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *BookmarkApiController) Update(ctx *gin.Context) {
	var ct modApiBmBookmark.CreateUpdateCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.Update(ctx, ct))
}

// UpdateByOrg 更新
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *BookmarkApiController) UpdateByOrg(ctx *gin.Context) {
	var ct modApiBmBookmark.CreateUpdateCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.Update(ctx, ct))
}

// SetTypeReading 设置读书状态
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *BookmarkApiController) SetTypeReading(ctx *gin.Context) {
	var ct modApiBmBookmark.TypeReadingCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.SetTypeReading(ctx, ct))
}

// Delete 逻辑删除
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *BookmarkApiController) Delete(ctx *gin.Context) {
	var ct model.BaseIdsCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.LogicalDeletion(ctx, ct.Ids))
}
