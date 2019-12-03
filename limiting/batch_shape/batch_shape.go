/**
将输入的元素，批量吐出
*/
package batch_shape

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"time"
)

type BatchShapeHandler func(items []interface{}) (err error)
type BatchShape struct {
	Size    int           // 批量数值
	Timeout time.Duration // 超过一定时间没达到BatchSize，则全部吐出
	C       chan interface{}
	Handler BatchShapeHandler
}

func NewBatchShape(
	size int,
	timeout time.Duration,
	handler BatchShapeHandler,
) (bs *BatchShape) {
	bs = &BatchShape{
		Size:    size,
		Timeout: timeout,
		Handler: handler,
		C:       make(chan interface{}, size),
	}
	return
}

func (m *BatchShape) Start() {
	go func() {
		// 设置定时器
		items := make([]interface{}, 0)
		lastHandleTime := time.Now()
		ticker := time.NewTicker(m.Timeout)

		handleItems := func() {
			lastHandleTime = time.Now()

			err := m.Handler(items)
			if nil != err {
				logrus.Errorf("handle batch items failed. error: %s.", err)
				return
			}

			items = make([]interface{}, 0)
		}

		for {
			func() {
				defer func() {
					if iRec := recover(); iRec != nil {
						logrus.Errorf("batch shape panic: %s.", iRec)
					}
				}()

				select {
				case item := <-m.C:
					items = append(items, item)
					if len(items) >= m.Size {
						fmt.Println("normal handle")
						handleItems()
					}

				case <-ticker.C: // tick
					now := time.Now()
					if now.Sub(lastHandleTime) >= m.Timeout && len(items) > 0 { // 超时
						fmt.Println("timeout handle")
						handleItems()
					}
				}

			}()

		}
	}()
}

func (m *BatchShape) PushItem(items ...interface{}) {
	fmt.Println("push")
	for _, item := range items {
		m.C <- item
	}
}
