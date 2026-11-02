package cronMidPg

import (
	"github.com/go-co-op/gocron/v2"
)

const DefaultKey = "default"

var cron gocron.Scheduler
var cronMap = make(map[string]gocron.Scheduler)

func Get(key string) (gocron.Scheduler, bool) {
	if scheduler, ok := cronMap[key]; ok {
		return scheduler, true
	}
	return nil, false
}
func GetDefault() gocron.Scheduler {
	return cronMap[DefaultKey]
}
func NewCronScheduler(keys []string) {
	for _, item := range keys {
		s, err := gocron.NewScheduler()
		if err != nil {
			// handle error
		} else {
			cronMap[item] = s
		}
	}
}
func NewCronSchedulerDefault() {
	s, err := gocron.NewScheduler()
	if err != nil {
		// handle error
	} else {
		cronMap[DefaultKey] = s
	}
}
func NewJob(jobDefinition gocron.JobDefinition, task gocron.Task, option ...gocron.JobOption) (gocron.Job, error) {
	return cronMap[DefaultKey].NewJob(jobDefinition, task, option...)
}

func Start() {
	cronMap[DefaultKey].Start()
}

func Shutdown() error {
	return cronMap[DefaultKey].Shutdown()
}
