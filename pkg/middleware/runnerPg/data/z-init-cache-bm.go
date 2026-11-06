package data

import (
	"context"

	"github.com/hongmengzhu/xianfu-blog-go/app/event/bm/service/bmBookmarkEvent"
	"github.com/hongmengzhu/xianfu-blog-go/app/event/bm/service/bmToolEvent"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/middleware/event"
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
		log.Error(ctx, log.TagAppDef, err, "书签.分类.缓存.error")
	}
	log.Info(ctx, log.TagAppDef, "[init].[工具.分类.缓存]===================")
	err = bmToolEvent.NewStartInit(c.Bus).Processor(context.Background())
	if err != nil {
		log.Error(ctx, log.TagAppDef, err, "工具.分类.缓存.error")
	}
	return nil
}
