package blog

import (
	"strings"
	"time"

	"github.com/hongmengzhu/xianfu-blog-go/app/core/blog/serviceCore"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/domainBlog/service"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/utils/webPg"
	"github.com/hongmengzhu/xianfu-blog-go/app/handler/vhost"
	"github.com/hongmengzhu/xianfu-blog-go/app/middleware/authPg"
	"github.com/hongmengzhu/xianfu-blog-go/app/middleware/staterEchoV5"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/blog/modBlogArticleWeb"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/templatePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"
	"github.com/labstack/echo/v5"
	"github.com/pangu-2/go-tools/tools/datetimePg"
	"github.com/pangu-2/go-tools/tools/strPg"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(ArchivesController)).Name("BlogHostArchivesController").Export(gs.As[staterEchoV5.RouteRegistrar]())
}

// ArchivesController 归档
type ArchivesController struct {
	staterEchoV5.RouteRegistrar
	Sp       *authPg.GroupWebVhostMiddlewareSp `autowire:"?"`
	sv       *service.BlogArticleWebService    `autowire:"?"`
	catCache *serviceCore.CoreArticleCategory  `autowire:"?"`
}

// Hosts 模块绑定域名
func (c *ArchivesController) Host() string {
	return vhost.Host_blog
}

// RegisterRoutes
//
//	@Description: 注册路由
//	@receiver c
//	@param g
func (c *ArchivesController) RegisterRoutes(g *echo.Echo) {
	group := g.Group("", authPg.GroupWebVhostMiddlewareEcho(c.Sp))
	group.GET("/archives/date/:code", c.List)
}

func (c *ArchivesController) List(ctx *echo.Context) error {
	dataIs := false
	var data any
	param := ctx.Param("code")
	param = strings.TrimSpace(param)
	tenantNo := webPg.GetTenantNoEcho(ctx)
	if strPg.IsNotBlank(param) {
		flexible, b := datetimePg.IsValidYearMonthFlexible(param)
		if b {
			firstDay := time.Date(flexible.Year(), flexible.Month(), 1, 0, 0, 0, 0, flexible.Location())
			monthEnd := firstDay.AddDate(0, 1, 0).Add(-1 * time.Second)
			var ct modBlogArticleWeb.QueryCt
			ct.CreateAtStart = new(typePg.Time(firstDay))
			ct.CreateAtEnd = new(typePg.Time(monthEnd))
			//
			rt := c.sv.QueryByContext(ctx.Request().Context(), tenantNo, ct)
			if rt.SuccessIs() {
				dataIs = true
				data = rt.Data
			}
		}
	}
	//
	tree, _ := c.catCache.FormatTree(ctx.Request().Context(), tenantNo)
	// 模版
	return templatePg.HTMLEcho(ctx, "blog/archive",
		templatePg.WithDataByResult(dataIs, data),
		templatePg.WithHtmlObjSet("categorys", tree),
		templatePg.WithHtmlObjSet("subTitle", param),
		templatePg.WithSitePage(templatePg.SitePage{
			Title:       "归档",
			Description: "博客",
			Keywords:    "博客",
			SiteName:    "博客",
		}))
}
