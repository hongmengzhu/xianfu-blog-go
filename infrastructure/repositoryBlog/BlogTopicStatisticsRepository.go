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
	gs.Provide(new(BlogTopicStatisticsRepository))

	gs.Provide(new(support.BaseService[BlogTopicStatisticsRepository]))
}

type BlogTopicStatisticsRepository struct {
	repositoryPg.BaseRepository[entityBlog.BlogTopicStatisticsEntity, int64]
}

func (c *BlogTopicStatisticsRepository) FindByTopicNo(ctx context.Context, no string) (info *entityBlog.BlogTopicStatisticsEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("topic_no=?", no).First(&info)
	if tx.Error != nil {
		log.Errorf(ctx, log.TagAppDef, "", tx.Error)
		return nil, false
	}
	if 0 == tx.RowsAffected {
		return nil, false
	}
	return info, true
}

func (c *BlogTopicStatisticsRepository) FindAllByTopicNoIn(ctx context.Context, no []string) (info []*entityBlog.BlogTopicStatisticsEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("topic_no in ?", no).Find(&info)
	if tx.Error != nil {
		log.Errorf(ctx, log.TagAppDef, "", tx.Error)
		return nil, false
	}
	if 0 == tx.RowsAffected {
		return nil, false
	}
	return info, true
}
