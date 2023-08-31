package models

import "errors"

var ErrNotFound = errors.New("not found metric struct")

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

func (m *Metrics) NewMetricValue(value float64) {
	m.Value = &value
}

func (m *Metrics) NewMetricDelta(delta int64) {
	m.Delta = &delta
}
