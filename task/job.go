package task

import (
	"context"
	"sync"
	"time"

	"github.com/haozzzzzzzz/go-rapid-development/utils/uerrors"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

type Locker interface {
	LockTask(taskName string, execTimeout time.Duration) (success bool)
	UnlockTask(taskName string) (success bool)
}

type Checker interface {
	Before()
	After(err error)
}

type TaskJob struct {
	// required
	Spec     string
	Schedule cron.Schedule

	TaskName string                                // required
	Handler  func(ctx context.Context) (err error) // required

	RetryTimes    uint32        // 重试次数
	RetryInterval time.Duration // 重试间隔

	ExecTimeout  time.Duration
	CheckerMaker func() (checker Checker)
	Locker       Locker

	Meta sync.Map // custom meta
}

func (m *TaskJob) DoJob() (err error) {
	// context
	var ctx context.Context
	var cancelCtx func(error)

	if m.ExecTimeout <= 0 {
		ctx, cancelCtx = NewBackgroundContext(m.TaskName)
	} else {
		ctx, cancelCtx = NewBackgroundContextWithTimeout(m.TaskName, m.ExecTimeout)
	}

	defer func() {
		cancelCtx(err)
	}()

	// check
	if m.CheckerMaker != nil {
		checker := m.CheckerMaker()
		checker.Before()
		defer func() {
			checker.After(err)
		}()
	}

	defer func() {
		if iRecover := recover(); iRecover != nil {
			err = uerrors.Newf("panic: %s", iRecover)
		}
	}()

	// do job
	startTime := time.Now()
	logrus.Infof("do cron task job start. task_name: %s, start_time: %s", m.TaskName, startTime)
	defer func() {
		endTime := time.Now()
		logrus.Infof("do cron task job end. task_name: %s, end_time: %s, dur: %s, error: %s", m.TaskName, endTime, endTime.Sub(startTime), err)
	}()

	err = m.Handler(ctx)
	if nil != err {
		logrus.Errorf("do Handler failed. error: %s.", err)
		return
	}

	return
}

func (m *TaskJob) Run() {
	curRetryTimes := uint32(0)

	// lock
	if m.Locker != nil {
		if !m.Locker.LockTask(m.TaskName, m.ExecTimeout) {
			return
		}

		defer func() {
			success := m.Locker.UnlockTask(m.TaskName)
			if !success {
				logrus.Errorf("task Locker unlock failed")
			}
		}()

	}

	for {
		err := m.DoJob()
		if err == nil {
			break
		}

		// 出错
		logrus.Errorf("do job failed. task_name: %s, error: %s.", m.TaskName, err)

		// 尝试重试
		if m.RetryTimes == 0 || curRetryTimes >= m.RetryTimes {
			break
		}

		curRetryTimes++

		if m.RetryInterval > 0 {
			time.Sleep(m.RetryInterval)
		}
	}

}
