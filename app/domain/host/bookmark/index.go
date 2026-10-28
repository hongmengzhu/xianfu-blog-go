package bookmark

import (
	"net/http"

	"github.com/hongmengzhu/xianfu-blog-go/app/core/blog/serviceCore"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/web/blog/model/modBlogArticle"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/web/blog/service"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/web/utils/webPg"
	"github.com/hongmengzhu/xianfu-blog-go/app/middleware/authPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/routerPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/templatePg"
	"github.com/labstack/echo/v4"
	"github.com/pangu-2/go-tools/tools/strPg"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(IndexController)).Name("BookmarkHostIndexController").Export(gs.As[routerPg.HostRouteRegistrar]())
}

// IndexController 书签首页
type IndexController struct {
	Sp       *authPg.GroupWebMiddlewareSp     `autowire:"?"`
	sv       *service.ArticleService          `autowire:"?"`
	catCache *serviceCore.CoreArticleCategory `autowire:"?"`
}

// Hosts 模块绑定域名
func (c *IndexController) Hosts() []string {
	return bookmarkHosts
}

// RegisterRoutes
//
//	@Description: 注册路由
//	@receiver c
//	@param g
func (c *IndexController) RegisterRoutes(g *echo.Group) {
	group := g.Group("", authPg.GroupWebMiddlewareEcho(c.Sp))
	group.GET("/", c.Index)
	group.GET("/page/:page", c.Page)
}

func (c *IndexController) Index(ctx echo.Context) error {
	var ct modBlogArticle.QueryCt
	if err := ctx.Bind(&ct); nil != err {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	//
	//tenantNo := webPg.GetTenantNoEcho(ctx)
	//rt := c.sv.QueryByContext(ctx.Request().Context(), tenantNo, ct)
	//log.Info(context.Background(), log.TagBizDef, "Data", log.Any("Pageable", rt.Data.Pageable))
	////
	//tree, _ := c.catCache.FormatTree(ctx.Request().Context(), tenantNo)
	return echo.NewHTTPError(http.StatusOK, "kkkkkk")
	// 模版
	//return templatePg.HTMLEcho(ctx, "bookmark/index",
	//	//templatePg.WithDataByResult(rt.SuccessIs(), rt.Data),
	//	//templatePg.WithHtmlObjSet("categorys", tree),
	//	templatePg.WithHtmlObjSet("pageUrl", "page"),
	//	templatePg.WithSitePage(templatePg.SitePage{
	//		Title:       "博客",
	//		Description: "博客",
	//		Keywords:    "博客",
	//		SiteName:    "博客",
	//	}))
}

func (c *IndexController) Page(ctx echo.Context) error {
	param := ctx.Param("page")
	var ct modBlogArticle.QueryCt
	ct.PageSize = 20
	ct.PageNum = strPg.ToInt64(param)
	if ct.PageNum < 1 {
		ct.PageNum = 1
	}
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
