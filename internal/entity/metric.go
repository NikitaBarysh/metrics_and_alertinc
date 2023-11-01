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


