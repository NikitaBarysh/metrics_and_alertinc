// Package entity - здесь структура на основе которой мы формируем метрики
package entity

import (
	"strconv"

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

func NewMetric(id, value string, mType MType) (*Metric, error) {
	switch mType {
	case Gauge:
		val, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, models.ErrWrongValue
		}
		return &Metric{
			ID:    id,
			MType: mType,
			Value: val,
		}, nil
	case Counter:
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
