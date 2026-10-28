package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/web/api/service"
	"github.com/hongmengzhu/xianfu-blog-go/app/middleware/authPg"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/blog/modBlogBookmark/modApiBlogBookmark"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/model"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/routerPg"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(BookmarkApiController)).Export(gs.As[routerPg.RouteRegistrar]())
}

type BookmarkApiController struct {
	routerPg.RouteRegistrar
	Sp *authPg.GroupApiMiddlewareSp `autowire:"?"`
	sv *service.BookmarkService     `autowire:"?"`
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
	group.POST("/update", c.Update)
	group.POST("/delete", c.Delete)
}

// GetAll 获取所有
func (c *BookmarkApiController) GetAll(ctx *gin.Context) {
	ctx.JSON(200, c.sv.GetAll(ctx))
}

// GetMy 获取我的
func (c *BookmarkApiController) GetMy(ctx *gin.Context) {
	var ct modApiBlogBookmark.QueryCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.GetMy(ctx, ct))
}

// GetOrg 获取组织
func (c *BookmarkApiController) GetOrg(ctx *gin.Context) {
	ctx.JSON(200, c.sv.GetOrg(ctx))
}

// Create 创建
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *BookmarkApiController) Create(ctx *gin.Context) {
	var ct modApiBlogBookmark.CreateUpdateCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.Create(ctx, ct))
}

// Update 更新
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *BookmarkApiController) Update(ctx *gin.Context) {
	var ct modApiBlogBookmark.CreateUpdateCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.Update(ctx, ct))
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
