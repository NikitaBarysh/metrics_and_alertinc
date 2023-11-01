package service

import "github.com/NikitaBarysh/metrics_and_alertinc/internal/repository/mem_storage"

type MetricAction struct {
	storage *mem_storage.MemStorage
	sender  sender
}

func NewMetricAction(storage *mem_storage.MemStorage, sender sender) *MetricAction {
	return &MetricAction{
		storage: storage,
		sender:  sender,
	}
}
