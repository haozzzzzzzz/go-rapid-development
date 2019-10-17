package rabbitmq

import (
	"time"

	"encoding/json"

	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

// TODO 加监控
type PublishMsg struct {
	Exchange string      `json:"exchange"`
	Key      string      `json:"key"`
	Body     interface{} `json:"body"`
}

type ExchangePublisher struct {
	url        string
	exchange   string
	msgC       chan *PublishMsg
	channelNum uint32
}

func NewExchangePublisher(
	url string,
	exchange string,
	channelNum uint32, // 发送worker个数
	msgBufferSize uint32,
) *ExchangePublisher {
	return &ExchangePublisher{
		url:        url,
		exchange:   exchange,
		channelNum: channelNum,
		msgC:       make(chan *PublishMsg, msgBufferSize),
	}
}

func (m *ExchangePublisher) Start() {
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

func (m *ExchangePublisher) connect() (err error) {
	defer func() {
		if iRc := recover(); iRc != nil {
			logrus.Errorf("exchange publisher connection panic: %s.", iRc)
		}
	}()

	conn, err := amqp.Dial(m.url)
	if nil != err {
		logrus.Errorf("amqp dial rabbitmq failed. url: %s, exchange: %s, error: %s.", m.url, m.exchange, err)
		return
	}
	defer func() {
		errClose := conn.Close()
		if errClose != nil {
			logrus.Errorf("close connection failed. error: %s.", err)
		}
	}()

	logrus.Infof("connect ampq successfully. url: %s, exchange: %s", m.url, m.exchange)
	connClose := make(chan *amqp.Error, 10) // close by NotifyClose
	conn.NotifyClose(connClose)

	// publish worker
	for i := uint32(0); i < m.channelNum; i++ {
		go func(channelOrder uint32) {
			defer func() {
				if iRec := recover(); iRec != nil {
					logrus.Errorf("panic: %s.", iRec)
				}
			}()

			var errGo error
			channel, errGo := conn.Channel()
			if nil != errGo {
				logrus.Errorf("create home page exchange channel failed. error: %s.", err)
				return
			}
			defer func() {
				errClose := channel.Close()
				if errClose != nil {
					logrus.Errorf("close channel failed. error: %s.", err)
				}
			}()

			channelClose := make(chan *amqp.Error, 10) // 会被channel.NotifyClose关闭
			channel.NotifyClose(channelClose)

			for {
				select {
				case msg := <-m.msgC:
					if msg == nil {
						return
					}

					byteMessage, err := json.Marshal(msg)
					if nil != err {
						logrus.Errorf("marshal message failed. error: %s.", err)
						break
					}

					err = channel.Publish(
						msg.Exchange,
						msg.Key,
						false,
						false,
						amqp.Publishing{
							ContentType: "application.json",
							Body:        byteMessage,
						},
					)
					if nil != err {
						logrus.Errorf("publish message failed. error: %s.", err)
						break
					}

				case <-channelClose:
					logrus.Warnf("publish channel close. order: %d", channelOrder)
					return
				}
			}

		}(i)
	}

	err = <-connClose
	logrus.Warnf("connection close. %s", err)
	return
}

func (m *ExchangePublisher) Publish(key string, obj interface{}) {
	m.msgC <- &PublishMsg{
		Exchange: m.exchange,
		Key:      key,
		Body:     obj,
	}
}
