package data

import (
	"context"

	"github.com/hongmengzhu/xianfu-blog-go/app/event/ram/service/accountDomainInit"
	"go-spring.org/log"
	_ "go-spring.org/spring/gs"
)

// ZInitAccountAdmin
// @Description: 超管账号初始化
type ZInitAccountAdmin struct {
	sp *accountDomainInit.Sp `autowire:"?"`
}

func (b *ZInitAccountAdmin) Run(ctx context.Context) error {
	log.Infof(ctx, log.TagAppDef, "[init].[超管账号初始化]===================")
	accountDomainInit.NewInitAccount(b.sp).Processor(context.Background())
	return nil
}
