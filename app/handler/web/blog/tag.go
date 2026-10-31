package blog

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/core/cache/cacheBlog"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/domainBlog/service"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/utils/webPg"
	"github.com/hongmengzhu/xianfu-blog-go/app/middleware/authPg"
	"github.com/hongmengzhu/xianfu-blog-go/app/middleware/components/cachePg/cacheDomainPg"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/blog/modBlogArticleWeb"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/routerPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/templatePg"
	"github.com/pangu-2/go-tools/tools/strPg"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(TagController)).Export(gs.As[routerPg.RouteRegistrar]())
}

type TagController struct {
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
func (c *TagController) RegisterRoutes(e *gin.Engine) {
	group := e.Group("", authPg.GroupWebMiddleware(c.Sp))
	group.GET("/tag/:tag", c.List)
}

func (c *TagController) List(ctx *gin.Context) {
	dataIs := false
	var data any
	param := ctx.Param("tag")
	param = strings.TrimSpace(param)
	if strPg.IsNotBlank(param) {
		var ct modBlogArticleWeb.QueryCt
		ct.TagsQuery = make([]string, 0)
		ct.TagsQuery = append(ct.TagsQuery, param)
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
	templatePg.HTML(ctx, "blog/article_list",
		templatePg.WithDataByResult(dataIs, data),
		templatePg.WithHtmlObjSet("categorys", tree),
		templatePg.WithHtmlObjSet("pageUrl", "tag"),
		templatePg.WithSitePage(templatePg.SitePage{
			Title:       "标签",
			Description: "博客",
			Keywords:    "博客",
			SiteName:    "博客",
		}))
}
