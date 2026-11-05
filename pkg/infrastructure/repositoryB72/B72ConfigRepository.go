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
	gs.Provide(new(B72ConfigRepository))

	gs.Provide(new(support.BaseService[B72ConfigRepository]))
}

type B72ConfigRepository struct {
	repositoryPg.BaseRepository[entityB72.B72ConfigEntity, int64]
}

func (c *B72ConfigRepository) FindByTypeDomainTenantAndField(ctx context.Context, typeDomain string, tenantNo string, code string) (info *entityB72.B72ConfigEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("type_domain=?", typeDomain).Where("tenant_no=?", tenantNo).Where("field=?", code).First(&info)
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
func (c *B72ConfigRepository) FindByEventNoAndFieldIn(ctx context.Context, eventNo string, code []string) (info []*entityB72.B72ConfigEntity, result bool) {
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
func (c *B72ConfigRepository) FindByEventNo(ctx context.Context, eventNo string) (info []*entityB72.B72ConfigEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("event_no=?", eventNo).Find(&info)
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
func (c *B72ConfigRepository) FindByEventNoTenantAndFieldIn(ctx context.Context, eventNo string, tenantNo string, code []string) (info []*entityB72.B72ConfigEntity, result bool) {
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
func (c *B72ConfigRepository) UpdateByTenantEventNoAndNoAndValue(ctx context.Context, tenantNo string, eventNo string, no, value string) (info []*entityB72.B72ConfigEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("tenant_no=?", tenantNo).Where("event_no=?", eventNo).Where("no = ?", no).Update("value", value)
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

func (c *B72ConfigRepository) UpdateByEventNoAndNoAndValue(ctx context.Context, eventNo string, no, value string) (result bool) {
	tx := c.DbModel().WithContext(ctx).Where("event_no=?", eventNo).Where("no = ?", no).Update("value", value)
	if tx.Error != nil {
		// record not found 跳过日志
		if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			log.Error(ctx, log.TagAppDef, tx.Error, "err")
		}
		return false
	}
	if 0 == tx.RowsAffected {
		return false
	}
	return true
}

// FindAllByModelNo 查询模型下所有配置（含停用），按 sort、创建时间排序。
func (c *B72ConfigRepository) FindAllByModelNo(ctx context.Context, modelNo string) (info []*entityB72.B72ConfigEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("model_no=?", modelNo).Order("sort asc,create_at").Find(&info)
	if tx.Error != nil {
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

// FindAllEnabledByTenantNoAndModel 查询租户 + model 下所有启用配置。
func (c *B72ConfigRepository) FindAllEnabledByTenantNoAndModel(ctx context.Context, tenantNo, model string) (info []*entityB72.B72ConfigEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).
		Where("tenant_no=?", tenantNo).
		Where("model=?", model).
		Where("state=?", 1).
		Order("sort asc,create_at").
		Find(&info)
	if tx.Error != nil {
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

// FindAllEnabled 查询所有启用配置，用于缓存重建。
func (c *B72ConfigRepository) FindAllEnabled(ctx context.Context) (info []*entityB72.B72ConfigEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("state=?", 1).Order("sort asc,create_at").Find(&info)
	if tx.Error != nil {
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

// FindByTenantNoAndModelNoAndField 查询租户 + 模型编号 + 字段下唯一配置。
func (c *B72ConfigRepository) FindByTenantNoAndModelNoAndField(ctx context.Context, tenantNo, modelNo, field string) (info *entityB72.B72ConfigEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).
		Where("tenant_no=?", tenantNo).
		Where("model_no=?", modelNo).
		Where("field=?", field).
		First(&info)
	if tx.Error != nil {
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

// FindByTenantNoAndModelAndField 查询租户 + model + 字段下唯一配置。
func (c *B72ConfigRepository) FindByTenantNoAndModelAndField(ctx context.Context, tenantNo, model, field string) (info *entityB72.B72ConfigEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).
		Where("tenant_no=?", tenantNo).
		Where("model=?", model).
		Where("field=?", field).
		First(&info)
	if tx.Error != nil {
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

// DeleteAllByModelNo 物理删除模型下所有配置，用于模型删除/类型切换。
func (c *B72ConfigRepository) DeleteAllByModelNo(ctx context.Context, modelNo string) error {
	tx := c.DbModel().WithContext(ctx).Where("model_no=?", modelNo).Delete(&c.Entity)
	return tx.Error
}
