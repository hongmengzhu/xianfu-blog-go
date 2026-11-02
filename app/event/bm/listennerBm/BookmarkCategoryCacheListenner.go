package listennerBm

import (
	"context"

	"github.com/hongmengzhu/xianfu-blog-go/app/event/bm/service/bmBookmarkEvent"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/cache/bm/modEventBmBookmarkCategory"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/event"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
)

func init() {
	// 注册订阅者 Bean，导出为 event.Listener 供容器自动收集。
	gs.Provide(new(BookmarkCategoryCacheListener)).Export(gs.As[event.Listener]())
}

// BookmarkCategoryCacheListener 书签分类处理
type BookmarkCategoryCacheListener struct {
	sp *bmBookmarkEvent.Sp `autowire:"?"`
}

func (c *BookmarkCategoryCacheListener) Register(bus event.Bus) {
	// 创建适配器
	event.Subscribe(bus, func(ctx context.Context, e modEventBmBookmarkCategory.CacheDto) error {
		err := bmBookmarkEvent.NewCategoryCache(c.sp, e).Processor(context.Background())
		if nil != err {
			log.Error(ctx, log.TagAppDef, err, "书签.分类.缓存")
		}
		return nil
	})
}
