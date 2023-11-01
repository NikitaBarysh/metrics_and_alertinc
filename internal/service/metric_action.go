package service

import "github.com/NikitaBarysh/metrics_and_alertinc/internal/repository/memstorage"

type MetricAction struct {
	storage *memstorage.MemStorage
	sender  sender
}

func NewMetricAction(storage *memstorage.MemStorage, sender sender) *MetricAction {
	return &MetricAction{
		storage: storage,
		sender:  sender,
	}
}
