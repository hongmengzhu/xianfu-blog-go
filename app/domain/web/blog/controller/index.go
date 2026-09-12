package controller

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/core/blog/serviceCore"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/web/blog/model/modBlogArticle"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/web/blog/service"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/web/utils/webPg"
	"github.com/hongmengzhu/xianfu-blog-go/app/middleware/authPg"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/tc/cacheTc"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/routerPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/templatePg"
	"github.com/pangu-2/go-tools/tools/strPg"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(IndexController)).Export(gs.As[routerPg.RouteRegistrar]())
}

type IndexController struct {
	routerPg.RouteRegistrar
	Sp       *authPg.GroupWebMiddlewareSp     `autowire:""`
	ca       *cacheTc.TenantDomainCache       `autowire:"?"`
	sv       *service.ArticleService          `autowire:"?"`
	catCache *serviceCore.CoreArticleCategory `autowire:"?"`
}

// RegisterRoutes
//
//	@Description: 注册路由
//	@receiver c
//	@param e
func (c *IndexController) RegisterRoutes(e *gin.Engine) {
	group := e.Group("", authPg.GroupWebMiddleware(c.Sp))
	group.GET("/", c.Index)
	group.GET("/page/:page", c.Page)
}

func (c *IndexController) Index(ctx *gin.Context) {
	var ct modBlogArticle.QueryCt
	ctx.Bind(&ct)
	//
	rt := c.sv.Query(ctx, ct)
	log.Infof(context.Background(), log.TagBizDef, "Data=%+v", rt.Data.Pageable)
	//
	tenantNo := webPg.GetTenantNo(ctx)
	tree, _ := c.catCache.FormatTree(ctx, tenantNo)

	//fmt.Printf("Data: %+v\n", tree)
	// 模版
	templatePg.HTML(ctx, "blog/index",
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

func (c *IndexController) Page(ctx *gin.Context) {
	param := ctx.Param("page")
	var ct modBlogArticle.QueryCt
	ct.PageSize = 20
	ct.PageNum = strPg.ToInt64(param)
	if ct.PageNum < 1 {
		ct.PageNum = 1
	}
	//
	rt := c.sv.Query(ctx, ct)
	// 模版
	templatePg.HTML(ctx, "blog/article_list",
		templatePg.WithDataByResult(rt.SuccessIs(), rt.Data),
		templatePg.WithHtmlObjSet("pageUrl", "page"),
		templatePg.WithSitePage(templatePg.SitePage{
			Title:       "博客",
			Description: "博客",
			Keywords:    "博客",
			SiteName:    "博客",
		}))
}
