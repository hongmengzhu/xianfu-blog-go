package bookmark

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/core/cache/cacheBookmark"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/domainBookmark/service"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/utils/webPg"
	"github.com/hongmengzhu/xianfu-blog-go/app/middleware/authPg"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/bm/modBmBookmarkWeb"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/routerPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/templatePg"
	"github.com/pangu-2/go-tools/tools/strPg"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(IndexController)).Name("BookmarkIndexController").Export(gs.As[routerPg.RouteRegistrar]())
}

type IndexController struct {
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
func (c *IndexController) RegisterRoutes(e *gin.Engine) {
	group := e.Group("", authPg.GroupWebMiddleware(c.Sp))
	group.GET("/bookmark", c.Index)
	group.GET("/bookmark/page/:page", c.Page)
}

func (c *IndexController) Index(ctx *gin.Context) {
	var ct modBmBookmarkWeb.QueryCt
	ctx.Bind(&ct)
	//
	rt := c.sv.Query(ctx, ct)
	log.Info(context.Background(), log.TagBizDef, "Data", log.Any("Pageable", rt.Data.Pageable))
	//
	tenantNo := webPg.GetTenantNo(ctx)
	tree, _ := c.catCache.FormatTree(ctx, tenantNo)

	//fmt.Printf("Data: %+v\n", tree)
	// 模版
	templatePg.HTML(ctx, "bookmark/index",
		templatePg.WithDataByResult(rt.SuccessIs(), rt.Data),
		// 页面 分类
		templatePg.WithHtmlObjSet("pageCategory", tree),
		templatePg.WithHtmlObjSet("pageUrl", "page"),
		templatePg.WithSitePage(templatePg.SitePage{
			Title:       "博客",
			Description: "博客",
			Keywords:    "博客",
			SiteName:    "博客",
		}))
}

func (c *IndexController) Page(ctx *gin.Context) {
	param := ctx.Param("page")
	var ct modBmBookmarkWeb.QueryCt
	ct.PageSize = 20
	ct.PageNum = max(strPg.ToInt64(param), 1)
	//
	rt := c.sv.Query(ctx, ct)
	// 模版
	templatePg.HTML(ctx, "bookmark/index_page",
		templatePg.WithDataByResult(rt.SuccessIs(), rt.Data),
		templatePg.WithHtmlObjSet("pageUrl", "page"),
		templatePg.WithSitePage(templatePg.SitePage{
			Title:       "博客",
			Description: "博客",
			Keywords:    "博客",
			SiteName:    "博客",
		}))
}
