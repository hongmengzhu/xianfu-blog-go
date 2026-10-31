package bookmark

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/core/cache/cacheBlog"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/domainBlog/service"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/utils/webPg"
	"github.com/hongmengzhu/xianfu-blog-go/app/middleware/authPg"
	"github.com/hongmengzhu/xianfu-blog-go/app/middleware/components/cachePg/cacheDomainPg"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/blog/modBlogArticleWeb"
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
	Sp       *authPg.GroupWebMiddlewareSp        `autowire:"?"`
	ca       *cacheDomainPg.TenantDomainCache    `autowire:"?"`
	sv       *service.BlogArticleWebService      `autowire:"?"`
	catCache *cacheBlog.CacheBlogArticleCategory `autowire:"?"`
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
	var ct modBlogArticleWeb.QueryCt
	ctx.Bind(&ct)
	param := ctx.Param("code")
	log.Info(ctx, log.TagBizDef, "Param", log.Any("code", param))
	//
	rt := c.sv.Query(ctx, ct)
	log.Info(context.Background(), log.TagBizDef, "Data", log.Any("Pageable", rt.Data.Pageable))
	//
	tenantNo := webPg.GetTenantNo(ctx)
	tree, _ := c.catCache.FormatTree(ctx, tenantNo)

	//fmt.Printf("Data: %+v\n", tree)
	// 模版
	templatePg.HTML(ctx, "bookmark/detail",
		templatePg.WithDataByResult(rt.SuccessIs(), rt.Data),
		templatePg.WithHtmlObjSet("categorys", tree),
		templatePg.WithHtmlObjSet("pageUrl", "page"),
		templatePg.WithSitePage(templatePg.SitePage{
			Title:       "博客",
			Description: "博客",
			Keywords:    "博客",
			SiteName:    "博客",
		}))
}
