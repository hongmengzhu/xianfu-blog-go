package controller

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/core/blog/serviceCore"
	"github.com/hongmengzhu/xianfu-blog-go/app/web/blog/model/modBlogArticle"
	"github.com/hongmengzhu/xianfu-blog-go/app/web/blog/service"
	"github.com/hongmengzhu/xianfu-blog-go/app/web/utils/webPg"
	"github.com/hongmengzhu/xianfu-blog-go/middleware/authPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/routerPg"
	templatePg2 "github.com/hongmengzhu/xianfu-blog-go/pkg/tools/templatePg"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(ArticleController)).Export(gs.As[routerPg.RouteRegistrar]())
}

type ArticleController struct {
	routerPg.RouteRegistrar
	Sp       *authPg.GroupWebMiddlewareSp     `autowire:"?"`
	sv       *service.ArticleService          `autowire:"?"`
	catCache *serviceCore.CoreArticleCategory `autowire:"?"`
}

// RegisterRoutes
//
//	@Description: 注册路由
//	@receiver c
//	@param e
func (c *ArticleController) RegisterRoutes(e *gin.Engine) {
	group := e.Group("", authPg.GroupWebMiddleware(c.Sp))
	group.GET("/article/search", c.List)
	group.GET("/article/:id", c.Detail)
}

func (c *ArticleController) Detail(ctx *gin.Context) {
	dataIs := false
	var data any
	param := ctx.Param("id")
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
	templatePg2.HTML(ctx, "blog/detail",
		templatePg2.WithDataByResult(dataIs, data),
		templatePg2.WithHtmlObjSet("categorys", tree),
		templatePg2.WithHtmlObjSet("pageUrl", "article"),
		templatePg2.WithSitePage(templatePg2.SitePage{
			Title:       "详情",
			Description: "博客",
			Keywords:    "博客",
			SiteName:    "博客",
		}))
}

func (c *ArticleController) List(ctx *gin.Context) {
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
	templatePg2.HTML(ctx, "blog/article_list",
		templatePg2.WithDataByResult(rt.SuccessIs(), rt.Data),
		templatePg2.WithHtmlObjSet("categorys", tree),
		templatePg2.WithHtmlObjSet("pageUrl", "/article/search"),
		templatePg2.WithSitePage(templatePg2.SitePage{
			Title:       "博客",
			Description: "博客",
			Keywords:    "博客",
			SiteName:    "博客",
		}))
}
