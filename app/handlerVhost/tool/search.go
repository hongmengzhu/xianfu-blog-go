package tool

import (
	"html/template"
	"net/url"
	"strings"

	"github.com/hongmengzhu/xianfu-blog-go/app/domain/domainBm/service"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/utils/webPg"
	"github.com/hongmengzhu/xianfu-blog-go/app/handlerVhost"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/bm/modBmToolWeb"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/cacheService/cacheBm"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/middleware/authPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/middleware/staterEchoV5"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/templatePg"
	"github.com/labstack/echo/v5"
	"github.com/pangu-2/go-tools/tools/strPg"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(SearchController)).Name("ToolVhostSearchController").Export(gs.As[staterEchoV5.RouteRegistrar]())
}

// SearchController 工具搜索（独立域名 Tool.test）
type SearchController struct {
	staterEchoV5.RouteRegistrar
	Sp       *authPg.GroupWebVhostMiddlewareSp `autowire:"?"`
	sv       *service.BmToolWebService         `autowire:"?"`
	catCache *cacheBm.CacheToolCategory        `autowire:"?"`
}

// Host 模块绑定域名
func (c *SearchController) Host() string {
	return handlerVhost.Host_tool
}

// RegisterRoutes
//
//	@Description: 注册路由
//	@receiver c
//	@param g
func (c *SearchController) RegisterRoutes(g *echo.Echo) {
	group := g.Group("", authPg.GroupWebVhostMiddlewareEcho(c.Sp))
	group.GET("/search", c.List)
}

// List 搜索结果（分页）
//
// search 只按名称模糊匹配，QueryByContext 使用 ct.Wd 生成 name like 条件；
// 分页链接需要带上 q，因此额外输出 pageQuery 供 block_pagination 拼接。
func (c *SearchController) List(ctx *echo.Context) error {
	q := strings.TrimSpace(ctx.QueryParam("q"))
	//
	var ct modBmToolWeb.QueryCt
	ct.Wd = q
	ct.PageSize = 20
	ct.PageNum = max(strPg.ToInt64(ctx.QueryParam("pageNum")), 1)
	//
	cxt := ctx.Request().Context()
	tenantNo := webPg.GetTenantNoEcho(ctx)
	rt := c.sv.QueryByContext(cxt, tenantNo, ct)
	log.Info(cxt, log.TagBizDef, "Data", log.Any("q", q), log.Any("Pageable", rt.Data.Pageable))
	//
	tree, _ := c.catCache.FormatTree(cxt, tenantNo)
	//
	pageQuery := template.URL("")
	title := "搜索工具"
	description := "按名称搜索已收录的 AI 工具"
	if strPg.IsNotBlank(q) {
		// 用 template.URL 交给模板，避免 html/template 把查询串再次转义成 %25/%3d
		pageQuery = template.URL("q=" + url.QueryEscape(q) + "&")
		title = "搜索 " + q
		description = "与“" + q + "”相关的 AI 工具"
	}
	// 模版
	return templatePg.HTMLEcho(ctx, "tool/search",
		templatePg.WithDataByResult(rt.SuccessIs(), rt.Data),
		templatePg.WithHtmlObjSet("pageCategory", tree),
		templatePg.WithHtmlObjSet("pageUrl", "/search"),
		templatePg.WithHtmlObjSet("searchQuery", q),
		templatePg.WithHtmlObjSet("pageQuery", pageQuery),
		templatePg.WithSitePage(templatePg.SitePage{
			Title:       title,
			Description: description,
			Keywords:    "AI工具搜索,AI导航,工具检索",
			SiteName:    "分宝呀",
		}))
}
