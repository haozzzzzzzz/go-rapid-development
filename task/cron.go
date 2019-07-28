package task

import (
	"errors"

	"time"

	"github.com/robfig/cron/v3"
)

type CronTab struct {
	cron.Cron
}

var defaultParser = cron.NewParser(
	cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.DowOptional | cron.Descriptor,
)

func (m *CronTab) AddTaskJob(job *TaskJob) (err error) {
	if job.Spec != "" {
		job.Schedule, err = defaultParser.Parse(job.Spec)
		if err != nil {
			return err
		}
	}

	if job.Schedule == nil {
		err = errors.New("schedule should not be nil")
		return
	}

	m.Schedule(job.Schedule, job)
	return
}

func New() *CronTab {
	return NewWithLocation(time.Now().Location())
}

func NewWithLocation(location *time.Location) *CronTab {
	return &CronTab{
		Cron: *cron.New(cron.WithLocation(location)),
	}
}
