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
	gs.Provide(new(BasicConfigRepository))

	gs.Provide(new(support.BaseService[BasicConfigRepository]))
}

type BasicConfigRepository struct {
	repositoryPg.BaseRepository[entityBasic.BasicConfigEntity, int64]
}

func (c *BasicConfigRepository) FindByTypeDomainTenantAndField(ctx context.Context, typeDomain string, tenantNo string, code string) (info *entityBasic.BasicConfigEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("type_domain=?", typeDomain).Where("tenant_no=?", tenantNo).Where("field=?", code).First(&info)
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
func (c *BasicConfigRepository) FindByEventNoAndFieldIn(ctx context.Context, eventNo string, code []string) (info []*entityBasic.BasicConfigEntity, result bool) {
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
func (c *BasicConfigRepository) FindByEventNo(ctx context.Context, eventNo string) (info []*entityBasic.BasicConfigEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("event_no=?", eventNo).Find(&info)
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
func (c *BasicConfigRepository) FindByEventNoTenantAndFieldIn(ctx context.Context, eventNo string, tenantNo string, code []string) (info []*entityBasic.BasicConfigEntity, result bool) {
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
func (c *BasicConfigRepository) UpdateByTenantEventNoAndNoAndValue(ctx context.Context, tenantNo string, eventNo string, no, value string) (info []*entityBasic.BasicConfigEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("tenant_no=?", tenantNo).Where("event_no=?", eventNo).Where("no = ?", no).Update("value", value)
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

func (c *BasicConfigRepository) UpdateByEventNoAndNoAndValue(ctx context.Context, eventNo string, no, value string) (result bool) {
	tx := c.DbModel().WithContext(ctx).Where("event_no=?", eventNo).Where("no = ?", no).Update("value", value)
	if tx.Error != nil {
		// record not found 跳过日志
		if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			log.Errorf(ctx, log.TagAppDef, "err=%+v", tx.Error)
		}
		return false
	}
	if 0 == tx.RowsAffected {
		return false
	}
	return true
}
