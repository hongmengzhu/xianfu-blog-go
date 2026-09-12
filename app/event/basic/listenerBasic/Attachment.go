package listenerBasic

import (
	"context"

	"github.com/farseer-go/eventBus"
	"github.com/farseer-go/fs/core"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/manage/domainBasic/service/attachment"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityBasic"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/repositoryBasic"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/consts/constEventBusPg"
	"go-spring.org/log"
	_ "go-spring.org/spring/gs"
)

// AttachmentListener 附件处理
// @Description:
type AttachmentListener struct {
	dao *repositoryBasic.BasicAttachmentRepository `autowire:"?"`
}

// Run 启动加载
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *AttachmentListener) Run(ctx context.Context) error {
	log.Infof(context.Background(), log.TagAppDef, "eventBus.Register=%+v", constEventBusPg.BasicAttachmentCreate)
	eventBus.RegisterEvent(constEventBusPg.BasicAttachmentCreate).RegisterSubscribe(constEventBusPg.BasicAttachmentCreate, func(message any, _ core.EventArgs) {
		log.Infof(context.Background(), log.TagAppDef, "SchedulerEvent.event=%+v", message)
		dto := message.(entityBasic.BasicAttachmentEntity)
		if len(dto.File) > 0 {
			err := attachment.NewCreate(c.dao, dto).Processor(context.Background())
			if nil != err {
				log.Errorf(ctx, log.TagAppDef, "", err)
			}
			message = nil
		}
	})
	return nil
}
