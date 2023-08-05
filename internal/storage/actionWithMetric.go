package storage

import (
	"fmt"
	"sync"
)

func CreateMemStorage() *MemStorage {
	return &MemStorage{
		gauge: make(map[string]float64),
	}
}

type MemStorage struct {
	gauge map[string]float64
	mu    sync.RWMutex
}

func (m *MemStorage) Put(key string, value float64) {
	m.mu.Lock()
	m.mu.Unlock()
	m.gauge[key] = value
	return
}

func (m *MemStorage) Get() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	metricNameSlice := make([]string, 0, 30)
	for metricName, _ := range m.gauge {
		metricNameSlice = append(metricNameSlice, metricName)
	}
	fmt.Println(metricNameSlice)
	return metricNameSlice
}

func (m *MemStorage) Read(key string) (interface{}, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	value, ok := m.gauge[key]
	return value, ok
}
