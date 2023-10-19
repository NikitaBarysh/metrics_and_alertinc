package storage

import (
	"context"
	"github.com/NikitaBarysh/metrics_and_alertinc/config/server"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/entity"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/interface/models"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/repository/file_storage"
	"sync"
	"time"
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

func NewMemStorage(cfg *server.Config, file *file_storage.FileEngine) (*MemStorage, error) { // TODO ctx
	m := &MemStorage{}
	data := make(map[string]entity.Metric)
	if file != nil {
		data, _ = file.GetAllMetric()
		go m.syncData(cfg.StoreInterval)
	}
	m.mu = sync.RWMutex{}
	m.MetricMap = data
	m.FileEngine = file
	return m, nil
}

func (m *MemStorage) syncData(interval uint64) {
	timeTicker := time.NewTicker(time.Second * time.Duration(interval))
	defer timeTicker.Stop()
	for {
		select {
		case <-timeTicker.C:
			m.FileEngine.SetMetrics(m.MetricMap)
			//case <-ctx.Done():
			//	return
		}
	}
}

func (m *MemStorage) UpdateGaugeMetric(key string, value float64) {
	m.mu.Lock()
	m.MetricMap[key] = entity.Metric{ID: key, MType: "gauge", Value: value}
	m.mu.Unlock()

}

func (m *MemStorage) UpdateCounterMetric(key string, value int64) {
	m.mu.Lock()
	metricValue := m.MetricMap[key].Delta
	metricValue += value
	m.MetricMap[key] = entity.Metric{ID: key, MType: "counter", Delta: metricValue}
	m.mu.Unlock()
}

func (m *MemStorage) GetMetric(key string) (entity.Metric, error) {
	//fmt.Println(m)
	//fmt.Println(m.mu)
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

func (m *MemStorage) CheckPing(ctx context.Context) error {
	return nil
}
