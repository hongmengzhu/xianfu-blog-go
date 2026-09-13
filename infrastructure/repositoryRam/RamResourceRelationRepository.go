package repositoryRam

import (
	"context"
	"errors"

	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityRam"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/support"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
	"gorm.io/gorm"
)

func init() {
	gs.Provide(new(RamResourceRelationRepository))

	gs.Provide(new(support.BaseService[RamResourceRelationRepository]))
}

type RamResourceRelationRepository struct {
	repositoryPg.BaseRepository[entityRam.RamResourceRelationEntity, int64]
	//
}

func (c *RamResourceRelationRepository) DeleteByAuthorityId(ctx context.Context, code int64) {
	c.DbModel().WithContext(ctx).Where("authority_id=?", code).Delete(&entityRam.RamResourceRelationEntity{})
}

func (c *RamResourceRelationRepository) FindByMark(ctx context.Context, code string) (info *entityRam.RamResourceRelationEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("mark=?", code).First(&info)
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
