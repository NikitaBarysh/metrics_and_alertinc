package repositories

import (
	"fmt"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/models"
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
	onUpdate      func()
	mu            sync.RWMutex
}

func (m *MemStorage) SetOnUpdate(fn func()) {
	m.onUpdate = fn
}

func (m *MemStorage) UpdateGaugeMetric(key string, value float64) {
	m.mu.Lock()
	m.MemStorageMap[key] = MemStorageStruct{ID: key, MType: "gauge", Value: value}
	fn := m.onUpdate
	m.mu.Unlock()
	if fn != nil {
		fn()
	}
}

func (m *MemStorage) UpdateCounterMetric(key string, value int64) {
	m.mu.Lock()
	metricValue := m.MemStorageMap[key].Delta
	metricValue += value
	m.MemStorageMap[key] = MemStorageStruct{ID: key, MType: "counter", Delta: metricValue}
	fn := m.onUpdate
	m.mu.Unlock()
	if fn != nil {
		fn()
	}
}

func (m *MemStorage) ReadDefinitelyMetric(key string) (MemStorageStruct, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	metricStruct, ok := m.MemStorageMap[key]
	if ok {
		return metricStruct, nil
	}
	return MemStorageStruct{}, models.ErrNotFound
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
			metricSlice = append(metricSlice, fmt.Sprintf("%s = %d", metricName, metricValue.Delta))
		case "counter":
			metricSlice = append(metricSlice, fmt.Sprintf("%s = %2f", metricName, metricValue.Value))
		}
	}
	return metricSlice
}

func (m *MemStorage) PutMetricMap(data map[string]MemStorageStruct) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.MemStorageMap = data
}
