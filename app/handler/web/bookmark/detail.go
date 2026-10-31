package bookmark

import (
	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/core/cache/cacheBookmark"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/domainBookmark/service"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/utils/webPg"
	"github.com/hongmengzhu/xianfu-blog-go/app/middleware/authPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/routerPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/templatePg"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(DetailController)).Name("BookmarkDetailController").Export(gs.As[routerPg.RouteRegistrar]())
}

type DetailController struct {
	routerPg.RouteRegistrar
	Sp       *authPg.GroupWebMiddlewareSp         `autowire:"?"`
	sv       *service.BmBookmarkWebService        `autowire:"?"`
	catCache *cacheBookmark.CacheBookmarkCategory `autowire:"?"`
}

// RegisterRoutes
//
//	@Description: 注册路由
//	@receiver c
//	@param e
func (c *DetailController) RegisterRoutes(e *gin.Engine) {
	group := e.Group("", authPg.GroupWebMiddleware(c.Sp))
	group.GET("/bookmark/tool/:code", c.DetailCode)
}

func (c *DetailController) DetailCode(ctx *gin.Context) {
	param := ctx.Param("code")
	log.Info(ctx, log.TagBizDef, "Param", log.Any("code", param))
	dataIs := false
	var data any
	rt := c.sv.Detail(ctx, param)
	if rt.SuccessIs() {
		dataIs = true
		data = rt.Data
	}
	//
	tenantNo := webPg.GetTenantNo(ctx)
	tree, _ := c.catCache.FormatTree(ctx, tenantNo)
	//fmt.Printf("Data: %+v\n", data)
	// 模版
	templatePg.HTML(ctx, "bookmark/detail",
		templatePg.WithDataByResult(dataIs, data),
		templatePg.WithHtmlObjSet("pageCategory", tree),
		templatePg.WithHtmlObjSet("pageUrl", "article"),
		templatePg.WithSitePage(templatePg.SitePage{
			Title:       "详情",
			Description: "博客",
			Keywords:    "博客",
			SiteName:    "博客",
		}))
}
