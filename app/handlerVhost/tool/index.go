package tool

import (
	"slices"

	"github.com/hongmengzhu/xianfu-blog-go/app/domain/domainBm/service"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/utils/webPg"
	"github.com/hongmengzhu/xianfu-blog-go/app/handlerVhost"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/cacheService/cacheBm"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/model"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/middleware/authPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/middleware/staterEchoV5"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/templatePg"
	"github.com/labstack/echo/v5"
	"github.com/pangu-2/go-tools/tools/strPg"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(IndexController)).Name("ToolHostIndexController").Export(gs.As[staterEchoV5.RouteRegistrar]())
}

// IndexController 书签首页（独立域名 Tool.test）
type IndexController struct {
	staterEchoV5.RouteRegistrar
	Sp       *authPg.GroupWebVhostMiddlewareSp `autowire:"?"`
	sv       *service.BmToolWebService         `autowire:"?"`
	catCache *cacheBm.CacheToolCategory        `autowire:"?"`
}

// Host 模块绑定域名
func (c *IndexController) Host() string {
	return handlerVhost.Host_tool
}

// RegisterRoutes
//
//	@Description: 注册路由
//	@receiver c
//	@param g
func (c *IndexController) RegisterRoutes(g *echo.Echo) {
	group := g.Group("", authPg.GroupWebVhostMiddlewareEcho(c.Sp))
	group.GET("/", c.Index)
}

func (c *IndexController) Index(ctx *echo.Context) error {
	htmlOption := make([]templatePg.Option, 0)
	htmlOption = append(htmlOption, templatePg.WithHtmlObjSet("pageUrl", "/"))
	htmlOption = append(htmlOption, templatePg.WithHtmlObjSet("pageCategory", make([]model.BaseNodeTree, 0)))
	htmlOption = append(htmlOption, templatePg.WithSitePage(templatePg.SitePage{
		Title:       "分宝呀",
		Description: "AI 工具与资源导航",
		Keywords:    "AI工具,AI导航,AI教程",
		SiteName:    "分宝呀",
	}))
	//
	cxt := ctx.Request().Context()
	tenantNo := webPg.GetTenantNoEcho(ctx)
	nosCategory := make([]*model.BaseNodeDataTree, 0)
	nosIds := make([]string, 0)
	cate, b := c.catCache.GetAllByCache(cxt, tenantNo)
	if b {
		for _, item := range cate {
			if strPg.IsNotBlank(item.No) && !slices.Contains(nosIds, item.No) {
				nosIds = append(nosIds, item.No)
				nosCategory = append(nosCategory, &model.BaseNodeDataTree{
					No:       item.No,
					Name:     item.Name,
					ParentNo: item.ParentNo,
					Code:     item.Code,
					Extend:   item,
				})
			}
		}
		tree := c.catCache.ConvertToTree(cate)
		htmlOption = append(htmlOption, templatePg.WithHtmlObjSet("pageCategory", tree))
		//
		rt := c.sv.QueryListBYCategory(cxt, tenantNo, nosCategory, 20)
		if rt.SuccessIs() {
			log.Info(cxt, log.TagAppDef, "nosCategory", log.Any("count", len(rt.Data)))
			htmlOption = append(htmlOption, templatePg.WithDataByResult(rt.SuccessIs(), rt.Data))
		}
	}
	// 模版
	return templatePg.HTMLEcho(ctx, "tool/index", htmlOption...)
}
