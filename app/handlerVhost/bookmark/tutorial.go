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
	gs.Provide(new(TutorialController)).Name("BookmarkVhostTutorialController").Export(gs.As[staterEchoV5.RouteRegistrar]())
}

// TutorialController AI 教程（独立域名 bookmark.test）
type TutorialController struct {
	staterEchoV5.RouteRegistrar
	Sp       *authPg.GroupWebVhostMiddlewareSp   `autowire:"?"`
	sv       *service.BlogArticleWebService      `autowire:"?"`
	catCache *cacheBlog.CacheBlogArticleCategory `autowire:"?"`
}

// Host 模块绑定域名
func (c *TutorialController) Host() string {
	return handlerVhost.Host_bookmark
}

// RegisterRoutes
//
//	@Description: 注册路由
//	@receiver c
//	@param g
func (c *TutorialController) RegisterRoutes(g *echo.Echo) {
	group := g.Group("", authPg.GroupWebVhostMiddlewareEcho(c.Sp))
	group.GET("/tutorial", c.List)
	group.GET("/tutorial/:code", c.Code)
}

func (c *TutorialController) List(ctx *echo.Context) error {
	var ct modBlogArticleWeb.QueryCt
	if err := ctx.Bind(&ct); nil != err {
		return echo.NewHTTPError(400, err.Error())
	}
	// 只查教程模块
	ct.Module = modulePg.TUTORIAL.String()
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
	return templatePg.HTMLEcho(ctx, "bookmark/tutorial",
		templatePg.WithDataByResult(rt.SuccessIs(), rt.Data),
		templatePg.WithHtmlObjSet("categorys", tree),
		templatePg.WithHtmlObjSet("pageUrl", "/tutorial"),
		templatePg.WithSitePage(templatePg.SitePage{
			Title:       "AI 教程资源",
			Description: "精选 AI 工具使用教程、视频课程与学习资料，从入门到进阶系统提升你的 AI 应用能力",
			Keywords:    "AI教程,AI专栏,AI问答,AI名人堂",
			SiteName:    "分宝呀",
		}))
}

func (c *TutorialController) Code(ctx *echo.Context) error {
	param := ctx.Param("code")
	cxt := ctx.Request().Context()
	tenantNo := webPg.GetTenantNoEcho(ctx)
	log.Info(cxt, log.TagBizDef, "Param", log.Any("code", param))
	// 按 code 查详情
	rt := c.sv.DetailByCodeByContext(cxt, tenantNo, param)
	//
	tree, _ := c.catCache.FormatTree(cxt, tenantNo)
	// 模版
	return templatePg.HTMLEcho(ctx, "bookmark/tutorial-detail",
		templatePg.WithDataByResult(rt.SuccessIs(), rt.Data),
		templatePg.WithHtmlObjSet("categorys", tree),
		templatePg.WithHtmlObjSet("pageUrl", "/tutorial/"+param),
		templatePg.WithSitePage(templatePg.SitePage{
			Title:       "AI 教程资源",
			Description: "精选 AI 工具使用教程、视频课程与学习资料，从入门到进阶系统提升你的 AI 应用能力",
			Keywords:    "AI教程,AI专栏,AI问答,AI名人堂",
			SiteName:    "分宝呀",
		}))
}
