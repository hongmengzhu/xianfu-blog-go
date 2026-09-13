package listenerRam

import (
	"context"

	"github.com/hongmengzhu/xianfu-blog-go/app/event/ram/listenerRam/service"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/repositoryRam"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/event"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/sdk/ram/model/modRamAccount"
	"github.com/pangu-2/go-tools/tools/strPg"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
	_ "go-spring.org/spring/gs"
)

func init() {
	// 注册订阅者 Bean，导出为 event.Listener 供容器自动收集。
	gs.Provide(new(RamListener)).Export(gs.As[event.Listener]())
}

// RamListener ram相关
type RamListener struct {
	acc      *repositoryRam.RamAccountRepository         `autowire:"?"`
	loginLog *repositoryRam.RamAccountLoginLogRepository `autowire:"?"`
	session  *repositoryRam.RamAccountSessionRepository  `autowire:"?"`
}

func (c *RamListener) Register(bus event.Bus) {
	// 创建适配器
	event.Subscribe(bus, func(ctx context.Context, e modRamAccount.LoginLogDto) error {
		if strPg.IsNotBlank(e.Ano) {
			err := service.NewAccountLoginLog(c.acc, c.loginLog, c.session).Processor(context.Background(), e)
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
//func (c *RamListener) Run(ctx context.Context) error {
//	//账号 登录日志
//	eventBus.RegisterEvent(constEventBusPg.RamAccountLoginLog).RegisterSubscribe(constEventBusPg.RamAccountLoginLog, func(message any, _ core.EventArgs) {
//		log.Infof(ctx, log.TagAppDef, "SchedulerEvent[账号.登录日志].event=%+v", message)
//		dto := message.(modRamAccount.LoginLogDto)
//		if strPg.IsNotBlank(dto.Ano) {
//			err := service.NewAccountLoginLog(c.acc, c.loginLog, c.session).Processor(context.Background(), dto)
//			if nil != err {
//				log.Errorf(ctx, log.TagAppDef, "", err)
//			}
//			message = nil
//		}
//	})
//	return nil
//}
