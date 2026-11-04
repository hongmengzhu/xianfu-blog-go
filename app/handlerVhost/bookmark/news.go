package bookmark

import (
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/domainBlog/service"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/utils/webPg"
	"github.com/hongmengzhu/xianfu-blog-go/app/handlerVhost"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/blog/modBlogArticleWeb"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/cacheService/cacheBlog"
	modulePg "github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/blog/articleModulePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/middleware/authPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/middleware/staterEchoV5"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/templatePg"
	"github.com/labstack/echo/v5"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(NewsController)).Name("BookmarkVhostNewsController").Export(gs.As[staterEchoV5.RouteRegistrar]())
}

// NewsController AI 快讯（独立域名 bookmark.test）
type NewsController struct {
	staterEchoV5.RouteRegistrar
	Sp       *authPg.GroupWebVhostMiddlewareSp   `autowire:"?"`
	sv       *service.BlogArticleWebService      `autowire:"?"`
	catCache *cacheBlog.CacheBlogArticleCategory `autowire:"?"`
}

// Host 模块绑定域名
func (c *NewsController) Host() string {
	return handlerVhost.Host_bookmark
}

// RegisterRoutes
//
//	@Description: 注册路由
//	@receiver c
//	@param g
func (c *NewsController) RegisterRoutes(g *echo.Echo) {
	group := g.Group("", authPg.GroupWebVhostMiddlewareEcho(c.Sp))
	group.GET("/news", c.List)
}

func (c *NewsController) List(ctx *echo.Context) error {
	var ct modBlogArticleWeb.QueryCt
	if err := ctx.Bind(&ct); nil != err {
		return echo.NewHTTPError(400, err.Error())
	}
	// 只查新闻模块
	ct.Module = modulePg.NEW.String()
	// 分页默认值：未传 pageSize 时 PageSize=0 会导致总页溢出且查不出数据
	if ct.PageSize < 1 {
		ct.PageSize = 20
	}
	if ct.PageNum < 1 {
		ct.PageNum = 1
	}
	//
	cxt := ctx.Request().Context()
	tenantNo := webPg.GetTenantNoEcho(ctx)
	rt := c.sv.QueryByContext(cxt, tenantNo, ct)
	log.Info(cxt, log.TagBizDef, "Data", log.Any("Pageable", rt.Data.Pageable))
	//
	tree, _ := c.catCache.FormatTree(cxt, tenantNo)
	// 模版
	return templatePg.HTMLEcho(ctx, "bookmark/news",
		templatePg.WithDataByResult(rt.SuccessIs(), rt.Data),
		templatePg.WithHtmlObjSet("categorys", tree),
		templatePg.WithHtmlObjSet("pageUrl", "/news"),
		templatePg.WithSitePage(templatePg.SitePage{
			Title:       "AI 快讯",
			Description: "每日更新的 AI 行业动态，涵盖新品发布、技术突破与融资并购",
			Keywords:    "AI快讯,AI新闻,行业动态",
			SiteName:    "分宝呀",
		}))
}
