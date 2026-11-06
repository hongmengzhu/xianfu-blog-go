package repositoryTc

import (
	"context"
	"errors"

	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/entityTc"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/support"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
	"gorm.io/gorm"
)

func init() {
	gs.Provide(new(TcTenantDomainModuleRepository))

	gs.Provide(new(support.BaseService[TcTenantDomainModuleRepository]))
}

type TcTenantDomainModuleRepository struct {
	repositoryPg.BaseRepository[entityTc.TcTenantDomainModuleEntity, int64]
}

func (c *TcTenantDomainModuleRepository) FindAllByTenantNo(ctx context.Context, no string) (infos []*entityTc.TcTenantDomainModuleEntity, query bool) {
	tx := c.Db().WithContext(ctx).Where("tenant_no=?", no).Find(&infos)
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
	return infos, true
}

func (c *TcTenantDomainModuleRepository) SetDefaultedByTenantNo(ctx context.Context, def int8, no string) (infos []*entityTc.TcTenantDomainModuleEntity, query bool) {
	tx := c.Db().WithContext(ctx).Where("tenant_no=?", no).Update("defaulted", def)
	if tx.Error != nil {
		// record not found 跳过日志
		if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			log.Error(ctx, log.TagAppDef, tx.Error, "err")
		}
		return nil, false
	}
	if 0 == tx.RowsAffected {
		return nil, false
	}
	return infos, true
}
