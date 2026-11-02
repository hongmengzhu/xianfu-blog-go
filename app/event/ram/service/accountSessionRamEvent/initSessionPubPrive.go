package accountSessionRamEvent

import (
	"context"

	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/cacheService/cacheRam"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/consts/constsRam/typeDomainPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/state/clientPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/state/enumStatePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/entityTc"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/repositoryRam"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/repositoryTc"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg/optionsPg"
	"gorm.io/gorm"
)

// InitSessionPubPrive
// @Description: 加载密钥
type InitSessionPubPrive struct {
	sessionAk            *repositoryRam.RamAccountSessionAccessKeyRepository `autowire:"?"`
	cacheSessionPubPrive *cacheRam.CacheSessionPubPrive                      `autowire:"?"`
	tenant               *repositoryTc.TcTenantRepository                    `autowire:"?"`
}

func NewInitSessionPubPrive(
	sessionAk *repositoryRam.RamAccountSessionAccessKeyRepository,
	cacheSessionPubPrive *cacheRam.CacheSessionPubPrive,
	tenant *repositoryTc.TcTenantRepository,
) *InitSessionPubPrive {
	return &InitSessionPubPrive{
		sessionAk:            sessionAk,
		cacheSessionPubPrive: cacheSessionPubPrive,
		tenant:               tenant,
	}
}

func (c *InitSessionPubPrive) Processor(ctx context.Context) error {
	client := clientPg.Browser
	//系统 浏览器
	{
		c.keySystem(ctx, typeDomainPg.System, client)
	}
	//租户
	entity := entityTc.TcTenantEntity{State: enumStatePg.ENABLE.Index()}
	infos := c.tenant.FindAll(ctx, entity, optionsPg.WithCondition(func(db *gorm.DB) *gorm.DB {
		db = db.Order("create_at asc")
		return db
	}))
	if nil != infos {
		for _, item := range infos {
			//登录密钥
			//租户
			c.keyTenant(ctx, typeDomainPg.Manage, client, item.No)
		}
	}
	//

	return nil
}

// 租户
func (c *InitSessionPubPrive) keyTenant(ctx context.Context, domain typeDomainPg.TypeDomain, client clientPg.Client, tenantNo string) {
	c.cacheSessionPubPrive.PaseKeyAccessToken(ctx, domain, client, tenantNo, nil)
	c.cacheSessionPubPrive.PaseKeyRefreshToken(ctx, domain, client, tenantNo, nil)
	//登录密钥
	c.cacheSessionPubPrive.ManageLoginKey(ctx, client, tenantNo)
}

// 系统
func (c *InitSessionPubPrive) keySystem(ctx context.Context,
	domain typeDomainPg.TypeDomain, client clientPg.Client) {
	c.cacheSessionPubPrive.PaseKeyAccessToken(ctx, domain, client, typeDomainPg.System.Code(), nil)
	c.cacheSessionPubPrive.PaseKeyRefreshToken(ctx, domain, client, typeDomainPg.System.Code(), nil)
	//登录密钥
	//重新生成,判断密钥，是否需要保存
	c.cacheSessionPubPrive.SystemLoginKey(ctx, client)
}
