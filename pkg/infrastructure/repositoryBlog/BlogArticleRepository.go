package repositoryBlog

import (
	"context"

	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/entityBlog"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/support"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(BlogArticleRepository))

	gs.Provide(new(support.BaseService[BlogArticleRepository]))
}

type BlogArticleRepository struct {
	repositoryPg.BaseRepository[entityBlog.BlogArticleEntity, int64]
}

// FindByCategoryNoIn 一次查询多个分类下的文章，避免按分类逐个 FindAll 的 N+1 查询。
func (c *BlogArticleRepository) FindByCategoryNoIn(ctx context.Context, categoryNos []string) []*entityBlog.BlogArticleEntity {
	if len(categoryNos) == 0 {
		return make([]*entityBlog.BlogArticleEntity, 0)
	}
	list := make([]*entityBlog.BlogArticleEntity, 0)
	tx := c.DbModel().WithContext(ctx).Where("category_no IN ?", categoryNos).Find(&list)
	if tx.Error != nil {
		log.Error(ctx, log.TagAppDef, tx.Error, "BlogArticle.FindByCategoryNoIn failed")
	}
	return list
}
