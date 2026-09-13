package repositoryBasic

import (
	"context"
	"errors"

	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityBasic"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/support"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
	"gorm.io/gorm"
)

func init() {
	gs.Provide(new(BasicConfigListRepository))

	gs.Provide(new(support.BaseService[BasicConfigListRepository]))
}

type BasicConfigListRepository struct {
	repositoryPg.BaseRepository[entityBasic.BasicConfigListEntity, int64]
}

func (c *BasicConfigListRepository) FindByEventNo(ctx context.Context, eventNo string) (info *entityBasic.BasicConfigListEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("event_no=?", eventNo).First(&info)
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

func (c *BasicConfigListRepository) FindByTenantNoAndEventNoAndIdNot(ctx context.Context, tenantNo, eventNo string, id string) (info *entityBasic.BasicConfigListEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("tenant_no=?", tenantNo).Where("event_no=?", eventNo).Where("id != ?", id).First(&info)
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
