// Package models  здесь сущность метрики для JSON и кастомные ошибки
package models

import "errors"

const (
	Gauge   = "gauge"
	Counter = "counter"
)

// Кастомные ошибки
var (
	ErrNotFound    = errors.New("not found metric")
	ErrWrongValue  = errors.New("entity: ParseValue")
	ErrUnknownType = errors.New("entity: New: UnknownType")
)

// Metrics - cущность для JSON
type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func NewMetric(id, mType string, delta *int64, value *float64) Metrics {
	return Metrics{
		ID:    id,
		MType: mType,
		Delta: delta,
		Value: value,
	}
}

// NewMetricValue - добавляет значение gauge метрики
func (m *Metrics) NewMetricValue(value float64) {
	m.Value = &value
}

// NewMetricDelta - добавляет значение counter метрики
func (m *Metrics) NewMetricDelta(delta int64) {
	m.Delta = &delta
}
