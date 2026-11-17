package service

import (
	"context"

	"github.com/hongmengzhu/xianfu-blog-go/app/models/tc/dtoTenantDomain"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/state/enumStatePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/entityTc"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/repositoryTc"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/middleware/components/cachePg/cacheDomainPg"
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
	ca    *cacheDomainPg.TenantDomainCache       `autowire:"?"`
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
	log.Debug(ctx, log.TagAppDef, "init tenant domain cache.findAll", func() []log.Field {
		return []log.Field{
			log.Int("count", len(infos)),
		}
	})
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
			log.Info(context.Background(), log.TagAppDef, "domian=tenant,%+v=%+v", log.Any("Code", item.Code), log.Any("TenantNo", item.TenantNo))
			c.ca.Domain.Store(item.Code, item.TenantNo)
			dto := dtoTenantDomain.Dto{
				TenantNo: item.TenantNo,
				Domain:   item.Code,
				Theme:    item.Theme,
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
		log.Debug(ctx, log.TagAppDef, "init tenant domain cache", func() []log.Field {
			return []log.Field{
				log.Int("count", len(infos)),
			}
		})
		log.Info(ctx, log.TagAppDef, "init tenant domain cache", log.Int("count", len(infos)))
		infos = nil
		mapTenant = nil
	}
}

// Reload 重新加载域名缓存。
//
//	@Description: 域名增删改后调用，保证缓存与数据库一致（原实现只在启动时加载一次）。
//	@receiver c
//	@param ctx
func (c *TcTenantDomainCacheService) Reload(ctx context.Context) {
	// 清空旧缓存，避免已删除/禁用/改绑的域名继续解析
	c.clear()
	infos := c.sv.FindAllEnabledForCache(ctx)
	if len(infos) == 0 {
		log.Info(ctx, log.TagAppDef, "tenant domain cache reloaded", log.Int("count", 0))
		return
	}
	idsTenantNo := make([]string, 0)
	mapTenant := make(map[string]*entityTc.TcTenantEntity)
	for _, item := range infos {
		if strPg.IsNotBlank(item.TenantNo) {
			idsTenantNo = append(idsTenantNo, item.TenantNo)
		}
	}
	if len(idsTenantNo) > 0 {
		if info, result := c.tenDb.FindAllByNoIn(ctx, idsTenantNo); result {
			for _, item := range info {
				mapTenant[item.No] = item
			}
		}
	}
	for _, item := range infos {
		if strPg.IsBlank(item.Code) || strPg.IsBlank(item.TenantNo) {
			continue
		}
		c.ca.Domain.Store(item.Code, item.TenantNo)
		dto := dtoTenantDomain.Dto{
			TenantNo: item.TenantNo,
			Domain:   item.Code,
			Theme:    item.Theme,
		}
		if info, ok := mapTenant[item.TenantNo]; ok && strPg.IsNotBlank(info.Founder) {
			dto.Ano = info.Founder
		}
		c.ca.DomainData.Store(item.Code, dto)
	}
	log.Info(ctx, log.TagAppDef, "tenant domain cache reloaded", log.Int("count", len(infos)))
}

// clear 清空域名缓存
func (c *TcTenantDomainCacheService) clear() {
	keys := make([]string, 0)
	c.ca.Domain.Range(func(key string, _ string) bool {
		keys = append(keys, key)
		return true
	})
	for _, key := range keys {
		c.ca.Domain.Delete(key)
	}
	keysData := make([]string, 0)
	c.ca.DomainData.Range(func(key string, _ dtoTenantDomain.Dto) bool {
		keysData = append(keysData, key)
		return true
	})
	for _, key := range keysData {
		c.ca.DomainData.Delete(key)
	}
}
