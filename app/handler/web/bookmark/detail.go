package bookmark

import (
	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/core/cache/cacheBookmark"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/domainBookmark/service"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/utils/webPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/middleware/authPg"
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
	htmlOption := make([]templatePg.Option, 0)
	htmlOption = append(htmlOption, templatePg.WithHtmlObjSet("pageUrl", "page"))
	htmlOption = append(htmlOption, templatePg.WithSitePage(templatePg.SitePage{
		Title:       "博客",
		Description: "博客",
		Keywords:    "博客",
		SiteName:    "博客",
	}))
	//
	//
	param := ctx.Param("code")
	log.Info(ctx, log.TagBizDef, "Param", log.Any("code", param))
	//
	//
	tenantNo := webPg.GetTenantNo(ctx)
	//
	tree, _ := c.catCache.FormatTree(ctx, tenantNo)
	htmlOption = append(htmlOption, templatePg.WithHtmlObjSet("pageCategory", tree))
	//
	rt := c.sv.Detail(ctx, param)
	htmlOption = append(htmlOption, templatePg.WithDataByResult(rt.SuccessIs(), rt.Data))
	//
	// 模版
	templatePg.HTML(ctx, "bookmark/detail", htmlOption...)
}
