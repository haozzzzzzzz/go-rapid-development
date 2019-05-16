package task

import (
	"context"
	"fmt"
	"testing"
)

func TestNewTaskJob(t *testing.T) {
	crontab := New()
	crontab.AddTaskJob(&TaskJob{
		Spec:     "*/2 * * * * *",
		TaskName: "test",
		Handler: func(ctx context.Context) (err error) {
			fmt.Println("tick")
			return
		},
	})

	crontab.Run()
}
