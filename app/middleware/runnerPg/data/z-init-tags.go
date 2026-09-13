package data

import (
	"context"

	"github.com/hongmengzhu/xianfu-blog-go/app/event/basic/service/tagsBasicEvent"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/event"
	"go-spring.org/log"
	_ "go-spring.org/spring/gs"
)

// ZInitTagsCache
// @Description: 启动后初始化一些数据
type ZInitTagsCache struct {
	Bus event.Bus          `autowire:"?"`
	sp  *tagsBasicEvent.Sp `autowire:"?"`
}

func (c *ZInitTagsCache) Run(ctx context.Context) error {
	log.Infof(ctx, log.TagAppDef, "[init].[标签缓存]===================")
	err := tagsBasicEvent.NewStartInit(c.sp).Processor(context.Background())
	if err != nil {
		log.Errorf(ctx, log.TagAppDef, "error:", err)
	}
	return nil
}
