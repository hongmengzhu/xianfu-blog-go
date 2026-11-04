package bookmark

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/domainBlog/service"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/utils/webPg"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/blog/modBlogArticleWeb"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/cacheService/cacheBlog"
	modulePg "github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/blog/articleModulePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/routerPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/middleware/authPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/middleware/components/cachePg/cacheDomainPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/templatePg"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(NewsController)).Name("BookmarkNewsController").Export(gs.As[routerPg.RouteRegistrar]())
}

type NewsController struct {
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
func (c *NewsController) RegisterRoutes(e *gin.Engine) {
	group := e.Group("", authPg.GroupWebMiddleware(c.Sp))
	group.GET("/bookmark/news", c.List)
}

func (c *NewsController) List(ctx *gin.Context) {
	var ct modBlogArticleWeb.QueryCt
	ctx.Bind(&ct)
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
	rt := c.sv.Query(ctx, ct)
	log.Info(context.Background(), log.TagBizDef, "Data", log.Any("Pageable", rt.Data.Pageable))
	//
	tenantNo := webPg.GetTenantNo(ctx)
	tree, _ := c.catCache.FormatTree(ctx, tenantNo)

	//fmt.Printf("Data: %+v\n", tree)
	// 模版
	templatePg.HTML(ctx, "bookmark/news",
		templatePg.WithDataByResult(rt.SuccessIs(), rt.Data),
		templatePg.WithHtmlObjSet("categorys", tree),
		templatePg.WithHtmlObjSet("pageUrl", "page"),
		templatePg.WithSitePage(templatePg.SitePage{
			Title:       "AI 快讯",
			Description: "每日更新的 AI 行业动态，涵盖新品发布、技术突破与融资并购",
			Keywords:    "AI快讯,AI新闻,行业动态",
			SiteName:    "博客",
		}))
}
