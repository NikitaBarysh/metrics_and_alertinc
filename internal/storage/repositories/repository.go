package repositories

import (
	"fmt"
	"sync"
)

func NewMemStorageStruct(id, mType string, delta int64, value float64) MemStorageStruct {
	return MemStorageStruct{
		ID:    id,
		MType: mType,
		Delta: delta,
		Value: value,
	}
}

type MemStorageStruct struct {
	ID    string  `json:"id"`              // имя метрики
	MType string  `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		MemStorageMap: make(map[string]MemStorageStruct),
	}
}

type MemStorage struct {
	MemStorageMap map[string]MemStorageStruct
	mu            sync.RWMutex
}

func (m *MemStorage) UpdateGaugeMetric(key string, value float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	metricType := m.MemStorageMap[key].MType
	if metricType == "gauge" {
		m.MemStorageMap[key] = MemStorageStruct{ID: key, MType: "gauge", Value: value}
	}
}

func (m *MemStorage) UpdateCounterMetric(key string, value int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	metricType := m.MemStorageMap[key].MType
	metricValue := m.MemStorageMap[key].Delta
	if metricType == "counter" {
		metricValue = value + 1
		m.MemStorageMap[key] = MemStorageStruct{ID: key, MType: "counter", Delta: metricValue}
	}
}
func (m *MemStorage) ReadMetric() map[string]MemStorageStruct {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.MemStorageMap
}

func (m *MemStorage) GetAllMetric() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	metricSlice := make([]string, 0, len(m.MemStorageMap))
	for metricName, metricValue := range m.MemStorageMap {
		metricType := m.MemStorageMap[metricName].MType
		switch metricType {
		case "gauge":
			metricSlice = append(metricSlice, fmt.Sprintf("%s = %2f", metricName, metricValue.Delta))
		case "counter":
			metricSlice = append(metricSlice, fmt.Sprintf("%s = %d", metricName, metricValue.Value))
		}
	}
	return metricSlice
}
