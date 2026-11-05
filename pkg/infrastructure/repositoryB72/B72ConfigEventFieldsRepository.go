package repositoryB72

import (
	"context"
	"errors"

	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/entityB72"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/support"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
	"gorm.io/gorm"
)

func init() {
	gs.Provide(new(B72ConfigEventFieldsRepository))

	gs.Provide(new(support.BaseService[B72ConfigEventFieldsRepository]))
}

type B72ConfigEventFieldsRepository struct {
	repositoryPg.BaseRepository[entityB72.B72ConfigEventFieldsEntity, int64]
}

func (c *B72ConfigEventFieldsRepository) FindAllByModelNo(ctx context.Context, no string) (info []*entityB72.B72ConfigEventFieldsEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("model_no=?", no).Order("sort asc,create_at").Find(&info)
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
	return info, true
}

func (c *B72ConfigEventFieldsRepository) DeleteAllByModelNoAndIds(ctx context.Context, no string, ids []string) (info []*entityB72.B72ConfigEventFieldsEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("id in ?", ids).Where("model_no=?", no).Delete(&c.Entity)
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
	return info, true
}
func (c *B72ConfigEventFieldsRepository) FindAllByEventNo(ctx context.Context, no string) (info []*entityB72.B72ConfigEventFieldsEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("event_no=?", no).Order("sort asc,create_at").Find(&info)
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
	return info, true
}

func (c *B72ConfigEventFieldsRepository) DeleteAllByEventNoAndIds(ctx context.Context, no string, ids []string) (info []*entityB72.B72ConfigEventFieldsEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("id in ?", ids).Where("event_no=?", no).Delete(&c.Entity)
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
	return info, true
}

func (c *B72ConfigEventFieldsRepository) FindByEventNoAndFieldIn(ctx context.Context, eventNo string, code []string) (info []*entityB72.B72ConfigEventFieldsEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("event_no=?", eventNo).Where("field in ?", code).Find(&info)
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
	return info, true
}
func (c *B72ConfigEventFieldsRepository) FindByEventNoTenantAndFieldIn(ctx context.Context, eventNo string, tenantNo string, code []string) (info []*entityB72.B72ConfigEventFieldsEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("event_no=?", eventNo).Where("tenant_no=?", tenantNo).Where("field in ?", code).Find(&info)
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
	return info, true
}
