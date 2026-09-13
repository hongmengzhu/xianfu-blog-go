package listenerBasic

import (
	"context"

	"github.com/hongmengzhu/xianfu-blog-go/app/event/basic/model/modEventBasicEvent"
	"github.com/hongmengzhu/xianfu-blog-go/app/event/basic/service/eventBasicEvent"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/event"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
	_ "go-spring.org/spring/gs"
)

func init() {
	// 注册订阅者 Bean，导出为 event.Listener 供容器自动收集。
	gs.Provide(new(EventCacheListener)).Export(gs.As[event.Listener]())
}

type EventCacheListener struct {
	sp *eventBasicEvent.Sp `autowire:"?"`
}

func (c *EventCacheListener) Register(bus event.Bus) {
	log.Infof(context.Background(), log.TagAppDef, "[init].listener.[基础.模型事件.缓存]===================")
	// 创建适配器
	//模型事件
	event.Subscribe(bus, func(ctx context.Context, e modEventBasicEvent.EventDto) error {
		log.Infof(ctx, log.TagAppDef, "listener.[基础.模型事件.缓存]1===================")
		err := eventBasicEvent.NewEventMakeCache(c.sp, e).Processor(context.Background())
		if nil != err {
			log.Errorf(ctx, log.TagAppDef, "基础.模型事件.缓存:%+v", err)
		}
		return nil
	})
	//模型事件字段
	event.Subscribe(bus, func(ctx context.Context, e modEventBasicEvent.FieldDto) error {
		log.Infof(ctx, log.TagAppDef, "listener.[基础.模型事件字段.缓存]2===================")
		err := eventBasicEvent.NewEventFieldMakeCache(c.sp, e).Processor(context.Background())
		if nil != err {
			log.Errorf(ctx, log.TagAppDef, "基础.模型事件字段.缓存:%+v", err)
		}
		return nil
	})
}

// Run 启动加载
//
//	@Description:
//	@receiver c
//	@param ctx
//func (c *EventCacheListener) Run(ctx context.Context) error {
//	log.Infof(ctx, log.TagAppDef, "[init].listener.[基础.模型事件.缓存]===================")
//	//模型事件
//	eventBus.RegisterEvent(constEventBusPg.BasicConfigEventCache).RegisterSubscribe(constEventBusPg.BasicConfigEventCache, func(message any, _ core.EventArgs) {
//		log.Infof(ctx, log.TagAppDef, "listener.[基础.模型事件.缓存]22===================")
//		dto := message.(modEventBasicEvent.EventDto)
//		//log.Infof(ctx, log.TagAppDef,"dto=%+v", dto)
//		err := eventBasicEvent.NewEventMakeCache(c.sp, dto).Processor(context.Background())
//		if nil != err {
//			log.Errorf(ctx, log.TagAppDef, "基础.模型事件.缓存:%+v", err)
//		}
//		message = nil
//	})
//	//模型事件字段
//	eventBus.RegisterEvent(constEventBusPg.BasicConfigEventFieldCache).RegisterSubscribe(constEventBusPg.BasicConfigEventFieldCache, func(message any, _ core.EventArgs) {
//		log.Infof(ctx, log.TagAppDef, "listener.[基础.模型事件字段.缓存]22===================")
//		dto := message.(modEventBasicEvent.FieldDto)
//		//log.Infof(ctx, log.TagAppDef,"dto=%+v", dto)
//		err := eventBasicEvent.NewEventFieldMakeCache(c.sp, dto).Processor(context.Background())
//		if nil != err {
//			log.Errorf(ctx, log.TagAppDef, "基础.模型事件字段.缓存:%+v", err)
//		}
//		message = nil
//	})
//	return nil
//}
