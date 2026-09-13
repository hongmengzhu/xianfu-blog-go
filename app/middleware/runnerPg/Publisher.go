package runnerPg

import (
	"context"

	"github.com/hongmengzhu/xianfu-blog-go/pkg/event"
)

// Publisher 配置变更发布者。
// 实现 gs.Runner 接口，注入事件总线和所有订阅者，启动后发布事件。
type Publisher struct {
	Bus       event.Bus        `autowire:"?"` // 按类型注入事件总线
	Listeners []event.Listener `autowire:"?"` // 收集所有 Listener Bean
}

// Run 注册所有订阅者后异步发布事件。
// Runner 必须立即返回，不能阻塞。
func (p *Publisher) Run(ctx context.Context) error {
	// 注册所有订阅者到事件总线
	for _, l := range p.Listeners {
		l.Register(p.Bus)
	}
	return nil
}
