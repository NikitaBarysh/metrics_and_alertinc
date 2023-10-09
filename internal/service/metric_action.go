package service

import (
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/repository/storage"
)

type MetricAction struct {
	//storage repository.Storage
	MemStorage *storage.MemStorage
	sender     sender
}

func NewMetricAction(MemStorage *storage.MemStorage, sender sender) *MetricAction {
	return &MetricAction{
		MemStorage: MemStorage,
		sender:     sender,
	}
}
