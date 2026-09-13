package controller

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/core/blog/serviceCore"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/web/blog/model/modBlogArticle"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/web/blog/service"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/web/utils/webPg"
	"github.com/hongmengzhu/xianfu-blog-go/app/middleware/authPg"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/tc/cacheTc"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/routerPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/templatePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"
	"github.com/pangu-2/go-tools/tools/datetimePg"
	"github.com/pangu-2/go-tools/tools/strPg"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(ArchivesController)).Export(gs.As[routerPg.RouteRegistrar]())
}

// ArchivesController 归档
type ArchivesController struct {
	routerPg.RouteRegistrar
	Sp       *authPg.GroupWebMiddlewareSp     `autowire:"?"`
	ca       *cacheTc.TenantDomainCache       `autowire:"?"`
	sv       *service.ArticleService          `autowire:"?"`
	catCache *serviceCore.CoreArticleCategory `autowire:"?"`
}

// RegisterRoutes
//
//	@Description: 注册路由
//	@receiver c
//	@param e
func (c *ArchivesController) RegisterRoutes(e *gin.Engine) {
	group := e.Group("", authPg.GroupWebMiddleware(c.Sp))
	group.GET("/archives/date/:code", c.List)
}

func (c *ArchivesController) List(ctx *gin.Context) {
	dataIs := false
	var data any
	param := ctx.Param("code")
	param = strings.TrimSpace(param)
	if strPg.IsNotBlank(param) {
		flexible, b := datetimePg.IsValidYearMonthFlexible(param)
		if b {
			firstDay := time.Date(flexible.Year(), flexible.Month(), 1, 0, 0, 0, 0, flexible.Location())
			monthEnd := firstDay.AddDate(0, 1, 0).Add(-1 * time.Second)
			var ct modBlogArticle.QueryCt
			ct.CreateAtStart = new(typePg.Time(firstDay))
			ct.CreateAtEnd = new(typePg.Time(monthEnd))
			//
			//
			rt := c.sv.Query(ctx, ct)
			if rt.SuccessIs() {
				dataIs = true
				data = rt.Data
			}
		}
	}
	//
	tenantNo := webPg.GetTenantNo(ctx)
	tree, _ := c.catCache.FormatTree(ctx, tenantNo)
	// 模版
	templatePg.HTML(ctx, "blog/archive",
		templatePg.WithDataByResult(dataIs, data),
		templatePg.WithHtmlObjSet("categorys", tree),
		templatePg.WithHtmlObjSet("subTitle", param),
		templatePg.WithSitePage(templatePg.SitePage{
			Title:       "归档",
			Description: "博客",
			Keywords:    "博客",
			SiteName:    "博客",
		}))
}
