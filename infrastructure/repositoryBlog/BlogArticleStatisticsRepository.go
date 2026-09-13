package repositoryBlog

import (
	"context"
	"errors"

	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityBlog"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/support"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
	"gorm.io/gorm"
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
		// record not found 跳过日志
		if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			log.Errorf(ctx, log.TagAppDef, "err=%+v", tx.Error)
		}
		return
	}
	if 0 == tx.RowsAffected {
		return nil, false
	}
	return info, true
}

func (c *BlogArticleStatisticsRepository) FindAllByArticleNoIn(ctx context.Context, no []string) (info []*entityBlog.BlogArticleStatisticsEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("article_no in ?", no).Find(&info)
	if tx.Error != nil {
		// record not found 跳过日志
		if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			log.Errorf(ctx, log.TagAppDef, "err=%+v", tx.Error)
		}
		return
	}
	if 0 == tx.RowsAffected {
		return nil, false
	}
	return info, true
}
