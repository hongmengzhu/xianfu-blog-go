package data

import (
	"context"

	"github.com/hongmengzhu/xianfu-blog-go/app/manage/domainTc/service"
	"go-spring.org/log"
	_ "go-spring.org/spring/gs"
)

// InitTenantDomain
// @Description: 初始化租户域名
type InitTenantDomain struct {
	domain *service.TcTenantDomainCacheService `autowire:"?"`
}

func (b *InitTenantDomain) Run(ctx context.Context) error {
	log.Infof(context.Background(), log.TagAppDef, "初始化 => 域名与租户的关系")
	b.domain.InitTenantDomain(context.Background())
	return nil
}
