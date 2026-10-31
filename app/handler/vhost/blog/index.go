package blog

import (
	"context"
	"net/http"

	"github.com/hongmengzhu/xianfu-blog-go/app/core/cache/cacheBlog"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/domainBlog/service"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/utils/webPg"
	"github.com/hongmengzhu/xianfu-blog-go/app/handler/vhost"
	"github.com/hongmengzhu/xianfu-blog-go/app/middleware/authPg"
	"github.com/hongmengzhu/xianfu-blog-go/app/middleware/staterEchoV5"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/blog/modBlogArticleWeb"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/templatePg"
	"github.com/labstack/echo/v5"
	"github.com/pangu-2/go-tools/tools/strPg"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(IndexController)).Name("BlogHostIndexController").Export(gs.As[staterEchoV5.RouteRegistrar]())
}

// IndexController 首页
type IndexController struct {
	staterEchoV5.RouteRegistrar
	Sp       *authPg.GroupWebVhostMiddlewareSp   `autowire:"?"`
	sv       *service.BlogArticleWebService      `autowire:"?"`
	catCache *cacheBlog.CacheBlogArticleCategory `autowire:"?"`
}

// Hosts 模块绑定域名
func (c *IndexController) Host() string {
	return vhost.Host_blog
}

// RegisterRoutes
//
//	@Description: 注册路由
//	@receiver c
//	@param g
func (c *IndexController) RegisterRoutes(g *echo.Echo) {
	group := g.Group("", authPg.GroupWebVhostMiddlewareEcho(c.Sp))
	group.GET("/", c.Index)
	group.GET("/page/:page", c.Page)
}

func (c *IndexController) Index(ctx *echo.Context) error {
	var ct modBlogArticleWeb.QueryCt
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
	return templatePg.HTMLEcho(ctx, "blog/index",
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

func (c *IndexController) Page(ctx *echo.Context) error {
	param := ctx.Param("page")
	var ct modBlogArticleWeb.QueryCt
	ct.PageSize = 20
	ct.PageNum = max(strPg.ToInt64(param), 1)
	//
	rt := c.sv.QueryByContext(ctx.Request().Context(), webPg.GetTenantNoEcho(ctx), ct)
	// 模版
	return templatePg.HTMLEcho(ctx, "blog/article_list",
		templatePg.WithDataByResult(rt.SuccessIs(), rt.Data),
		templatePg.WithHtmlObjSet("pageUrl", "page"),
		templatePg.WithSitePage(templatePg.SitePage{
			Title:       "博客",
			Description: "博客",
			Keywords:    "博客",
			SiteName:    "博客",
		}))
}
