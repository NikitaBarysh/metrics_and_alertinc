package storage

import (
	"fmt"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/entity"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/interface/models"
	"sync"
)

type MemStorage struct {
	MetricMap map[string]entity.Metric
	onUpdate  func()
	mu        sync.RWMutex
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		MetricMap: make(map[string]entity.Metric),
	}
}

func (m *MemStorage) SetOnUpdate(fn func()) {
	m.onUpdate = fn
}

func (m *MemStorage) UpdateGaugeMetric(key string, value float64) {
	m.mu.Lock()
	m.MetricMap[key] = entity.Metric{ID: key, MType: "gauge", Value: value}
	fn := m.onUpdate
	m.mu.Unlock()
	if fn != nil {
		fn()
	}
}

func (m *MemStorage) UpdateCounterMetric(key string, value int64) {
	m.mu.Lock()
	metricValue := m.MetricMap[key].Delta
	metricValue += value
	m.MetricMap[key] = entity.Metric{ID: key, MType: "counter", Delta: metricValue}
	fn := m.onUpdate
	m.mu.Unlock()
	if fn != nil {
		fn()
	}
}

func (m *MemStorage) ReadDefinitelyMetric(key string) (entity.Metric, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	metricStruct, ok := m.MetricMap[key]
	if ok {
		return metricStruct, nil
	}
	return metricStruct, models.ErrNotFound
}

func (m *MemStorage) ReadMetric() map[string]entity.Metric {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.MetricMap
}

func (m *MemStorage) GetAllMetric() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	metricSlice := make([]string, 0, len(m.MetricMap))
	for metricName, metricValue := range m.MetricMap {
		metricType := m.MetricMap[metricName].MType
		switch metricType {
		case "gauge":
			metricSlice = append(metricSlice, fmt.Sprintf("%s = %d", metricName, metricValue.Delta))
		case "counter":
			metricSlice = append(metricSlice, fmt.Sprintf("%s = %2f", metricName, metricValue.Value))
		}
	}
	return metricSlice
}

func (m *MemStorage) PutMetricMap(data map[string]entity.Metric) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.MetricMap = data
}
