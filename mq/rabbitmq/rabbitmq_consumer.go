package rabbitmq

import (
	"time"

	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

// TODO 加监控
type QueueConsumer struct {
	url          string
	queue        string
	channelCount uint32
	workerCount  uint32
	msgC         chan *amqp.Delivery
	handler      func(msg []byte) (err error)
}

func NewQueueConsumer(
	url string,
	queue string,
	channelCount uint32,
	workerCount uint32,
	handler func(msg []byte) (err error),
) *QueueConsumer {
	return &QueueConsumer{
		url:          url,
		queue:        queue,
		channelCount: channelCount,
		workerCount:  workerCount,
		msgC:         make(chan *amqp.Delivery, workerCount),
		handler:      handler,
	}
}

func (m *QueueConsumer) Start() {
	// run wokers
	for i := uint32(0); i < m.workerCount; i++ {
		go func(workerOrder uint32) {
			for msg := range m.msgC {
				func(msg *amqp.Delivery) {
					if msg == nil {
						return
					}

					var errGo error
					defer func() {
						if iRec := recover(); iRec != nil {
							logrus.Errorf("handler panic: %s.", iRec)
						}
					}()

					defer func() {
						if errGo == nil {
							msg.Ack(false)
						} else { // 如果失败则丢弃
							msg.Nack(false, false)
						}
					}()

					errGo = m.handler(msg.Body)
					if errGo != nil {
						logrus.Errorf("handle msg failed. error: %s.", errGo)
					}
				}(msg)
			}
		}(i)
	}

	// connect
	go func() {
		for {
			err := m.connect()
			if nil != err {
				logrus.Errorf("connection is lost, trying to reconnecting. error: %s.", err)
			}

			time.Sleep(1 * time.Second)
		}
	}()
}

func (m *QueueConsumer) connect() (err error) {
	defer func() {
		if iRc := recover(); iRc != nil {
			logrus.Errorf("connection panic: %s.", iRc)
		}
	}()

	// init connection
	conn, err := amqp.Dial(m.url)
	if nil != err {
		logrus.Errorf("amqp dial rabbitmq failed. url: %s. error: %s.", m.url, err)
		return
	}
	defer conn.Close()

	logrus.Infof("connect ampq successfully. url: %s", m.url)

	connClose := make(chan *amqp.Error, 10)
	conn.NotifyClose(connClose)

	// channels
	for i := uint32(0); i < m.channelCount; i++ {
		go func(channelOrder uint32) {
			var errGo error
			defer func() {
				if iRec := recover(); iRec != nil {
					logrus.Errorf("channel panic: %s.", iRec)
				}
			}()

			logrus.Debugf("running channel : %d", channelOrder)

			channel, errGo := conn.Channel()
			if errGo != nil {
				logrus.Errorf("create connection channel failed. error: %s.", err)
			}
			defer channel.Close()

			channelClose := make(chan *amqp.Error) // 会被channel.NotifyClose关闭
			channel.NotifyClose(channelClose)

			errGo = channel.Qos(1, 0, false)
			if nil != errGo {
				logrus.Errorf("set channel qos failed. error: %s.", err)
				return
			}

			msgs, errGo := channel.Consume(
				m.queue,
				fmt.Sprintf("channel_%d", channelOrder),
				false,
				false,
				false,
				false,
				nil,
			)
			if nil != errGo {
				logrus.Errorf("channel consume queue failed. queue_name: %s. error: %s.", m.queue, err)
				return
			}

			for {
				select {
				case msg := <-msgs:
					m.msgC <- &msg

				case <-channelClose:
					logrus.Warnf("consume channel close. order: %d", channelOrder)
					return
				}
			}
		}(i)
	}

	err = <-connClose
	logrus.Errorf("connection closed. error: %s.", err)

	return
}
