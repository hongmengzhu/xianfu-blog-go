package repositoryTc

import (
	"context"
	"errors"

	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/state/enumStatePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/entityTc"
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
			log.Error(ctx, log.TagAppDef, tx.Error, "err")
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
			log.Error(ctx, log.TagAppDef, tx.Error, "err")
		}
		return nil, false
	}
	if 0 == tx.RowsAffected {
		return nil, false
	}
	return infos, true
}

// FindByCodeGlobal 全局（跨租户）按域名 code 查询，用于唯一性校验。
// 直接使用 Db() 绕过多租户 scope：域名必须全局唯一，不能被租户过滤掉。
func (c *TcTenantDomainRepository) FindByCodeGlobal(ctx context.Context, code string) (*entityTc.TcTenantDomainEntity, bool) {
	var info *entityTc.TcTenantDomainEntity
	tx := c.Db().WithContext(ctx).Where("code=?", code).First(&info)
	if tx.Error != nil {
		if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			log.Error(ctx, log.TagAppDef, tx.Error, "err")
		}
		return nil, false
	}
	if tx.RowsAffected == 0 {
		return nil, false
	}
	return info, true
}

// FindAllEnabledForCache 缓存初始化/刷新专用：读取全部启用域名，不受租户 scope 影响。
func (c *TcTenantDomainRepository) FindAllEnabledForCache(ctx context.Context) []*entityTc.TcTenantDomainEntity {
	infos := make([]*entityTc.TcTenantDomainEntity, 0)
	tx := c.Db().WithContext(ctx).Where("state=?", enumStatePg.ENABLE.Index()).Find(&infos)
	if tx.Error != nil {
		if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			log.Error(ctx, log.TagAppDef, tx.Error, "err")
		}
		return infos
	}
	return infos
}
