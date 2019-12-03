package batch_shape

import (
	"fmt"
	"testing"
	"time"
)

func TestBatchShape(t *testing.T) {
	batchShape := NewBatchShape(5, 10*time.Second, func(items []interface{}) (err error) {
		time.Sleep(3 * time.Second)
		fmt.Printf("%#v\n", items)
		return
	})
	batchShape.Start()

	i := 0
	for {
		i++
		batchShape.PushItem(i)
		if i%6 == 0 {
			time.Sleep(20 * time.Second)
		} else {
			time.Sleep(1 * time.Second)
		}
	}
}
