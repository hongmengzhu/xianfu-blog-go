package listenerBasic

import (
	"context"

	"github.com/hongmengzhu/xianfu-blog-go/app/event/basic/model/modEventBasicTags"
	"github.com/hongmengzhu/xianfu-blog-go/app/event/basic/service/tagsBasicEvent"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/event"
	"github.com/pangu-2/go-tools/tools/strPg"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
	_ "go-spring.org/spring/gs"
)

func init() {
	// 注册订阅者 Bean，导出为 event.Listener 供容器自动收集。
	gs.Provide(new(TagsListener)).Export(gs.As[event.Listener]())
}

// TagsListener 标签处理
// @Description:
type TagsListener struct {
	sp *tagsBasicEvent.Sp `autowire:"?"`
}

func (c *TagsListener) Register(bus event.Bus) {
	// 创建适配器
	event.Subscribe(bus, func(ctx context.Context, e modEventBasicTags.TagsRelation) error {
		if strPg.IsNotBlank(e.Category) {
			err := tagsBasicEvent.NewSaveByCategory(c.sp, e).Processor()
			if nil != err {
				log.Errorf(ctx, log.TagAppDef, "err:=%+v", err)
			}
		}
		return nil
	})
}

// Run 启动加载
//
//	@Description:
//	@receiver c
//	@param ctx
//func (c *TagsListener) Run(ctx context.Context) error {
//	log.Infof(context.Background(), log.TagAppDef, "eventBus.Register=%+v", constEventBusPg.BlogArticle)
//	//博客文章
//	eventBus.RegisterEvent(constEventBusPg.BlogArticle).RegisterSubscribe(constEventBusPg.BlogArticle, func(message any, _ core.EventArgs) {
//		log.Infof(ctx, log.TagAppDef, "SchedulerEvent.event=%+v", message)
//		dto := message.(modEventBasicTags.TagsRelation)
//		if strPg.IsNotBlank(dto.Category) {
//			err := tagsBasicEvent.NewSaveByCategory(c.sp, dto).Processor()
//			if nil != err {
//				log.Errorf(ctx, log.TagAppDef, "", err)
//			}
//			message = nil
//		}
//	})
//	return nil
//}
