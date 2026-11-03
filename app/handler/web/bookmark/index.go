package bookmark

import (
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/domainBookmark/service"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/utils/webPg"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/bm/modBmBookmarkWeb"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/cacheService/cacheBookmark"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/model"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/routerPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/middleware/authPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/templatePg"
	"github.com/pangu-2/go-tools/tools/strPg"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(IndexController)).Name("BookmarkIndexController").Export(gs.As[routerPg.RouteRegistrar]())
}

type IndexController struct {
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
func (c *IndexController) RegisterRoutes(e *gin.Engine) {
	group := e.Group("", authPg.GroupWebMiddleware(c.Sp))
	group.GET("/bookmark", c.Index)
	//group.GET("/bookmark/page/:page", c.Page)
}

func (c *IndexController) Index(ctx *gin.Context) {
	htmlOption := make([]templatePg.Option, 0)
	htmlOption = append(htmlOption, templatePg.WithHtmlObjSet("pageUrl", "/bookmark"))
	htmlOption = append(htmlOption, templatePg.WithHtmlObjSet("pageCategory", make([]model.BaseNodeTree, 0)))
	htmlOption = append(htmlOption, templatePg.WithSitePage(templatePg.SitePage{
		Title:       "博客",
		Description: "博客",
		Keywords:    "博客",
		SiteName:    "博客",
	}))
	//
	tenantNo := webPg.GetTenantNo(ctx)
	nosCategory := make([]*model.BaseNodeDataTree, 0)
	nosIds := make([]string, 0)
	cate, b := c.catCache.GetAllByCache(ctx, tenantNo)
	//log.Info(ctx, log.TagAppDef, "nosCategory", log.Any("tenantNo", tenantNo))
	//log.Info(ctx, log.TagAppDef, "cate", log.Any("b", b), log.Any("cate", cate))
	if b {
		//
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
		//
		tree := c.catCache.ConvertToTree(cate)
		htmlOption = append(htmlOption, templatePg.WithHtmlObjSet("pageCategory", tree))
		//
		//log.Info(ctx, log.TagAppDef, "nosCategory", log.Any("nosCategory", nosCategory))
		//
		rt := c.sv.QueryListBYCategory(ctx, tenantNo, nosCategory, 20)
		if rt.SuccessIs() {
			log.Info(ctx, log.TagAppDef, "nosCategory", log.Any("dddddd", rt.Data))
			htmlOption = append(htmlOption, templatePg.WithDataByResult(rt.SuccessIs(), rt.Data))
		}
	}
	// 模版
	templatePg.HTML(ctx, "bookmark/index", htmlOption...)
}

func (c *IndexController) Page(ctx *gin.Context) {
	param := ctx.Param("page")
	var ct modBmBookmarkWeb.QueryCt
	ct.PageSize = 20
	ct.PageNum = max(strPg.ToInt64(param), 1)
	//
	rt := c.sv.Query(ctx, ct)
	// 模版
	templatePg.HTML(ctx, "bookmark/index_page",
		templatePg.WithDataByResult(rt.SuccessIs(), rt.Data),
		templatePg.WithHtmlObjSet("pageUrl", "page"),
		templatePg.WithSitePage(templatePg.SitePage{
			Title:       "博客",
			Description: "博客",
			Keywords:    "博客",
			SiteName:    "博客",
		}))
}
