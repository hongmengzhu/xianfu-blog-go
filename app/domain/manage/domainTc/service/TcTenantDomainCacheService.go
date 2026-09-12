package service

import (
	"context"

	"github.com/hongmengzhu/xianfu-blog-go/app/models/tc/cacheTc"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/tc/dtoTenantDomain"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityTc"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/repositoryTc"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/state/enumStatePg"
	"github.com/pangu-2/go-tools/tools/strPg"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(TcTenantDomainCacheService))
}

type TcTenantDomainCacheService struct {
	sv    *repositoryTc.TcTenantDomainRepository `autowire:"?"`
	tenDb *repositoryTc.TcTenantRepository       `autowire:"?"`
	ca    *cacheTc.TenantDomainCache             `autowire:"?"`
}

// InitTenantDomain
//
//	@Description: 初始化加载域名缓存
//	@receiver c
func (c *TcTenantDomainCacheService) InitTenantDomain(ctx context.Context) {
	query := entityTc.TcTenantDomainEntity{
		State: enumStatePg.ENABLE.Index(),
	}
	infos := c.sv.FindAll(ctx, query)
	if len(infos) > 0 {
		idsTenantNo := make([]string, 0)
		mapTenant := make(map[string]*entityTc.TcTenantEntity)
		for _, item := range infos {
			if strPg.IsNotBlank(item.TenantNo) {
				idsTenantNo = append(idsTenantNo, item.TenantNo)
			}
		}
		//
		{
			if len(idsTenantNo) > 0 {
				info, result := c.tenDb.FindAllByNoIn(ctx, idsTenantNo)
				if result {
					for _, item := range info {
						mapTenant[item.No] = item
					}
				}
			}
		}
		//
		for _, item := range infos {
			log.Infof(context.Background(), log.TagAppDef, "domian=tenant,%+v=%+v", item.Code, item.TenantNo)
			c.ca.Domain.Store(item.Code, item.TenantNo)
			dto := dtoTenantDomain.Dto{
				TenantNo: item.TenantNo,
				Domain:   item.Code,
			}
			if strPg.IsNotBlank(item.TenantNo) {
				if info, ok := mapTenant[item.TenantNo]; ok {
					if strPg.IsNotBlank(info.Founder) {
						dto.Ano = info.Founder
					}
				}
			}
			c.ca.DomainData.Store(item.Code, dto)
		}
		infos = nil
		mapTenant = nil
	}
}
