package data

import (
	"context"

	"github.com/hongmengzhu/xianfu-blog-go/app/event/blog/service/articleBlogEvent"
	"go-spring.org/log"
	_ "go-spring.org/spring/gs"
)

// ZInitCacheBlog
// @Description: 启动后初始化一些数据
type ZInitCacheBlog struct {
}

func (c *ZInitCacheBlog) Run(ctx context.Context) error {
	log.Infof(ctx, log.TagAppDef, "[init].[博客.分类.缓存]===================")
	err := articleBlogEvent.NewStartInit().Processor(context.Background())
	if err != nil {
		log.Errorf(ctx, log.TagAppDef, "error:", err)
	}
	return nil
}
