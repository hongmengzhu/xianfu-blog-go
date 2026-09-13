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
	gs.Provide(new(BlogTopicRelationRepository))

	gs.Provide(new(support.BaseService[BlogTopicRelationRepository]))
}

type BlogTopicRelationRepository struct {
	repositoryPg.BaseRepository[entityBlog.BlogTopicRelationEntity, int64]
}

func (c *BlogTopicRelationRepository) FindAllByTopicNo(ctx context.Context, no string) (info []*entityBlog.BlogTopicRelationEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("topic_no=?", no).Find(&info)
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

func (c *BlogTopicRelationRepository) FindByTopicNoAndArticleNo(ctx context.Context, topicNo, no string) (info *entityBlog.BlogTopicRelationEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("topic_no=?", topicNo).Where("article_no=?", no).First(&info)
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
func (c *BlogTopicRelationRepository) FindAllByArticleNo(ctx context.Context, no string) (info []*entityBlog.BlogTopicRelationEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("article_no=?", no).Find(&info)
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
