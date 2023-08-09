package storage

import (
	"sync"
)

func CreateMemStorage() *MemStorage {
	return &MemStorage{
		gauge:   make(map[string]float64),
		counter: make(map[string]int64),
	}
}

type MemStorage struct {
	gauge   map[string]float64
	counter map[string]int64
	mu      sync.RWMutex
}

func (m *MemStorage) UpdateGaugeMetric(key string, value float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.gauge[key] = value
	return
}

func (m *MemStorage) UpdateCounterMetric(key string, value int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.counter[key] += value
	return
}

//func (m *MemStorage) GetGaugeMetric() []string {
//	m.mu.RLock()
//	defer m.mu.RUnlock()
//	metricNameSlice := make([]string, 0, 30)
//	for metricName, _ := range m.gauge {
//		metricNameSlice = append(metricNameSlice, metricName)
//	}
//	return metricNameSlice
//}
//
//func (m *MemStorage) GetCounterMetric() []string {
//	m.mu.RLock()
//	defer m.mu.RUnlock()
//	metricNameSlice := make([]string, 0, 30)
//	for metricName, _ := range m.counter {
//		metricNameSlice = append(metricNameSlice, metricName)
//	}
//	return metricNameSlice
//}

func (m *MemStorage) ReadGaugeMetric() map[string]float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	// TODO copymap
	return m.gauge
}

func (m *MemStorage) ReadCounterMetric() map[string]int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.counter
}
