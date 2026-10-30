package blog

import (
	"strings"

	"github.com/hongmengzhu/xianfu-blog-go/app/core/blog/serviceCore"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/host"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/web/blog/model/modBlogArticle"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/web/blog/service"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/web/utils/webPg"
	"github.com/hongmengzhu/xianfu-blog-go/app/middleware/authPg"
	"github.com/hongmengzhu/xianfu-blog-go/app/middleware/staterEchoV5"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/routerPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/templatePg"
	"github.com/labstack/echo/v5"
	"github.com/pangu-2/go-tools/tools/strPg"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(CategoryController)).Name("BlogHostCategoryController").Export(gs.As[routerPg.HostRouteRegistrar]())
}

// CategoryController 分类
type CategoryController struct {
	staterEchoV5.RouteRegistrar
	Sp       *authPg.GroupWebMiddlewareSp     `autowire:"?"`
	sv       *service.ArticleService          `autowire:"?"`
	catCache *serviceCore.CoreArticleCategory `autowire:"?"`
}

// Hosts 模块绑定域名
func (c *CategoryController) Hosts() string {
	return host.Host_blog
}

// RegisterRoutes
//
//	@Description: 注册路由
//	@receiver c
//	@param g
func (c *CategoryController) RegisterRoutes(g *echo.Group) {
	group := g.Group("", authPg.GroupWebMiddlewareEcho(c.Sp))
	group.GET("/category/:cat", c.List)
}

func (c *CategoryController) List(ctx *echo.Context) error {
	dataIs := false
	var data any
	param := ctx.Param("cat")
	param = strings.TrimSpace(param)
	tenantNo := webPg.GetTenantNoEcho(ctx)
	if strPg.IsNotBlank(param) {
		var ct modBlogArticle.QueryCt
		ct.CategoryQuery = make([]string, 0)
		ct.CategoryQuery = append(ct.CategoryQuery, param)
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
		templatePg.WithHtmlObjSet("pageUrl", "category"),
		templatePg.WithSitePage(templatePg.SitePage{
			Title:       "标签",
			Description: "博客",
			Keywords:    "博客",
			SiteName:    "博客",
		}))
}
