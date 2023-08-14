package repositories

import (
	"fmt"
	"sync"
)

func NewMemStorage() *MemStorage {
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
}

func (m *MemStorage) UpdateCounterMetric(key string, value int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.counter[key] += value
}

func (m *MemStorage) ReadGaugeMetric() map[string]float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.gauge
}

func (m *MemStorage) ReadCounterMetric() map[string]int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.counter
}

func (m *MemStorage) GetAllMetric() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	metricSlice := make([]string, 0, len(m.counter)+len(m.gauge))
	for metricName, metricValue := range m.gauge {
		metricSlice = append(metricSlice, fmt.Sprintf("%s = %2f", metricName, metricValue))
	}
	for metricName, metricValue := range m.counter {
		metricSlice = append(metricSlice, fmt.Sprintf("%s = %d", metricName, metricValue))
	}
	return metricSlice
}