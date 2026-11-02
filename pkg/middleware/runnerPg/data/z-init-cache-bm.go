package data

import (
	"context"

	"github.com/hongmengzhu/xianfu-blog-go/app/event/bm/service/bmBookmarkEvent"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/event"
	"go-spring.org/log"
)

// ZInitCacheBm
// @Description: 启动后初始化一些数据
type ZInitCacheBm struct {
	Bus event.Bus `autowire:"?"`
}

func (c *ZInitCacheBm) Run(ctx context.Context) error {
	log.Info(ctx, log.TagAppDef, "[init].[书签.分类.缓存]===================")
	err := bmBookmarkEvent.NewStartInit(c.Bus).Processor(context.Background())
	if err != nil {
		log.Error(ctx, log.TagAppDef, err, "error")
	}
	return nil
}
