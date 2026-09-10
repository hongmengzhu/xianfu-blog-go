package repositoryBlog

import (
	"context"

	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityBlog"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/support"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(BlogArticleStatisticsRepository))

	gs.Provide(new(support.BaseService[BlogArticleStatisticsRepository]))
}

type BlogArticleStatisticsRepository struct {
	repositoryPg.BaseRepository[entityBlog.BlogArticleStatisticsEntity, int64]
}

func (c *BlogArticleStatisticsRepository) FindByArticleNo(ctx context.Context, no string) (info *entityBlog.BlogArticleStatisticsEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("article_no=?", no).First(&info)
	if tx.Error != nil {
		log.Errorf(ctx, log.TagAppDef, "", tx.Error)
		return nil, false
	}
	if 0 == tx.RowsAffected {
		return nil, false
	}
	return info, true
}

func (c *BlogArticleStatisticsRepository) FindAllByArticleNoIn(ctx context.Context, no []string) (info []*entityBlog.BlogArticleStatisticsEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("article_no in ?", no).Find(&info)
	if tx.Error != nil {
		log.Errorf(ctx, log.TagAppDef, "", tx.Error)
		return nil, false
	}
	if 0 == tx.RowsAffected {
		return nil, false
	}
	return info, true
}
