package repositories

import (
	"fmt"
	"sync"

	"github.com/NikitaBarysh/metrics_and_alertinc/internal/models"
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

type keeper interface {
	Restore() (map[string]MemStorageStruct, error)
	Flush(data map[string]MemStorageStruct) error
}

type MemStorage struct {
	MemStorageMap map[string]MemStorageStruct
	Keeper        keeper
	mu            sync.RWMutex
}

func NewMemStorage(keeper keeper) *MemStorage {
	var data map[string]MemStorageStruct
	if keeper == nil {
		data = make(map[string]MemStorageStruct)
	} else {
		data, _ = keeper.Restore()
	}
	storage := &MemStorage{
		MemStorageMap: data,
		Keeper:        keeper,
	}
	return storage
}

func (m *MemStorage) SaveData() {
	m.mu.RLock()
	defer m.mu.RUnlock()
	m.Keeper.Flush(m.MemStorageMap)
}

func (m *MemStorage) UpdateGaugeMetric(key string, value float64) {
	m.mu.Lock()
	m.MemStorageMap[key] = MemStorageStruct{ID: key, MType: "gauge", Value: value}
	m.mu.Unlock()
}

func (m *MemStorage) UpdateCounterMetric(key string, value int64) {
	m.mu.Lock()
	metricValue := m.MemStorageMap[key].Delta
	metricValue += value
	m.MemStorageMap[key] = MemStorageStruct{ID: key, MType: "counter", Delta: metricValue}
	m.mu.Unlock()

}

func (m *MemStorage) ReadDefinitelyMetric(key string) (MemStorageStruct, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	metricStruct, ok := m.MemStorageMap[key]
	if ok {
		return metricStruct, nil
	}
	return metricStruct, models.ErrNotFound
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
