package bookmark

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/core/blog/serviceCore"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/domainBlog/service"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/utils/webPg"
	"github.com/hongmengzhu/xianfu-blog-go/app/middleware/authPg"
	"github.com/hongmengzhu/xianfu-blog-go/app/middleware/components/cachePg/cacheDomainPg"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/blog/modBlogArticleWeb"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/routerPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/templatePg"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(CategoryController)).Name("BookmarkCategoryController").Export(gs.As[routerPg.RouteRegistrar]())
}

type CategoryController struct {
	routerPg.RouteRegistrar
	Sp       *authPg.GroupWebMiddlewareSp     `autowire:"?"`
	ca       *cacheDomainPg.TenantDomainCache `autowire:"?"`
	sv       *service.BlogArticleWebService   `autowire:"?"`
	catCache *serviceCore.CoreArticleCategory `autowire:"?"`
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

func (c *CategoryController) Code(ctx *gin.Context) {
	var ct modBlogArticleWeb.QueryCt
	ctx.Bind(&ct)
	param := ctx.Param("code")
	log.Info(ctx, log.TagBizDef, "Param", log.Any("code", param))
	//
	rt := c.sv.Query(ctx, ct)
	log.Info(context.Background(), log.TagBizDef, "Data", log.Any("Pageable", rt.Data.Pageable))
	//
	tenantNo := webPg.GetTenantNo(ctx)
	tree, _ := c.catCache.FormatTree(ctx, tenantNo)

	//fmt.Printf("Data: %+v\n", tree)
	// 模版
	templatePg.HTML(ctx, "bookmark/category-code",
		templatePg.WithDataByResult(rt.SuccessIs(), rt.Data),
		templatePg.WithHtmlObjSet("categorys", tree),
		templatePg.WithHtmlObjSet("pageUrl", "page"),
		templatePg.WithSitePage(templatePg.SitePage{
			Title:       "博客",
			Description: "博客",
			Keywords:    "博客",
			SiteName:    "博客",
		}))
}
func (c *CategoryController) List(ctx *gin.Context) {
	var ct modBlogArticleWeb.QueryCt
	ctx.Bind(&ct)
	//
	rt := c.sv.Query(ctx, ct)
	log.Info(context.Background(), log.TagBizDef, "Data", log.Any("Pageable", rt.Data.Pageable))
	//
	tenantNo := webPg.GetTenantNo(ctx)
	tree, _ := c.catCache.FormatTree(ctx, tenantNo)

	//fmt.Printf("Data: %+v\n", tree)
	// 模版
	templatePg.HTML(ctx, "bookmark/category",
		templatePg.WithDataByResult(rt.SuccessIs(), rt.Data),
		templatePg.WithHtmlObjSet("categorys", tree),
		templatePg.WithHtmlObjSet("pageUrl", "page"),
		templatePg.WithSitePage(templatePg.SitePage{
			Title:       "博客",
			Description: "博客",
			Keywords:    "博客",
			SiteName:    "博客",
		}))
}
