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
	gs.Provide(new(BlogCollectRepository))

	gs.Provide(new(support.BaseService[BlogCollectRepository]))
}

type BlogCollectRepository struct {
	repositoryPg.BaseRepository[entityBlog.BlogCollectEntity, int64]
}

func (c *BlogCollectRepository) FindAllByUrlSourceMd5(ctx context.Context, code string) (infos []*entityBlog.BlogCollectEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("url_source_md5 = ?", code).Find(&infos)
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
	return infos, true
}

func (c *BlogCollectRepository) FindAllByUrlSourceMd5In(ctx context.Context, code []string) (infos []*entityBlog.BlogCollectEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("url_source_md5 in ?", code).Find(&infos)
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
	return infos, true
}
