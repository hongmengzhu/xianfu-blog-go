package service

import (
	"github.com/gin-gonic/gin"
	modBlogArticleCategory2 "github.com/hongmengzhu/xianfu-blog-go/app/domain/web/api/model/modBlogArticleCategory"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityBlog"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/repositoryBlog"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/model"
	"github.com/jinzhu/copier"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(ArticleCategoryService))
}

type ArticleCategoryService struct {
	sv *repositoryBlog.BlogArticleCategoryRepository `autowire:"?"`
}

// SelectNodeAllPublic 查询
//
//	@Description:
//	@receiver c
//	@param ct
func (c *ArticleCategoryService) SelectNodeAllPublic(ctx *gin.Context, ct modBlogArticleCategory2.QueryPublicCt) (rt rg.Rs[[]model.BaseNodeNo]) {
	var query entityBlog.BlogArticleCategoryEntity
	copier.Copy(&query, &ct)
	slice := make([]model.BaseNodeNo, 0)
	rt.Data = slice
	infos := c.sv.FindAll(ctx, query)
	if len(infos) > 0 {
		for _, item := range infos {
			var vo modBlogArticleCategory2.Vo
			copier.Copy(&vo, &item)
			code := model.BaseNodeNo{
				Value:    item.No,
				Label:    item.Name,
				ParentNo: item.ParentNo,
				Extend:   vo,
			}
			slice = append(slice, code)
		}
		rt.Data = slice
	}
	return rt.Ok()
}
