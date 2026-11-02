package data

import (
	"context"

	"github.com/hongmengzhu/xianfu-blog-go/app/event/ram/service/accountSessionRamEvent"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/cacheService/cacheRam"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/repositoryRam"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/repositoryTc"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/middleware/event"
	"go-spring.org/log"
	_ "go-spring.org/spring/gs"
)

// 加载 密钥缓存
type InitSessionPubPrive struct {
	Bus                  event.Bus                                           `autowire:"?"`
	sessionAk            *repositoryRam.RamAccountSessionAccessKeyRepository `autowire:"?"`
	cacheSessionPubPrive *cacheRam.CacheSessionPubPrive                      `autowire:"?"`
	tenant               *repositoryTc.TcTenantRepository                    `autowire:"?"`
}

func (b *InitSessionPubPrive) Run(ctx context.Context) error {
	log.Info(context.Background(), log.TagAppDef, "[init].[密钥缓存初始化]===================")
	accountSessionRamEvent.NewInitSessionPubPrive(b.sessionAk, b.cacheSessionPubPrive, b.tenant).Processor(context.Background())
	return nil
}
