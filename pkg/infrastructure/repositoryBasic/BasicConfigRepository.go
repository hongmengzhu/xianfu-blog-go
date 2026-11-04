package repositoryBasic

import (
	"context"
	"errors"

	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/entityBasic"
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
			log.Error(ctx, log.TagAppDef, tx.Error, "err")
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
			log.Error(ctx, log.TagAppDef, tx.Error, "err")
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
			log.Error(ctx, log.TagAppDef, tx.Error, "err")
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
			log.Error(ctx, log.TagAppDef, tx.Error, "err")
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
			log.Error(ctx, log.TagAppDef, tx.Error, "err")
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
func (c *BasicConfigRepository) FindAllByModelNo(ctx context.Context, modelNo string) (info []*entityBasic.BasicConfigEntity, result bool) {
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
func (c *BasicConfigRepository) FindAllEnabledByTenantNoAndModel(ctx context.Context, tenantNo, model string) (info []*entityBasic.BasicConfigEntity, result bool) {
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
func (c *BasicConfigRepository) FindAllEnabled(ctx context.Context) (info []*entityBasic.BasicConfigEntity, result bool) {
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
func (c *BasicConfigRepository) FindByTenantNoAndModelNoAndField(ctx context.Context, tenantNo, modelNo, field string) (info *entityBasic.BasicConfigEntity, result bool) {
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
func (c *BasicConfigRepository) FindByTenantNoAndModelAndField(ctx context.Context, tenantNo, model, field string) (info *entityBasic.BasicConfigEntity, result bool) {
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
func (c *BasicConfigRepository) DeleteAllByModelNo(ctx context.Context, modelNo string) error {
	tx := c.DbModel().WithContext(ctx).Where("model_no=?", modelNo).Delete(&c.Entity)
	return tx.Error
}
