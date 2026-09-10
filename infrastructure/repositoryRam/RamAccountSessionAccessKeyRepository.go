package repositoryRam

import (
	"context"

	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityRam"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/support"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(RamAccountSessionAccessKeyRepository))

	gs.Provide(new(support.BaseService[RamAccountSessionAccessKeyRepository]))
}

type RamAccountSessionAccessKeyRepository struct {
	repositoryPg.BaseRepository[entityRam.RamAccountSessionAccessKeyEntity, int64]
}

func (c *RamAccountSessionAccessKeyRepository) FindByAno(ctx context.Context, no string) (info *entityRam.RamAccountSessionAccessKeyEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("ano=?", no).First(&info)
	if tx.Error != nil {
		log.Errorf(ctx, log.TagAppDef, "", tx.Error)
		return nil, false
	}
	if 0 == tx.RowsAffected {
		return nil, false
	}
	return info, true
}

func (c *RamAccountSessionAccessKeyRepository) FindByAnoAndAppNo(ctx context.Context, no, appNo string) (info *entityRam.RamAccountSessionAccessKeyEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("ano=?", no).Where("app_no=?", no).First(&info)
	if tx.Error != nil {
		log.Errorf(ctx, log.TagAppDef, "", tx.Error)
		return nil, false
	}
	if 0 == tx.RowsAffected {
		return nil, false
	}
	return info, true
}

func (c *RamAccountSessionAccessKeyRepository) FindByNoAndState(ctx context.Context, no string) (info *entityRam.RamAccountSessionAccessKeyEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("no=?", no).Where("state=1").First(&info)
	if tx.Error != nil {
		log.Errorf(ctx, log.TagAppDef, "", tx.Error)
		return nil, false
	}
	if 0 == tx.RowsAffected {
		return nil, false
	}
	return info, true
}
func (c *RamAccountSessionAccessKeyRepository) FindByNoAndClientAndState(ctx context.Context, no string, client string) (info *entityRam.RamAccountSessionAccessKeyEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("no=?", no).Where("state=1").Where("client=?", client).First(&info)
	if tx.Error != nil {
		log.Errorf(ctx, log.TagAppDef, "", tx.Error)
		return nil, false
	}
	if 0 == tx.RowsAffected {
		return nil, false
	}
	return info, true
}

func (c *RamAccountSessionAccessKeyRepository) FindByTenantNoAndNoAndState(ctx context.Context, tno, no string) (info *entityRam.RamAccountSessionAccessKeyEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("tenant_no=?", tno).Where("no=?", no).Where("state=1").First(&info)
	if tx.Error != nil {
		log.Errorf(ctx, log.TagAppDef, "", tx.Error)
		return nil, false
	}
	if 0 == tx.RowsAffected {
		return nil, false
	}
	return info, true
}

func (c *RamAccountSessionAccessKeyRepository) FindByTypeDomainAndState(ctx context.Context, domain string) (info *entityRam.RamAccountSessionAccessKeyEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("type_domain = ?", domain).Where("state=1").First(&info)
	if tx.Error != nil {
		log.Errorf(ctx, log.TagAppDef, "", tx.Error)
		return nil, false
	}
	if 0 == tx.RowsAffected {
		return nil, false
	}
	return info, true
}

func (c *RamAccountSessionAccessKeyRepository) FindByTypeDomainInAndState(ctx context.Context, domain []string) (info []*entityRam.RamAccountSessionAccessKeyEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("type_domain in ?", domain).Where("state=1").Find(&info)
	if tx.Error != nil {
		log.Errorf(ctx, log.TagAppDef, "", tx.Error)
		return nil, false
	}
	if 0 == tx.RowsAffected {
		return nil, false
	}
	return info, true
}

func (c *RamAccountSessionAccessKeyRepository) FindByTypeDomainInAndClientAndState(ctx context.Context, domain, client []string) (info []*entityRam.RamAccountSessionAccessKeyEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("type_domain in ?", domain).Where("client in ?", client).Where("state=1").Find(&info)
	if tx.Error != nil {
		log.Errorf(ctx, log.TagAppDef, "", tx.Error)
		return nil, false
	}
	if 0 == tx.RowsAffected {
		return nil, false
	}
	return info, true
}
func (c *RamAccountSessionAccessKeyRepository) FindByTenantNoAndTypeDomainInAndClientAndState(ctx context.Context, tenantNo, domain, client string) (info *entityRam.RamAccountSessionAccessKeyEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("tenant_no = ?", tenantNo).Where("type_domain = ?", domain).Where("client = ?", client).Where("state=1").First(&info)
	if tx.Error != nil {
		log.Errorf(ctx, log.TagAppDef, "", tx.Error)
		return nil, false
	}
	if 0 == tx.RowsAffected {
		return nil, false
	}
	return info, true
}
func (c *RamAccountSessionAccessKeyRepository) FindByTenantNoAndTypeDomainInAndClientAndTypeAndState(ctx context.Context, tenantNo, domain, client, keyType string) (info *entityRam.RamAccountSessionAccessKeyEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("tenant_no = ?", tenantNo).Where("type_domain = ?", domain).Where("client = ?", client).Where("type = ?", keyType).Where("state=1").First(&info)
	if tx.Error != nil {
		log.Errorf(ctx, log.TagAppDef, "", tx.Error)
		return nil, false
	}
	if 0 == tx.RowsAffected {
		return nil, false
	}
	return info, true
}
func (c *RamAccountSessionAccessKeyRepository) DeleteByTypeDomainInAndClientAndState(ctx context.Context, domain, client []string) (result bool) {
	tx := c.DbModel().WithContext(ctx).Where("type_domain in ?", domain).Where("client in ?", client).Where("state=1").Delete(c.Entity)
	if tx.Error != nil {
		log.Errorf(ctx, log.TagAppDef, "", tx.Error)
		return false
	}
	return true
}

func (c *RamAccountSessionAccessKeyRepository) DeleteByTypeDomainAndClientAndState(ctx context.Context, domain string, client []string) (result bool) {
	tx := c.DbModel().WithContext(ctx).Where("type_domain = ?", domain).Where("client in ?", client).Where("state=1").Delete(c.Entity)
	if tx.Error != nil {
		log.Errorf(ctx, log.TagAppDef, "", tx.Error)
		return false
	}
	return true
}

func (c *RamAccountSessionAccessKeyRepository) DeleteByTypeDomainAndClientAndTypeAndState(ctx context.Context, domain, client, keyType string) (result bool) {
	tx := c.DbModel().WithContext(ctx).Where("type_domain = ?", domain).Where("client = ?", client).Where("type=?", keyType).Where("state=1").Delete(c.Entity)
	if tx.Error != nil {
		log.Errorf(ctx, log.TagAppDef, "", tx.Error)
		return false
	}
	return true
}
func (c *RamAccountSessionAccessKeyRepository) DeleteByTenantNoAndTypeDomainAndClientAndTypeAndState(ctx context.Context, tenantNo, domain, client, keyType string) (result bool) {
	tx := c.DbModel().WithContext(ctx).Where("tenant_no=?", tenantNo).Where("type_domain = ?", domain).Where("type_domain = ?", domain).Where("client = ?", client).Where("type=?", keyType).Where("state=1").Delete(c.Entity)
	if tx.Error != nil {
		log.Errorf(ctx, log.TagAppDef, "", tx.Error)
		return false
	}
	return true
}
