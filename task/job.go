package task

import (
	"context"
	"time"

	"github.com/haozzzzzzzz/go-rapid-development/aws/xray"
	"github.com/haozzzzzzzz/go-rapid-development/utils/uerrors"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
)

type Locker interface {
	LockTask(taskName string, expire time.Duration) (success bool)
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

	ExecTimeout  time.Duration
	CheckerMaker func() (checker Checker)
	Locker       Locker
}

func (m *TaskJob) DoJob() (err error) {
	// lock
	if m.Locker != nil {
		if !m.Locker.LockTask(m.TaskName, m.ExecTimeout) {
			return
		}

		defer func() {
			success := m.Locker.UnlockTask(m.TaskName)
			if !success {
				logrus.Errorf("task Locker unlock failed. error: %s.", err)
			}
		}()

	}

	// context
	var ctx context.Context
	var cancelCtx func(err error)
	if m.ExecTimeout <= 0 {
		ctx, _, cancelCtx = xray.NewBackgroundContext(m.TaskName)
	} else {
		ctx, _, cancelCtx = xray.NewBackgroundContextWithTimeout(m.TaskName, m.ExecTimeout)
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
	err = m.Handler(ctx)
	if nil != err {
		logrus.Errorf("do Handler failed. error: %s.", err)
		return
	}

	return
}

func (m *TaskJob) Run() {
	err := m.DoJob()
	if nil != err {
		logrus.Errorf("do job failed. task_name: %s, error: %s.", m.TaskName, err)
		return
	}
}
