/**
limit execution times
调用频次限制
*/
package current_limiting

import (
	"github.com/haozzzzzzzz/go-rapid-development/v2/limiting/store"
	"time"

	"github.com/sirupsen/logrus"
)

// 任务
type Handler func(datas ...interface{}) (err error)

// 每分钟时间内执行times次，每次返回batch_size条数据
type MinuteFrequencyLimiting struct {
	WorkerNumber uint          // worker数目
	Times        uint          // 频次。要大于0
	WaitInterval time.Duration // 如果没有数据，等待间隔
	MaxBatchSize uint          // 一次多少条数据
	Store        store.Store   // 数据存储
	Handler      Handler       // 任务处理函数
}

func (m *MinuteFrequencyLimiting) Start() (err error) {

	// frequency Control
	datasC := m.Control()

	// workers
	for i := uint(0); i < m.WorkerNumber; i++ {
		go func(workerId uint) {
			logrus.Printf("running minute frequency limiting worker_%d", workerId)
			var err error
			for { // 循环做任务
				datas, ok := <-datasC
				if !ok { // 关闭
					return
				}

				err = m.Handler(datas...)
				if nil != err {
					logrus.Errorf("handle data failed. %s", err)
				}
			}
		}(i)
	}

	return
}

func (m *MinuteFrequencyLimiting) Control() (datasC chan []interface{}) {
	datasC = make(chan []interface{}, 100)
	go func() {
		lastMinuteTime := time.Now()
		var lastMinuteTimes uint
		for {
			datas := m.Store.Pop(m.MaxBatchSize)
			if len(datas) == 0 && m.WaitInterval > 0 {
				time.Sleep(m.WaitInterval)
			} else {
				datasC <- datas

				now := time.Now()
				if now.Minute() != lastMinuteTime.Minute() || now.Sub(lastMinuteTime).Minutes() >= 1 { // 重新标记当前分钟的计数
					lastMinuteTimes = 0
					lastMinuteTime = now
				}

				lastMinuteTimes++

				if lastMinuteTimes >= m.Times { // 达到当前分钟数的限制
					// 等待到下一分钟
					time.Sleep(time.Duration((60 - now.Second())) * time.Second)
				}
			}
		}
	}()

	return
}

func (m *MinuteFrequencyLimiting) AcceptData(datas ...interface{}) (err error) {
	m.Store.Push(datas...)
	return
}
