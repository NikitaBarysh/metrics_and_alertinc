package service

import (
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/repository/storage"
)

type MetricAction struct {
	storage *storage.MemStorage
	sender  sender
}

func NewMetricAction(storage *storage.MemStorage, sender sender) *MetricAction {
	return &MetricAction{
		storage: storage,
		sender:  sender,
	}
}
