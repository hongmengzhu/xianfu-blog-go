package repositoryBm

import (
	"context"
	"errors"

	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityBm"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/support"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
	"gorm.io/gorm"
)

func init() {
	gs.Provide(new(BmBookmarkRepository))

	gs.Provide(new(support.BaseService[BmBookmarkRepository]))
}

type BmBookmarkRepository struct {
	repositoryPg.BaseRepository[entityBm.BmBookmarkEntity, int64]
}

func (c *BmBookmarkRepository) FindAllByUrlSourceMd5(ctx context.Context, code string) (infos []*entityBm.BmBookmarkEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("url_source_md5 = ?", code).Find(&infos)
	if tx.Error != nil {
		// record not found 跳过日志
		if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			log.Error(ctx, log.TagAppDef, tx.Error, "err")
		}
		return
	}
	if 0 == tx.RowsAffected {
		return nil, false
	}
	return infos, true
}

func (c *BmBookmarkRepository) FindAllByUrlSourceMd5In(ctx context.Context, code []string) (infos []*entityBm.BmBookmarkEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("url_source_md5 in ?", code).Find(&infos)
	if tx.Error != nil {
		// record not found 跳过日志
		if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			log.Error(ctx, log.TagAppDef, tx.Error, "err")
		}
		return
	}
	if 0 == tx.RowsAffected {
		return nil, false
	}
	return infos, true
}
