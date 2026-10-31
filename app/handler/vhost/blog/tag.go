package blog

import (
	"strings"

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
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(TagController)).Name("BlogHostTagController").Export(gs.As[staterEchoV5.RouteRegistrar]())
}

// TagController 标签
type TagController struct {
	staterEchoV5.RouteRegistrar
	Sp       *authPg.GroupWebVhostMiddlewareSp   `autowire:"?"`
	sv       *service.BlogArticleWebService      `autowire:"?"`
	catCache *cacheBlog.CacheBlogArticleCategory `autowire:"?"`
}

// Hosts 模块绑定域名
func (c *TagController) Host() string {
	return vhost.Host_blog
}

// RegisterRoutes
//
//	@Description: 注册路由
//	@receiver c
//	@param g
func (c *TagController) RegisterRoutes(g *echo.Echo) {
	group := g.Group("", authPg.GroupWebVhostMiddlewareEcho(c.Sp))
	group.GET("/tag/:tag", c.List)
}

func (c *TagController) List(ctx *echo.Context) error {
	dataIs := false
	var data any
	param := ctx.Param("tag")
	param = strings.TrimSpace(param)
	tenantNo := webPg.GetTenantNoEcho(ctx)
	if strPg.IsNotBlank(param) {
		var ct modBlogArticleWeb.QueryCt
		ct.TagsQuery = make([]string, 0)
		ct.TagsQuery = append(ct.TagsQuery, param)
		//
		rt := c.sv.QueryByContext(ctx.Request().Context(), tenantNo, ct)
		if rt.SuccessIs() {
			dataIs = true
			data = rt.Data
		}
	}
	//
	tree, _ := c.catCache.FormatTree(ctx.Request().Context(), tenantNo)
	// 模版
	return templatePg.HTMLEcho(ctx, "blog/article_list",
		templatePg.WithDataByResult(dataIs, data),
		templatePg.WithHtmlObjSet("categorys", tree),
		templatePg.WithHtmlObjSet("pageUrl", "tag"),
		templatePg.WithSitePage(templatePg.SitePage{
			Title:       "标签",
			Description: "博客",
			Keywords:    "博客",
			SiteName:    "博客",
		}))
}
