package task

import (
	"context"
	"fmt"
	"service/common/proxy/operations_api"
	"testing"
	"time"

	"github.com/haozzzzzzzz/go-rapid-development/http"
)

func TestNewTaskJob(t *testing.T) {
	operations_api.SetApiUrlPrefix("http://127.0.0.1:18116")
	InitMetrics("test", "1")
	job := NewTaskJob(
		"test",
		"TestNewTaskJob",
		0,
		false,
		false,
		func(ctx context.Context) (err error) {
			for i := 0; i < 10; i++ {
				fmt.Println("tick")
				time.Sleep(1 * time.Second)
			}

			req, err := http.NewRequest("http://127.0.0.1:18116/api/operations_rpc/metrics", ctx, http.NoTimeoutRequestClient)
			if nil != err {
				t.Error(err)
				return
			}

			txt, err := req.GetText()
			if nil != err {
				t.Error(err)
				return
			}

			fmt.Println(txt)
			return
		},
	)

	job.Run()
}
