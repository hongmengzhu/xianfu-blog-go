package blog

import (
	"context"
	"net/http"

	"github.com/hongmengzhu/xianfu-blog-go/app/core/blog/serviceCore"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/host"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/web/blog/model/modBlogArticle"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/web/blog/service"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/web/utils/webPg"
	"github.com/hongmengzhu/xianfu-blog-go/app/middleware/authPg"
	"github.com/hongmengzhu/xianfu-blog-go/app/middleware/staterEchoV5"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/templatePg"
	"github.com/labstack/echo/v5"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(ArticleController)).Name("BlogHostArticleController").Export(gs.As[staterEchoV5.RouteRegistrar]())
}

// ArticleController 文章
type ArticleController struct {
	staterEchoV5.RouteRegistrar
	Sp       *authPg.GroupWebVhostMiddlewareSp `autowire:"?"`
	sv       *service.ArticleService           `autowire:"?"`
	catCache *serviceCore.CoreArticleCategory  `autowire:"?"`
}

// Hosts 模块绑定域名
func (c *ArticleController) Host() string {
	return host.Host_blog
}

// RegisterRoutes
//
//	@Description: 注册路由
//	@receiver c
//	@param g
func (c *ArticleController) RegisterRoutes(g *echo.Echo) {
	group := g.Group("", authPg.GroupWebVhostMiddlewareEcho(c.Sp))
	group.GET("/article/search", c.List)
	group.GET("/article/:id", c.Detail)
}

func (c *ArticleController) Detail(ctx *echo.Context) error {
	dataIs := false
	var data any
	param := ctx.Param("id")
	tenantNo := webPg.GetTenantNoEcho(ctx)
	rt := c.sv.DetailByContext(ctx.Request().Context(), tenantNo, param)
	if rt.SuccessIs() {
		dataIs = true
		data = rt.Data
	}
	//
	tree, _ := c.catCache.FormatTree(ctx.Request().Context(), tenantNo)
	// 文章不存在时返回主题化 404 页
	if !dataIs {
		return templatePg.HTMLEcho(ctx, "blog/404",
			templatePg.WithCode(http.StatusNotFound),
			templatePg.WithHtmlObjSet("categorys", tree),
			templatePg.WithHtmlObjSet("pageUrl", "article"),
			templatePg.WithSitePage(templatePg.SitePage{
				Title:       "页面不存在",
				Description: "博客",
				Keywords:    "博客",
				SiteName:    "博客",
			}))
	}
	// 模版
	return templatePg.HTMLEcho(ctx, "blog/detail",
		templatePg.WithDataByResult(dataIs, data),
		templatePg.WithHtmlObjSet("categorys", tree),
		templatePg.WithHtmlObjSet("pageUrl", "article"),
		templatePg.WithSitePage(templatePg.SitePage{
			Title:       "详情",
			Description: "博客",
			Keywords:    "博客",
			SiteName:    "博客",
		}))
}

func (c *ArticleController) List(ctx *echo.Context) error {
	var ct modBlogArticle.QueryCt
	if err := ctx.Bind(&ct); nil != err {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	//
	tenantNo := webPg.GetTenantNoEcho(ctx)
	rt := c.sv.QueryByContext(ctx.Request().Context(), tenantNo, ct)
	log.Info(context.Background(), log.TagBizDef, "Data", log.Any("Pageable", rt.Data.Pageable))
	//
	tree, _ := c.catCache.FormatTree(ctx.Request().Context(), tenantNo)

	// 模版
	return templatePg.HTMLEcho(ctx, "blog/article_list",
		templatePg.WithDataByResult(rt.SuccessIs(), rt.Data),
		templatePg.WithHtmlObjSet("categorys", tree),
		templatePg.WithHtmlObjSet("pageUrl", "/article/search"),
		templatePg.WithSitePage(templatePg.SitePage{
			Title:       "博客",
			Description: "博客",
			Keywords:    "博客",
			SiteName:    "博客",
		}))
}
