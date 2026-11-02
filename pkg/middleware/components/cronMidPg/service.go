package cronMidPg

import (
	"time"

	"github.com/go-co-op/gocron/v2"
)

// Cron
// @Description: 工具
type Cron struct {
	scheduler gocron.Scheduler
}

func NewCron(scheduler gocron.Scheduler) *Cron {
	return &Cron{
		scheduler: scheduler,
	}
}

// NewJobDelay
//
//	@Description: 延迟 几个秒/分/小时 后执行
//	@param add
//	@param task
//	@param option
//	@return error
func (c *Cron) NewJobDelay(add time.Duration, task gocron.Task, option ...gocron.JobOption) (gocron.Job, error) {
	return c.scheduler.NewJob(gocron.OneTimeJob(
		gocron.OneTimeJobStartDateTime(time.Now().Add(add)),
	), task, option...)
}
