// Package entity - здесь структура на основе которой мы формируем метрики
package entity

import (
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/interface/models"
)

type MType string

const (
	Gauge   MType = "gauge"
	Counter MType = "counter"
)

// Metric - сущность метрики
type Metric struct {
	ID    string  // имя метрики
	MType MType   // параметр, принимающий значение gauge или counter
	Delta int64   // значение метрики в случае передачи counter
	Value float64 // значение метрики в случае передачи gauge
}

func (m *Metric) SetMetric(id string, mType MType, delta int64, value float64) (*Metric, error) {
	switch mType {
	case Gauge:
		return &Metric{
			ID:    id,
			MType: mType,
			Value: value,
		}, nil
	case Counter:
		return &Metric{
			ID:    id,
			MType: mType,
			Delta: delta,
		}, nil
	}
	return nil, models.ErrUnknownType
}
