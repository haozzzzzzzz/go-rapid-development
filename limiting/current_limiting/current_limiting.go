package current_limiting

import (
	"time"

	"sync"

	"github.com/sirupsen/logrus"
)

type Limiting interface {
	AcceptData(data interface{}) (err error)
}

type Store interface {
	Push(data interface{}) (err error)
	Pop() (data interface{}, err error)
}

// 数据存储
// TODO 使用数据结构
// https://github.com/Workiva/go-datastructures
type MemoryStore struct {
	Mutex sync.Mutex
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{}
}

func (m *MemoryStore) Push(data interface{}) (err error) {
	return
}

func (m *MemoryStore) Pop() (data interface{}, err error) {
	return
}

// 任务
type Handler func(data interface{}) (err error)

// TODO 需要实现
type FrequencyLimiting struct {
	WorkerNumber uint32        // worker数目
	Duration     time.Duration // 时长
	Times        uint32        // 频次
	Store        Store         // 数据存储

	Handler Handler // 任务处理函数
}

func (m *FrequencyLimiting) Start() (err error) {
	for i := uint32(0); i < m.WorkerNumber; i++ {
		go func(workerId uint32) {
			for { // 循环做任务
				data, err := m.Store.Pop()
				if nil != err {
					logrus.Errorf("pop store failed. error: %s.", err)
					continue
				}

				err = m.Handler(data)
				if nil != err {
					logrus.Errorf("handle data failed. error: %s.", err)
					continue
				}
			}
		}(i)
	}

	return
}

func (m *FrequencyLimiting) AcceptData(data interface{}) (err error) {
	m.Store.Push(data)
	return
}
