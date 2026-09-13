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
	gs.Provide(new(BasicConfigEventFieldsRepository))

	gs.Provide(new(support.BaseService[BasicConfigEventFieldsRepository]))
}

type BasicConfigEventFieldsRepository struct {
	repositoryPg.BaseRepository[entityBasic.BasicConfigEventFieldsEntity, int64]
}

func (c *BasicConfigEventFieldsRepository) FindAllByModelNo(ctx context.Context, no string) (info []*entityBasic.BasicConfigEventFieldsEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("model_no=?", no).Order("sort asc,create_at").Find(&info)
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

func (c *BasicConfigEventFieldsRepository) DeleteAllByModelNoAndIds(ctx context.Context, no string, ids []string) (info []*entityBasic.BasicConfigEventFieldsEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("id in ?", ids).Where("model_no=?", no).Delete(&c.Entity)
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
func (c *BasicConfigEventFieldsRepository) FindAllByEventNo(ctx context.Context, no string) (info []*entityBasic.BasicConfigEventFieldsEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("event_no=?", no).Order("sort asc,create_at").Find(&info)
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

func (c *BasicConfigEventFieldsRepository) DeleteAllByEventNoAndIds(ctx context.Context, no string, ids []string) (info []*entityBasic.BasicConfigEventFieldsEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("id in ?", ids).Where("event_no=?", no).Delete(&c.Entity)
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

func (c *BasicConfigEventFieldsRepository) FindByEventNoAndFieldIn(ctx context.Context, eventNo string, code []string) (info []*entityBasic.BasicConfigEventFieldsEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("event_no=?", eventNo).Where("field in ?", code).Find(&info)
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
func (c *BasicConfigEventFieldsRepository) FindByEventNoTenantAndFieldIn(ctx context.Context, eventNo string, tenantNo string, code []string) (info []*entityBasic.BasicConfigEventFieldsEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("event_no=?", eventNo).Where("tenant_no=?", tenantNo).Where("field in ?", code).Find(&info)
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
