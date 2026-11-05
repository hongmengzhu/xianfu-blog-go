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
	gs.Provide(new(B72ConfigModelFieldsRepository))

	gs.Provide(new(support.BaseService[B72ConfigModelFieldsRepository]))
}

type B72ConfigModelFieldsRepository struct {
	repositoryPg.BaseRepository[entityB72.B72ConfigModelFieldsEntity, int64]
}

func (c *B72ConfigModelFieldsRepository) FindAllByModelNo(ctx context.Context, no string) (info []*entityB72.B72ConfigModelFieldsEntity, result bool) {
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

func (c *B72ConfigModelFieldsRepository) DeleteAllByModelNoAndIds(ctx context.Context, no string, ids []string) (info []*entityB72.B72ConfigModelFieldsEntity, result bool) {
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

func (c *B72ConfigModelFieldsRepository) FindByEventNoAndFieldIn(ctx context.Context, eventNo string, code []string) (info []*entityB72.B72ConfigModelFieldsEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("model_no=?", eventNo).Where("field in ?", code).Find(&info)
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
func (c *B72ConfigModelFieldsRepository) FindByEventNoTenantAndFieldIn(ctx context.Context, eventNo string, tenantNo string, code []string) (info []*entityB72.B72ConfigModelFieldsEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("model_no=?", eventNo).Where("tenant_no=?", tenantNo).Where("field in ?", code).Find(&info)
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

// FindByModelNoAndField 查询模型下指定字段。
func (c *B72ConfigModelFieldsRepository) FindByModelNoAndField(ctx context.Context, modelNo, field string) (info *entityB72.B72ConfigModelFieldsEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("model_no=?", modelNo).Where("field=?", field).First(&info)
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

// DeleteAllByModelNo 物理删除模型下所有字段。
func (c *B72ConfigModelFieldsRepository) DeleteAllByModelNo(ctx context.Context, modelNo string) error {
	tx := c.DbModel().WithContext(ctx).Where("model_no=?", modelNo).Delete(&c.Entity)
	return tx.Error
}
