package bookmark

import (
	"slices"

	"github.com/hongmengzhu/xianfu-blog-go/app/domain/domainBookmark/service"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/utils/webPg"
	"github.com/hongmengzhu/xianfu-blog-go/app/handlerVhost"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/bm/modBmBookmarkWeb"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/cacheService/cacheBookmark"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/model"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/middleware/authPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/middleware/staterEchoV5"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/templatePg"
	"github.com/labstack/echo/v5"
	"github.com/pangu-2/go-tools/tools/strPg"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(CategoryController)).Name("BookmarkVhostCategoryController").Export(gs.As[staterEchoV5.RouteRegistrar]())
}

// CategoryController 书签分类（独立域名 bookmark.test）
type CategoryController struct {
	staterEchoV5.RouteRegistrar
	Sp       *authPg.GroupWebVhostMiddlewareSp    `autowire:"?"`
	sv       *service.BmBookmarkWebService        `autowire:"?"`
	catCache *cacheBookmark.CacheBookmarkCategory `autowire:"?"`
}

// Host 模块绑定域名
func (c *CategoryController) Host() string {
	return handlerVhost.Host_bookmark
}

// RegisterRoutes
//
//	@Description: 注册路由
//	@receiver c
//	@param g
func (c *CategoryController) RegisterRoutes(g *echo.Echo) {
	group := g.Group("", authPg.GroupWebVhostMiddlewareEcho(c.Sp))
	group.GET("/category", c.List)
	group.GET("/category/:code", c.Code)
}

// List 分类总览
func (c *CategoryController) List(ctx *echo.Context) error {
	htmlOption := make([]templatePg.Option, 0)
	htmlOption = append(htmlOption, templatePg.WithHtmlObjSet("pageUrl", "/category"))
	htmlOption = append(htmlOption, templatePg.WithHtmlObjSet("pageCategory", make([]model.BaseNodeTree, 0)))
	htmlOption = append(htmlOption, templatePg.WithSitePage(templatePg.SitePage{
		Title:       "AI 工具分类",
		Description: "按分类浏览 AI 工具",
		Keywords:    "AI工具分类,AI导航",
		SiteName:    "分宝呀",
	}))
	//
	cxt := ctx.Request().Context()
	tenantNo := webPg.GetTenantNoEcho(ctx)
	nosCategory := make([]*model.BaseNodeDataTree, 0)
	nosIds := make([]string, 0)
	cate, b := c.catCache.GetAllByCache(cxt, tenantNo)
	if b {
		for _, item := range cate {
			if strPg.IsNotBlank(item.No) && !slices.Contains(nosIds, item.No) {
				nosIds = append(nosIds, item.No)
				nosCategory = append(nosCategory, &model.BaseNodeDataTree{
					No:       item.No,
					Name:     item.Name,
					ParentNo: item.ParentNo,
					Code:     item.Code,
					Extend:   item,
				})
			}
		}
		htmlOption = append(htmlOption, templatePg.WithHtmlObjSet("pageCategory", c.catCache.ConvertToTree(cate)))
		if rt := c.sv.QueryListBYCategory(cxt, tenantNo, nosCategory, 20); rt.SuccessIs() {
			htmlOption = append(htmlOption, templatePg.WithDataByResult(rt.SuccessIs(), rt.Data))
		}
	}
	// 模版
	return templatePg.HTMLEcho(ctx, "bookmark/category", htmlOption...)
}

// Code 分类详情（分页）
func (c *CategoryController) Code(ctx *echo.Context) error {
	param := ctx.Param("code")
	htmlOption := make([]templatePg.Option, 0)
	htmlOption = append(htmlOption, templatePg.WithHtmlObjSet("pageUrl", "/category/"+param))
	htmlOption = append(htmlOption, templatePg.WithSitePage(templatePg.SitePage{
		Title:       "AI 工具分类",
		Description: "按分类浏览 AI 工具",
		Keywords:    "AI工具分类,AI导航",
		SiteName:    "分宝呀",
	}))
	//
	cxt := ctx.Request().Context()
	tenantNo := webPg.GetTenantNoEcho(ctx)
	tree, _ := c.catCache.FormatTree(cxt, tenantNo)
	htmlOption = append(htmlOption, templatePg.WithHtmlObjSet("pageCategory", tree))
	//
	var ct modBmBookmarkWeb.QueryCt
	ct.PageSize = 20
	ct.PageNum = max(strPg.ToInt64(ctx.QueryParam("pageNum")), 1)
	//
	log.Info(cxt, log.TagBizDef, "Param", log.Any("code", param))
	if strPg.IsNotBlank(param) {
		if key, b := c.catCache.GetByKey(cxt, param, tenantNo); b {
			htmlOption = append(htmlOption, templatePg.WithHtmlObjSet("pageCategoryItem", key))
			ct.CategoryNo = key.No
			rt := c.sv.QueryByContext(cxt, tenantNo, ct)
			log.Info(cxt, log.TagBizDef, "Data", log.Any("Pageable", rt.Data.Pageable))
			htmlOption = append(htmlOption, templatePg.WithDataByResult(rt.SuccessIs(), rt.Data))
		}
	}
	// 模版
	return templatePg.HTMLEcho(ctx, "bookmark/category-code", htmlOption...)
}
