package auditPg

import (
	"context"
	"time"

	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/repositoryBasic"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
)

const (
	cleanInterval = 24 * time.Hour
	cleanBatch    = 5000
)

func init() {
	gs.Provide(new(AuditCleaner)).Export(gs.As[gs.Server]())
}

// AuditCleaner 按保留期清理过期审计日志。
// 保留期可通过 pg.audit.retention-days 配置，默认 180 天。
type AuditCleaner struct {
	repo *repositoryBasic.BasicOperationLogRepository `autowire:"?"`

	RetentionDays int `value:"${pg.audit.retention-days:=180}"`
}

// Run 实现 gs.Server。
func (c *AuditCleaner) Run(ctx context.Context, sig gs.ReadySignal) error {
	<-sig.TriggerAndWait()
	// 启动后延迟 10 分钟执行首次清理，避开启动高峰。
	timer := time.NewTimer(10 * time.Minute)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return nil
	case <-timer.C:
		c.clean(ctx)
	}
	ticker := time.NewTicker(cleanInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			c.clean(ctx)
		}
	}
}

// Stop 实现 gs.Server；清理任务无需额外收尾。
func (c *AuditCleaner) Stop(ctx context.Context) error { return nil }

func (c *AuditCleaner) clean(ctx context.Context) {
	if c.repo == nil {
		return
	}
	days := c.RetentionDays
	if days <= 0 {
		days = 180
	}
	cutoff := time.Now().AddDate(0, 0, -days)
	total, err := c.repo.DeleteBefore(ctx, cutoff, cleanBatch)
	if err != nil {
		log.Error(ctx, log.TagAppDef, err, "审计日志清理失败")
		return
	}
	if total > 0 {
		log.Info(ctx, log.TagAppDef, "审计日志清理完成", log.Any("deleted", total), log.Any("cutoff", cutoff))
	}
}
