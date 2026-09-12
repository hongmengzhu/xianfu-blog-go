package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/manage/domainBlog/service"
	"github.com/hongmengzhu/xianfu-blog-go/app/middleware/authPg"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/blog/modBlogBookmarkCategory"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/state/enumStatePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/model"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/routerPg"
	"github.com/pangu-2/go-tools/tools/strPg"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(BookmarkCategoryController)).Name("ManageBookmarkCategoryController").Export(gs.As[routerPg.RouteRegistrar]())
}

// BookmarkCategoryController 书签分类
// @Description:
type BookmarkCategoryController struct {
	routerPg.RouteRegistrar
	Sp *authPg.GroupManageMiddlewareSp      `autowire:""`
	sv *service.BlogBookmarkCategoryService `autowire:"?"`
}

// RegisterRoutes 注册路由
//
//	@Description:
//	@receiver c
//	@param e
func (c *BookmarkCategoryController) RegisterRoutes(e *gin.Engine) {
	group := e.Group("/xianfu/manage/blog/bookmark-category", authPg.GroupManageMiddleware(c.Sp))
	group.POST("/create", c.Create)
	group.POST("/update", c.Update)
	group.GET("/detail/:id", c.Detail)
	group.POST("/enable", c.Enable)
	group.POST("/disable", c.Disable)
	group.POST("/state", c.State)
	group.POST("/delete", c.Delete)
	group.POST("/recovery", c.Recovery)
	group.POST("/physicalDeletion", c.PhysicalDeletion)
	group.POST("/query", c.Query)
	group.POST("/queryAll", c.QueryAll)
	group.POST("/selectNodeAll", c.SelectNodeAll)
	group.POST("/selectNodeAllPublic", c.SelectNodeAllPublic)
	group.POST("/existName", c.ExistName)
	group.POST("/existNo", c.ExistNo)
}

// Create 创建
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *BookmarkCategoryController) Create(ctx *gin.Context) {
	var ct modBlogBookmarkCategory.CreateCt
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
func (c *BookmarkCategoryController) Update(ctx *gin.Context) {
	var ct modBlogBookmarkCategory.UpdateCt
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
func (c *BookmarkCategoryController) Delete(ctx *gin.Context) {
	var ct model.BaseIdsCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.LogicalDeletion(ctx, ct.Ids))
}

// Recovery 逻辑删除恢复
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *BookmarkCategoryController) Recovery(ctx *gin.Context) {
	var ct model.BaseIdsCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.LogicalRecovery(ctx, ct.Ids))
}

// PhysicalDeletion 物理删除
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *BookmarkCategoryController) PhysicalDeletion(ctx *gin.Context) {
	var ct model.BaseIdsCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.PhysicalDeletion(ctx, ct.Ids))
}

// Detail 详情
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *BookmarkCategoryController) Detail(ctx *gin.Context) {
	param := ctx.Param("id")
	if "" == param {
		ctx.JSON(200, rg.Error[string]("id不能为空"))
		return
	}
	ctx.JSON(200, c.sv.Detail(ctx, strPg.ToInt64(param)))
}

// Enable 启用
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *BookmarkCategoryController) Enable(ctx *gin.Context) {
	var ct model.BaseIdsCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.Enable(ctx, ct))
}

// Disable 禁用
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *BookmarkCategoryController) Disable(ctx *gin.Context) {
	var ct model.BaseIdsCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.Disable(ctx, ct))
}

// State 状态
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *BookmarkCategoryController) State(ctx *gin.Context) {
	var ct model.BaseStateIdsCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	state, ok := enumStatePg.IsExistInt64(ct.State)
	if !ok {
		ctx.JSON(200, rg.Error[string]("类型不正确"))
		return
	}
	ctx.JSON(200, c.sv.StateEnableDisable(ctx, ct.Ids, state))
}

// Query 查询列表
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *BookmarkCategoryController) Query(ctx *gin.Context) {
	var ct modBlogBookmarkCategory.QueryCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.Query(ctx, ct))
}

func (c *BookmarkCategoryController) SelectNodeAll(ctx *gin.Context) {
	var ct modBlogBookmarkCategory.QueryPublicCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ct.State = enumStatePg.ENABLE.IndexPg()
	ctx.JSON(200, c.sv.SelectNodeAll(ctx, ct))
}

func (c *BookmarkCategoryController) SelectNodeAllPublic(ctx *gin.Context) {
	var ct modBlogBookmarkCategory.QueryPublicCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ct.State = enumStatePg.ENABLE.IndexPg()
	ctx.JSON(200, c.sv.SelectNodeAllPublic(ctx, ct))
}

func (c *BookmarkCategoryController) QueryAll(ctx *gin.Context) {
	var ct modBlogBookmarkCategory.QueryPublicCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ct.State = enumStatePg.ENABLE.IndexPg()
	ctx.JSON(200, c.sv.QueryAll(ctx, ct))
}

// ExistName 查重
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *BookmarkCategoryController) ExistName(ctx *gin.Context) {
	var ct model.BaseExistWdCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.ExistName(ctx, ct))
}

// ExistNo 查重
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *BookmarkCategoryController) ExistNo(ctx *gin.Context) {
	var ct model.BaseExistWdCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.ExistNo(ctx, ct))
}
