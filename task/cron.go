package task

import (
	"errors"

	"time"

	"github.com/robfig/cron"
)

type CronTab struct {
	cron.Cron
}

func (m *CronTab) AddTaskJob(job *TaskJob) (err error) {
	if job.Spec != "" {
		job.Schedule, err = cron.Parse(job.Spec)
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
		Cron: *cron.NewWithLocation(location),
	}
}
