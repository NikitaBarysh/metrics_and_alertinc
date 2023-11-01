package service

import "github.com/NikitaBarysh/metrics_and_alertinc/internal/repository/memStorage"

type MetricAction struct {
	storage *memStorage.MemStorage
	sender  sender
}

func NewMetricAction(storage *memStorage.MemStorage, sender sender) *MetricAction {
	return &MetricAction{
		storage: storage,
		sender:  sender,
	}
}
