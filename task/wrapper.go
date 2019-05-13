package task

import (
	"context"

	"github.com/haozzzzzzzz/go-rapid-development/utils/uerrors"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
)

// cron: github.com/robfig/cron

type IJobContext interface {
	Before() (ctx context.Context) // before doing job
	After(err error)               // after doing job
}

type Handler func(ctx context.Context) (err error)

func WrapJob(
	jobCtx IJobContext,
	handler Handler,
) cron.FuncJob {
	return func() {
		var err error
		ctx := jobCtx.Before()
		defer func() {
			jobCtx.After(err)
		}()

		defer func() {
			if iRecover := recover(); iRecover != nil {
				err = uerrors.Newf("panic: %s", iRecover)
			}
		}()

		err = handler(ctx)
		if nil != err {
			logrus.Errorf("do handler failed. error: %s.", err)
			return
		}
	}
}
