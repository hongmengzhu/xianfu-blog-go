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
	gs.Provide(new(TcTenantDomainRepository))

	gs.Provide(new(support.BaseService[TcTenantDomainRepository]))
}

type TcTenantDomainRepository struct {
	repositoryPg.BaseRepository[entityTc.TcTenantDomainEntity, int64]
}

func (c *TcTenantDomainRepository) FindAllByTenantNo(ctx context.Context, no string) (infos []*entityTc.TcTenantDomainEntity, query bool) {
	tx := c.Db().WithContext(ctx).Where("tenant_no=?", no).Find(&infos)
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
	return infos, true
}

func (c *TcTenantDomainRepository) SetDefaultedByTenantNo(ctx context.Context, def int8, no string) (infos []*entityTc.TcTenantDomainEntity, query bool) {
	tx := c.Db().WithContext(ctx).Where("tenant_no=?", no).Update("defaulted", def)
	if tx.Error != nil {
		// record not found 跳过日志
		if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			log.Errorf(ctx, log.TagAppDef, "err=%+v", tx.Error)
		}
		return nil, false
	}
	if 0 == tx.RowsAffected {
		return nil, false
	}
	return infos, true
}
