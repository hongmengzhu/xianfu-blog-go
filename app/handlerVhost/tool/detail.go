package tool

import (
	"strings"

	"github.com/hongmengzhu/xianfu-blog-go/app/domain/domainBm/service"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/utils/webPg"
	"github.com/hongmengzhu/xianfu-blog-go/app/handlerVhost"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/cacheService/cacheBm"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/model"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/middleware/authPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/middleware/staterEchoV5"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/templatePg"
	"github.com/labstack/echo/v5"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(DetailController)).Name("ToolVhostDetailController").Export(gs.As[staterEchoV5.RouteRegistrar]())
}

// DetailController 书签详情（独立域名 Tool.test）
type DetailController struct {
	staterEchoV5.RouteRegistrar
	Sp       *authPg.GroupWebVhostMiddlewareSp `autowire:"?"`
	sv       *service.BmToolWebService         `autowire:"?"`
	catCache *cacheBm.CacheToolCategory        `autowire:"?"`
}

// Host 模块绑定域名
func (c *DetailController) Host() string {
	return handlerVhost.Host_tool
}

// RegisterRoutes
//
//	@Description: 注册路由
//	@receiver c
//	@param g
func (c *DetailController) RegisterRoutes(g *echo.Echo) {
	group := g.Group("", authPg.GroupWebVhostMiddlewareEcho(c.Sp))
	group.GET("/tool/:code", c.DetailCode)
	group.GET("/detail/:code", c.DetailNo)
}

// DetailCode 按 id 查详情（与 web 侧 /Tool/tool/:code 行为一致）
func (c *DetailController) DetailCode(ctx *echo.Context) error {
	param := strings.TrimSpace(ctx.Param("code"))
	return c.detail(ctx, param, "/tool/"+param, false)
}

// DetailNo 按编号查详情（与 web 侧 /Tool/detail/:code 行为一致）
func (c *DetailController) DetailNo(ctx *echo.Context) error {
	param := strings.TrimSpace(ctx.Param("code"))
	return c.detail(ctx, param, "/detail/"+param, true)
}

// detail 详情公共实现
//
//	@Description: 渲染详情页
//	@receiver c
//	@param ctx
//	@param param 参数值
//	@param pageUrl 当前页面地址
//	@param isNo true 按编号查询，false 按 id 查询
//	@return error
func (c *DetailController) detail(ctx *echo.Context, param string, pageUrl string, isNo bool) error {
	htmlOption := make([]templatePg.Option, 0)
	htmlOption = append(htmlOption, templatePg.WithHtmlObjSet("pageUrl", pageUrl))
	htmlOption = append(htmlOption, templatePg.WithHtmlObjSet("pageCategory", make([]model.BaseNodeTree, 0)))
	htmlOption = append(htmlOption, templatePg.WithSitePage(templatePg.SitePage{
		Title:       "AI 工具详情",
		Description: "AI 工具详情与使用说明",
		Keywords:    "AI工具,AI导航",
		SiteName:    "分宝呀",
	}))
	//
	cxt := ctx.Request().Context()
	tenantNo := webPg.GetTenantNoEcho(ctx)
	log.Info(cxt, log.TagBizDef, "Param", log.Any("code", param), log.Any("isNo", isNo))
	//
	tree, _ := c.catCache.FormatTree(cxt, tenantNo)
	htmlOption = append(htmlOption, templatePg.WithHtmlObjSet("pageCategory", tree))
	//
	rt := c.sv.DetailByContext(cxt, tenantNo, param, false, isNo)
	htmlOption = append(htmlOption, templatePg.WithDataByResult(rt.SuccessIs(), rt.Data))
	//
	return templatePg.HTMLEcho(ctx, "tool/detail", htmlOption...)
}
