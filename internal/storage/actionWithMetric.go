package storage

import (
	"sync"
)

func CreateMemStorage() *MemStorage {
	return &MemStorage{
		storage: make(map[string]interface{}),
	}
}

type MemStorage struct {
	storage map[string]interface{}
	mu      sync.RWMutex
}

func (m *MemStorage) Put(key string, value interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.storage[key] = value
	return
}

func (m *MemStorage) Get() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	metricNameSlice := make([]string, 0, 30)
	for metricName, _ := range m.storage {
		metricNameSlice = append(metricNameSlice, metricName)
	}
	return metricNameSlice
}

func (m *MemStorage) Read(key string) (interface{}, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	value, ok := m.storage[key]
	return value, ok
}
