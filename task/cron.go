package task

import (
	"errors"

	"github.com/robfig/cron"
)

type CronTab struct {
	cron.Cron
}

func (m *CronTab) AddTaskJob(job *TaskJob) (err error) {
	if job.StrSpec != "" {
		job.Schedule, err = cron.Parse(job.StrSpec)
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
