package store

import (
	"sync"
)

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
