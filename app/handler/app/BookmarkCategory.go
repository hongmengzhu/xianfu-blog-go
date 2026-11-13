package app

import (
	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/domainBm/service"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/bm/modBmBookmarkCategory/modApiBmBookmarkCategory"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/model"
	routerPg "github.com/hongmengzhu/xianfu-blog-go/pkg/core/routerPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/middleware/authPg"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(BookmarkCategoryAppController)).Export(gs.As[routerPg.RouteRegistrar]())
}

type BookmarkCategoryAppController struct {
	routerPg.RouteRegistrar
	Sp *authPg.GroupApiMiddlewareSp          `autowire:"?"`
	sv *service.BmBookmarkCategoryApiService `autowire:"?"`
}

// RegisterRoutes
//
//	@Description: 注册路由
//	@receiver c
//	@param e
func (c *BookmarkCategoryAppController) RegisterRoutes(e *gin.Engine) {
	group := e.Group("/api-app/bookmark-category", authPg.GroupApiMiddleware(c.Sp))
	group.POST("/getAll", c.GetAll)
	group.POST("/createUpdate", c.CreateUpdate)
	group.POST("/delete", c.Delete)
	group.POST("/getMy	", c.GetMy)
	group.POST("/getOrg	", c.GetOrg)
}

// GetAll 获取所有
func (c *BookmarkCategoryAppController) GetAll(ctx *gin.Context) {
	ctx.JSON(200, c.sv.GetAll(ctx))
}

// GetMy 获取我的
func (c *BookmarkCategoryAppController) GetMy(ctx *gin.Context) {
	ctx.JSON(200, c.sv.GetMy(ctx))
}

// GetOrg 获取组织
func (c *BookmarkCategoryAppController) GetOrg(ctx *gin.Context) {
	ctx.JSON(200, c.sv.GetOrg(ctx))
}

// CreateUpdate 创建
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *BookmarkCategoryAppController) CreateUpdate(ctx *gin.Context) {
	var ct modApiBmBookmarkCategory.CreateUpdateCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.CreateUpdate(ctx, ct))
}

// Delete 逻辑删除
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *BookmarkCategoryAppController) Delete(ctx *gin.Context) {
	var ct model.BaseIdsCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.LogicalDeletion(ctx, ct.Ids))
}
