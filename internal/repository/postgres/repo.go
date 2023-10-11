package postgres

import (
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/entity"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/interface/models"
	"sync"
)

type DBStorage struct {
	MetricMap map[string]entity.Metric

	mu sync.RWMutex
}

func NewDBStorage() *DBStorage {
	return &DBStorage{
		MetricMap: make(map[string]entity.Metric),
	}
}

//func (m *DBStorage) SetOnUpdate(fn func()) {
//	m.onUpdate = fn
//}

func (m *DBStorage) UpdateGaugeMetric(key string, value float64) {
	m.mu.Lock()
	m.MetricMap[key] = entity.Metric{ID: key, MType: "gauge", Value: value}
	//fn := m.onUpdate
	m.mu.Unlock()
	//if fn != nil {
	//	fn()
	//}
}

func (m *DBStorage) UpdateCounterMetric(key string, value int64) {
	m.mu.Lock()
	metricValue := m.MetricMap[key].Delta
	metricValue += value
	m.MetricMap[key] = entity.Metric{ID: key, MType: "counter", Delta: metricValue}
	//fn := m.onUpdate
	m.mu.Unlock()
	//if fn != nil {
	//	fn()
	//}
}

func (m *DBStorage) GetMetric(key string) (entity.Metric, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	metricStruct, ok := m.MetricMap[key]
	if ok {
		return metricStruct, nil
	}
	return metricStruct, models.ErrNotFound
}

func (m *DBStorage) ReadMetric() map[string]entity.Metric {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.MetricMap
}

func (m *DBStorage) GetAllMetric() []entity.Metric {
	m.mu.RLock()
	defer m.mu.RUnlock()
	metricSlice := make([]entity.Metric, 0, len(m.MetricMap))
	for _, metricValue := range m.MetricMap {
		metricSlice = append(metricSlice, metricValue)
	}
	return metricSlice
}

func (m *DBStorage) SetMetric(data map[string]entity.Metric) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.MetricMap = data
}
