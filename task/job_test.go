package task

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestNewTaskJob(t *testing.T) {
	crontab := New()
	crontab.AddTaskJob(&TaskJob{
		Spec:          "*/20 * * * * *",
		TaskName:      "test",
		RetryTimes:    2,
		RetryInterval: time.Second,
		Handler: func(ctx context.Context) (err error) {
			fmt.Println("tick")
			//err = errors.New("trigger retry")
			return
		},
	})

	crontab.Run()
}
