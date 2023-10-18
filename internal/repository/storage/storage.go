package storage

import (
	"fmt"
	"github.com/NikitaBarysh/metrics_and_alertinc/config/server"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/entity"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/interface/models"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/repository/file_storage"
	"sync"
)

type MemStorage struct {
	MetricMap  map[string]entity.Metric
	onUpdate   func()
	mu         sync.RWMutex
	FileEngine *file_storage.FileEngine
}

func NewAgentStorage() *MemStorage {
	return &MemStorage{
		MetricMap: make(map[string]entity.Metric),
	}
}

//fileEngine *service.FileEngine

func NewMemStorage(cfg *server.Config) (*MemStorage, error) {
	fileEngine, err := file_storage.NewFileEngine(cfg.StorePath)
	if err != nil {
		return nil, fmt.Errorf("storage: newMemStorage: NewFileEngine: %w", err)
	}
	return &MemStorage{
		MetricMap:  make(map[string]entity.Metric),
		FileEngine: fileEngine,
	}, nil
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

func (m *MemStorage) GetAllMetric() ([]entity.Metric, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	metricSlice := make([]entity.Metric, 0, len(m.MetricMap))
	for _, metricValue := range m.MetricMap {
		metricSlice = append(metricSlice, metricValue)
	}
	return metricSlice, nil
}

func (m *MemStorage) SetMetrics(metric []entity.Metric) error {
	for _, v := range metric {
		switch v.MType {
		case "counter":
			m.mu.Lock()
			metricValue := m.MetricMap[v.ID].Delta
			metricValue += v.Delta
			m.MetricMap[v.ID] = entity.Metric{ID: v.ID, MType: v.MType, Delta: metricValue}
			m.mu.Unlock()
		case "gauge":
			m.mu.Lock()
			m.MetricMap[v.ID] = v
			m.mu.Unlock()
		default:
			return models.ErrNotFound
		}
	}
	return nil
}
