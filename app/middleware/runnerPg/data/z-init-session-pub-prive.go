package data

import (
	"context"

	"github.com/hongmengzhu/xianfu-blog-go/app/core/cache/cacheRam"
	"github.com/hongmengzhu/xianfu-blog-go/app/event/ram/service/accountSessionRamEvent"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/repositoryRam"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/repositoryTc"
	"go-spring.org/log"
	_ "go-spring.org/spring/gs"
)

// 加载 密钥缓存
type InitSessionPubPrive struct {
	sessionAk            *repositoryRam.RamAccountSessionAccessKeyRepository `autowire:"?"`
	cacheSessionPubPrive *cacheRam.CacheSessionPubPrive                      `autowire:"?"`
	tenant               *repositoryTc.TcTenantRepository                    `autowire:"?"`
}

func (b *InitSessionPubPrive) Run(ctx context.Context) error {
	log.Infof(context.Background(), log.TagAppDef, "[init].[密钥缓存初始化]===================")
	accountSessionRamEvent.NewInitSessionPubPrive(b.sessionAk, b.cacheSessionPubPrive, b.tenant).Processor(context.Background())
	return nil
}
