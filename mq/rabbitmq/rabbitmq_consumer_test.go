package rabbitmq

import (
	"fmt"
	"testing"
	"time"
)

func TestQueueConsumer_Start(t *testing.T) {
	NewQueueConsumer(
		"amqp://admin:admin@localhost:5672",
		"vcoin_out_order",
		2,
		2,
		func(msg []byte) (err error) {
			fmt.Println(string(msg))
			return
		},
	).Start()
	time.Sleep(10 * time.Minute)
}
