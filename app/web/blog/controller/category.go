package controller

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/core/blog/serviceCore"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/tc/cacheTc"
	"github.com/hongmengzhu/xianfu-blog-go/app/web/blog/model/modBlogArticle"
	"github.com/hongmengzhu/xianfu-blog-go/app/web/blog/service"
	"github.com/hongmengzhu/xianfu-blog-go/app/web/utils/webPg"
	"github.com/hongmengzhu/xianfu-blog-go/middleware/authPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/routerPg"
	templatePg2 "github.com/hongmengzhu/xianfu-blog-go/pkg/tools/templatePg"
	"github.com/pangu-2/go-tools/tools/strPg"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(CategoryController)).Export(gs.As[routerPg.RouteRegistrar]())
}

type CategoryController struct {
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
func (c *CategoryController) RegisterRoutes(e *gin.Engine) {
	group := e.Group("", authPg.GroupWebMiddleware(c.Sp))
	group.GET("/category/:cat", c.List)
}

func (c *CategoryController) List(ctx *gin.Context) {
	dataIs := false
	var data any
	param := ctx.Param("cat")
	param = strings.TrimSpace(param)
	if strPg.IsNotBlank(param) {
		var ct modBlogArticle.QueryCt
		ct.CategoryQuery = make([]string, 0)
		ct.CategoryQuery = append(ct.CategoryQuery, param)
		//
		//
		rt := c.sv.Query(ctx, ct)
		if rt.SuccessIs() {
			dataIs = true
			data = rt.Data
		}
	}
	//
	tenantNo := webPg.GetTenantNo(ctx)
	tree, _ := c.catCache.FormatTree(ctx, tenantNo)
	// 模版
	templatePg2.HTML(ctx, "blog/article_list",
		templatePg2.WithDataByResult(dataIs, data),
		templatePg2.WithHtmlObjSet("categorys", tree),
		templatePg2.WithHtmlObjSet("pageUrl", "category"),
		templatePg2.WithSitePage(templatePg2.SitePage{
			Title:       "标签",
			Description: "博客",
			Keywords:    "博客",
			SiteName:    "博客",
		}))
}
