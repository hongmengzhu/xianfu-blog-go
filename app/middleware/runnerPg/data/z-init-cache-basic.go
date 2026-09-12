package data

import (
	"context"

	"github.com/hongmengzhu/xianfu-blog-go/app/event/basic/service/eventBasicEvent"
	"github.com/hongmengzhu/xianfu-blog-go/app/event/basic/service/eventBasicRules"
	"go-spring.org/log"
	_ "go-spring.org/spring/gs"
)

// ZInitCacheBasic
// @Description: 启动后初始化一些数据
type ZInitCacheBasic struct {
}

func (c *ZInitCacheBasic) Run(ctx context.Context) error {
	log.Infof(ctx, log.TagAppDef, "[init].[启动初始化.基础.缓存]===================")
	{
		err := eventBasicEvent.NewStartInit().Processor(context.Background())
		if err != nil {
			log.Errorf(ctx, log.TagAppDef, "error:", err)
		}
	}
	{
		err := eventBasicRules.NewStartInit().Processor(context.Background())
		if err != nil {
			log.Errorf(ctx, log.TagAppDef, "error:", err)
		}
	}
	return nil
}
