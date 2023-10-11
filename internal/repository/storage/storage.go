package storage

import (
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/entity"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/interface/models"
	"sync"
)

type MemStorage struct {
	MetricMap map[string]entity.Metric
	onUpdate  func()
	mu        sync.RWMutex
	//flush  *flusher.Flusher
}

//fileEngine *service.FileEngine

func NewMemStorage() *MemStorage {
	return &MemStorage{
		MetricMap: make(map[string]entity.Metric),
		//flush: fileEngine,
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

func (m *MemStorage) GetMetric(key string) (entity.Metric, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	metricStruct, ok := m.MetricMap[key]
	if ok {
		return metricStruct, nil
	}
	return metricStruct, models.ErrNotFound
}

func (m *MemStorage) GetAllMetric() []entity.Metric {
	m.mu.RLock()
	defer m.mu.RUnlock()
	metricSlice := make([]entity.Metric, 0, len(m.MetricMap))
	for _, metricValue := range m.MetricMap {
		metricSlice = append(metricSlice, metricValue)
	}
	return metricSlice
}

func (m *MemStorage) SetMetric(data []entity.Metric) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, value := range data {
		m.MetricMap[value.ID] = value
	}
}