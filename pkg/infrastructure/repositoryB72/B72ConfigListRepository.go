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
	gs.Provide(new(B72ConfigListRepository))

	gs.Provide(new(support.BaseService[B72ConfigListRepository]))
}

type B72ConfigListRepository struct {
	repositoryPg.BaseRepository[entityB72.B72ConfigListEntity, int64]
}

func (c *B72ConfigListRepository) FindByEventNo(ctx context.Context, eventNo string) (info *entityB72.B72ConfigListEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("event_no=?", eventNo).First(&info)
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

// FindByTenantNoAndModelNo 按租户 + 模型编号查询（配置模型保存后同步的配置列表行）
func (c *B72ConfigListRepository) FindByTenantNoAndModelNo(ctx context.Context, tenantNo, modelNo string) (info *entityB72.B72ConfigListEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("tenant_no=?", tenantNo).Where("model_no=?", modelNo).First(&info)
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

// DeleteAllByModelNo 按模型编号删除配置列表行
func (c *B72ConfigListRepository) DeleteAllByModelNo(ctx context.Context, modelNo string) error {
	if modelNo == "" {
		return nil
	}
	tx := c.DbModel().WithContext(ctx).Where("model_no=?", modelNo).Delete(&entityB72.B72ConfigListEntity{})
	if tx.Error != nil {
		log.Error(ctx, log.TagAppDef, tx.Error, "err")
	}
	return tx.Error
}

func (c *B72ConfigListRepository) FindByTenantNoAndEventNoAndIdNot(ctx context.Context, tenantNo, eventNo string, id string) (info *entityB72.B72ConfigListEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("tenant_no=?", tenantNo).Where("event_no=?", eventNo).Where("id != ?", id).First(&info)
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
