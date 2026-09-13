package listenerBasic

import (
	"context"

	"github.com/hongmengzhu/xianfu-blog-go/app/domain/manage/domainBasic/service/attachment"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityBasic"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/repositoryBasic"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/event"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
	_ "go-spring.org/spring/gs"
)

func init() {
	// 注册订阅者 Bean，导出为 event.Listener 供容器自动收集。
	gs.Provide(new(AttachmentListener)).Export(gs.As[event.Listener]())
}

// AttachmentListener 附件处理
// @Description:
type AttachmentListener struct {
	dao *repositoryBasic.BasicAttachmentRepository `autowire:"?"`
}

func (c *AttachmentListener) Register(bus event.Bus) {
	// 创建适配器
	event.Subscribe(bus, func(ctx context.Context, e entityBasic.BasicAttachmentEntity) error {
		if len(e.File) > 0 {
			err := attachment.NewCreate(c.dao, e).Processor(context.Background())
			if nil != err {
				log.Errorf(ctx, log.TagAppDef, "", err)
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
//func (c *AttachmentListener) Run(ctx context.Context) error {
//	log.Infof(context.Background(), log.TagAppDef, "eventBus.Register=%+v", constEventBusPg.BasicAttachmentCreate)
//	eventBus.RegisterEvent(constEventBusPg.BasicAttachmentCreate).RegisterSubscribe(constEventBusPg.BasicAttachmentCreate, func(message any, _ core.EventArgs) {
//		log.Infof(context.Background(), log.TagAppDef, "SchedulerEvent.event=%+v", message)
//		dto := message.(entityBasic.BasicAttachmentEntity)
//		if len(dto.File) > 0 {
//			err := attachment.NewCreate(c.dao, dto).Processor(context.Background())
//			if nil != err {
//				log.Errorf(ctx, log.TagAppDef, "", err)
//			}
//			message = nil
//		}
//	})
//	return nil
//}
