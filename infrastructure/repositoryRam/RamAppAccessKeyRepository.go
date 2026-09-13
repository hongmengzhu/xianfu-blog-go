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
	gs.Provide(new(RamAppAccessKeyRepository))

	gs.Provide(new(support.BaseService[RamAppAccessKeyRepository]))
}

type RamAppAccessKeyRepository struct {
	repositoryPg.BaseRepository[entityRam.RamAppAccessKeyEntity, int64]
}

func (c *RamAppAccessKeyRepository) FindByTenantNoAndAppNo(ctx context.Context, no, appNo string) (info *entityRam.RamAppAccessKeyEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("tenant_no=?", no).Where("app_no=?", appNo).First(&info)
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

func (c *RamAppAccessKeyRepository) UpdateAllByAppNoAndNoSetState(ctx context.Context, appNo, id string, state int8) (sum int64, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("app_no=?", appNo).Where("id=?", id).Updates(entityRam.RamAppAccessKeyEntity{State: state})
	if tx.Error != nil {
		log.Errorf(ctx, log.TagAppDef, "", tx.Error)
		return 0, false
	}
	if 0 == tx.RowsAffected {
		return 0, false
	}
	return tx.RowsAffected, true
}
