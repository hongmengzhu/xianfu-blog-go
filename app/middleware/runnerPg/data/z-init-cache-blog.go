package data

import (
	"context"

	"github.com/hongmengzhu/xianfu-blog-go/app/event/blog/service/articleBlogEvent"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/event"
	"go-spring.org/log"
	_ "go-spring.org/spring/gs"
)

// ZInitCacheBlog
// @Description: 启动后初始化一些数据
type ZInitCacheBlog struct {
	Bus event.Bus `autowire:"?"`
}

func (c *ZInitCacheBlog) Run(ctx context.Context) error {
	log.Infof(ctx, log.TagAppDef, "[init].[博客.分类.缓存]===================")
	err := articleBlogEvent.NewStartInit(c.Bus).Processor(context.Background())
	if err != nil {
		log.Errorf(ctx, log.TagAppDef, "error:%+v", err)
	}
	return nil
}
