package entity

import (
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/interface/models"
	"strconv"
)

type Metric struct {
	ID    string  // имя метрики
	MType string  // параметр, принимающий значение gauge или counter
	Delta int64   // значение метрики в случае передачи counter
	Value float64 // значение метрики в случае передачи gauge
}

func NewMetric(id, mType, value string) (*Metric, error) {
	switch mType {
	case "gauge":
		val, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, models.ErrWrongValue
		}
		return &Metric{
			ID:    id,
			MType: mType,
			Value: val,
		}, nil
	case "counter":
		val, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil, models.ErrWrongValue
		}
		return &Metric{
			ID:    id,
			MType: mType,
			Delta: val,
		}, nil
	}
	return nil, models.ErrUnknownType
}

//func NewMemStorage() *MemStorage {
//	return &MemStorage{
//		MetricMap: make(map[string]Metric),
//	}
//}
//
//type MemStorage struct {
//	MetricMap map[string]Metric
//	onUpdate  func()
//	mu        sync.RWMutex
//}
//
//func (m *MemStorage) SetOnUpdate(fn func()) {
//	m.onUpdate = fn
//}
//
//func (m *MemStorage) UpdateGaugeMetric(key string, value float64) {
//	m.mu.Lock()
//	m.MetricMap[key] = Metric{ID: key, MType: "gauge", Value: value}
//	fn := m.onUpdate
//	m.mu.Unlock()
//	if fn != nil {
//		fn()
//	}
//}
//
//func (m *MemStorage) UpdateCounterMetric(key string, value int64) {
//	m.mu.Lock()
//	metricValue := m.MetricMap[key].Delta
//	metricValue += value
//	m.MetricMap[key] = Metric{ID: key, MType: "counter", Delta: metricValue}
//	fn := m.onUpdate
//	m.mu.Unlock()
//	if fn != nil {
//		fn()
//	}
//}
//
//func (m *MemStorage) GetMetric(key string) (Metric, error) {
//	m.mu.RLock()
//	defer m.mu.RUnlock()
//	metricStruct, ok := m.MetricMap[key]
//	if ok {
//		return metricStruct, nil
//	}
//	return metricStruct, models.ErrNotFound
//}
//
//func (m *MemStorage) ReadMetric() map[string]Metric {
//	m.mu.RLock()
//	defer m.mu.RUnlock()
//	return m.MetricMap
//}
//
//func (m *MemStorage) GetAllMetric() []string {
//	m.mu.RLock()
//	defer m.mu.RUnlock()
//	metricSlice := make([]string, 0, len(m.MetricMap))
//	for metricName, metricValue := range m.MetricMap {
//		metricType := m.MetricMap[metricName].MType
//		switch metricType {
//		case "gauge":
//			metricSlice = append(metricSlice, fmt.Sprintf("%s = %d", metricName, metricValue.Delta))
//		case "counter":
//			metricSlice = append(metricSlice, fmt.Sprintf("%s = %2f", metricName, metricValue.Value))
//		}
//	}
//	return metricSlice
//}
//
//func (m *MemStorage) SetMetric(data map[string]Metric) {
//	m.mu.Lock()
//	defer m.mu.Unlock()
//	m.MetricMap = data
//}
