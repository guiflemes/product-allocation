package domain

import (
	"sync"
)

type Events struct {
	items []interface{}
	mu    sync.RWMutex
}

func (m *Events) Append(x interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.items = append(m.items, x)
}

func (m *Events) Read() []interface{} {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.items
}

func (m *Events) Len() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.items)
}

type Allocated struct {
	OrderId  string
	Sku      string
	Qty      int
	BatchRef string
}

type OutOfStock struct {
	Sku string
}

type Deallocate struct {
	OrderId string
	Sku     string
	qty     int
}
