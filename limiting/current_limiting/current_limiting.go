package current_limiting

import (
	"time"

	"sync"

	"fmt"

	"github.com/sirupsen/logrus"
)

type Limiting interface {
	AcceptData(datas ...interface{}) (err error)
}

type Store interface {
	Push(datas ...interface{})
	Pop(number uint) (datas []interface{})
}

// 数据存储
type MemoryStore struct {
	mutex sync.Mutex
	items []interface{}
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		items: make([]interface{}, 0),
	}
}

func (m *MemoryStore) Push(datas ...interface{}) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.items = append(m.items, datas...)

	return
}

func (m *MemoryStore) Pop(number uint) (datas []interface{}) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	datas = make([]interface{}, 0)
	lenItems := uint(len(m.items))
	if lenItems == 0 {
		return
	}

	if lenItems <= number {
		datas = make([]interface{}, lenItems)
		copy(datas, m.items)
		m.items = make([]interface{}, 0)
		return
	}

	datas = make([]interface{}, number)
	copy(datas, m.items[:number])

	m.items = m.items[number:]

	return
}

func (m *MemoryStore) Get() (items []interface{}) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	items = make([]interface{}, len(m.items))
	copy(items, m.items)
	return
}

// 任务
type Handler func(datas ...interface{}) (err error)

// 每分钟时间内执行times次，每次返回batch_size条数据
type MinuteFrequencyLimiting struct {
	WorkerNumber uint          // worker数目
	Times        uint          // 频次。要大于0
	WaitInterval time.Duration // 如果没有数据，等待间隔
	MaxBatchSize uint          // 一次多少条数据
	Store        Store         // 数据存储
	Handler      Handler       // 任务处理函数
}

func (m *MinuteFrequencyLimiting) Start() (err error) {

	// frequency Control
	datasC := m.Control()

	// workers
	for i := uint(0); i < m.WorkerNumber; i++ {
		go func(workerId uint) {
			fmt.Println("running worker")
			var err error
			for { // 循环做任务
				datas, ok := <-datasC
				if !ok { // 关闭
					return
				}

				fmt.Printf("worker: %d, data: %#v\n", workerId, datas)
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
