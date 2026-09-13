package repositoryTc

import (
	"context"
	"errors"

	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityTc"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/support"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
	"gorm.io/gorm"
)

func init() {
	gs.Provide(new(TcTenantRepository))

	gs.Provide(new(support.BaseService[TcTenantRepository]))
}

type TcTenantRepository struct {
	repositoryPg.BaseRepository[entityTc.TcTenantEntity, int64]
}

func (c *TcTenantRepository) FindByFounder(ctx context.Context, no string) (info *entityTc.TcTenantEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("founder=?", no).First(&info)
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

func (c *TcTenantRepository) FindByFounderAndNotIdString(ctx context.Context, no string, id string) (info *entityTc.TcTenantEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("founder=?", no).Where("id<>?", id).First(&info)
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

func (c *TcTenantRepository) FindByTenantAndFounder(ctx context.Context, no string) (info *entityTc.TcTenantEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("founder=?", no).First(&info)
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
