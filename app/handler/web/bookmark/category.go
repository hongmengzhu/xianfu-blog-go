package bookmark

import (
	"context"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/domainBookmark/service"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/utils/webPg"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/bm/modBmBookmarkWeb"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/cacheService/cacheBookmark"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/model"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/routerPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/middleware/authPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/templatePg"
	"github.com/pangu-2/go-tools/tools/strPg"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(CategoryController)).Name("BookmarkCategoryController").Export(gs.As[routerPg.RouteRegistrar]())
}

type CategoryController struct {
	routerPg.RouteRegistrar
	Sp       *authPg.GroupWebMiddlewareSp         `autowire:"?"`
	sv       *service.BmBookmarkWebService        `autowire:"?"`
	catCache *cacheBookmark.CacheBookmarkCategory `autowire:"?"`
}

// RegisterRoutes
//
//	@Description: 注册路由
//	@receiver c
//	@param e
func (c *CategoryController) RegisterRoutes(e *gin.Engine) {
	group := e.Group("", authPg.GroupWebMiddleware(c.Sp))
	group.GET("/bookmark/category", c.List)
	group.GET("/bookmark/category/:code", c.Code)
}

// Code
func (c *CategoryController) Code(ctx *gin.Context) {
	htmlOption := make([]templatePg.Option, 0)
	htmlOption = append(htmlOption, templatePg.WithHtmlObjSet("pageUrl", "page"))
	htmlOption = append(htmlOption, templatePg.WithSitePage(templatePg.SitePage{
		Title:       "博客",
		Description: "博客",
		Keywords:    "博客",
		SiteName:    "博客",
	}))
	//
	var ct modBmBookmarkWeb.QueryCt
	param := ctx.Param("code")
	//
	tenantNo := webPg.GetTenantNo(ctx)
	//
	tree, _ := c.catCache.FormatTree(ctx, tenantNo)
	htmlOption = append(htmlOption, templatePg.WithHtmlObjSet("pageCategory", tree))
	//
	log.Info(ctx, log.TagBizDef, "Param", log.Any("code", param))
	if strPg.IsNotBlank(param) {
		key, b := c.catCache.GetByKey(ctx, param, tenantNo)
		if b {
			//
			ct.CategoryNo = key.No
			rt := c.sv.Query(ctx, ct)
			log.Info(context.Background(), log.TagBizDef, "Data", log.Any("Pageable", rt.Data.Pageable))
			htmlOption = append(htmlOption, templatePg.WithDataByResult(rt.SuccessIs(), rt.Data))
		}
	}
	// 模版
	templatePg.HTML(ctx, "bookmark/category-code", htmlOption...)
}
func (c *CategoryController) List(ctx *gin.Context) {
	htmlOption := make([]templatePg.Option, 0)
	htmlOption = append(htmlOption, templatePg.WithHtmlObjSet("pageUrl", "page"))
	htmlOption = append(htmlOption, templatePg.WithHtmlObjSet("pageCategory", make([]model.BaseNodeTree, 0)))
	htmlOption = append(htmlOption, templatePg.WithSitePage(templatePg.SitePage{
		Title:       "博客",
		Description: "博客",
		Keywords:    "博客",
		SiteName:    "博客",
	}))
	//
	//
	tenantNo := webPg.GetTenantNo(ctx)
	nosCategory := make([]string, 0)
	cate, b := c.catCache.GetAllByCache(ctx, tenantNo)
	if b {
		//
		for _, item := range cate {
			if strPg.IsNotBlank(item.No) && !slices.Contains(nosCategory, item.No) {
				nosCategory = append(nosCategory, item.No)
			}
		}
		//
		tree := c.catCache.ConvertToTree(cate)
		htmlOption = append(htmlOption, templatePg.WithHtmlObjSet("pageCategory", tree))
		//
		rt := c.sv.QueryListBYCategory(ctx, tenantNo, nosCategory, 20)
		if rt.SuccessIs() {
			templatePg.WithDataByResult(rt.SuccessIs(), rt.Data)
		}
	}
	// 模版
	templatePg.HTML(ctx, "bookmark/category", htmlOption...)
}
