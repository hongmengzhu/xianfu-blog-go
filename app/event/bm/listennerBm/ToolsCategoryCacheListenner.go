package listennerBm

import (
	"context"

	"github.com/hongmengzhu/xianfu-blog-go/app/event/bm/service/bmToolEvent"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/event/bm/modEventBmToolCategory"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/middleware/event"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
)

func init() {
	// 注册订阅者 Bean，导出为 event.Listener 供容器自动收集。
	gs.Provide(new(ToolCategoryCacheListener)).Export(gs.As[event.Listener]())
}

// ToolCategoryCacheListener 工具-分类处理
type ToolCategoryCacheListener struct {
	sp *bmToolEvent.Sp `autowire:"?"`
}

func (c *ToolCategoryCacheListener) Register(bus event.Bus) {
	// 创建适配器
	event.Subscribe(bus, func(ctx context.Context, e modEventBmToolCategory.CacheDto) error {
		err := bmToolEvent.NewCategoryCache(c.sp, e).Processor(context.Background())
		if nil != err {
			log.Error(ctx, log.TagAppDef, err, "工具.分类.缓存")
		}
		return nil
	})
}
