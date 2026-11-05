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
	gs.Provide(new(B72ConfigModelRepository))

	gs.Provide(new(support.BaseService[B72ConfigModelRepository]))
}

type B72ConfigModelRepository struct {
	repositoryPg.BaseRepository[entityB72.B72ConfigModelEntity, int64]
}

func (c *B72ConfigModelRepository) FindByModel(ctx context.Context, code string) (info *entityB72.B72ConfigModelEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("model=?", code).First(&info)
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

func (c *B72ConfigModelRepository) FindByModelAndIdNot(ctx context.Context, code string, id string) (info *entityB72.B72ConfigModelEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("model=?", code).Where("id <>?", id).First(&info)
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

// FindByTenantNoAndModel 租户下按 model 唯一查询。
func (c *B72ConfigModelRepository) FindByTenantNoAndModel(ctx context.Context, tenantNo, model string) (info *entityB72.B72ConfigModelEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("tenant_no=?", tenantNo).Where("model=?", model).First(&info)
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

// FindByTenantNoAndModelAndIdNot 租户下按 model 查重，排除指定 id（更新时使用）。
func (c *B72ConfigModelRepository) FindByTenantNoAndModelAndIdNot(ctx context.Context, tenantNo, model, id string) (info *entityB72.B72ConfigModelEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).
		Where("tenant_no=?", tenantNo).
		Where("model=?", model).
		Where("id <>?", id).
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

// FindByTenantNoAndNo 租户下按模型编号查询。
func (c *B72ConfigModelRepository) FindByTenantNoAndNo(ctx context.Context, tenantNo, no string) (info *entityB72.B72ConfigModelEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("tenant_no=?", tenantNo).Where("no=?", no).First(&info)
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

// FindAllByTenantNo 查询租户下所有模型，用于缓存重建/启动校验。
func (c *B72ConfigModelRepository) FindAllByTenantNo(ctx context.Context, tenantNo string) (info []*entityB72.B72ConfigModelEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("tenant_no=?", tenantNo).Order("create_at desc").Find(&info)
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
