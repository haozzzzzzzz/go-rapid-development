package rabbitmq

import (
	"fmt"
	"testing"
	"time"
)

func TestExchangePublisher_Publish(t *testing.T) {
	publisher := NewExchangePublisher(
		"amqp://admin:admin@localhost:5672",
		"video_buddy_pay_out_order",
		2,
		100,
	)

	publisher.Start()
	for {
		publisher.Publish("test", map[string]interface{}{
			"hello": "world",
		})
		time.Sleep(200 * time.Millisecond)
	}
	fmt.Println("world")

}
