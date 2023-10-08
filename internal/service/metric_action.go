package service

import (
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/repository/storage"
)

type MetricAction struct {
	MemStorage *storage.MemStorage
	sender     sender
}

func NewMetricAction(memStorage *storage.MemStorage, sender sender) *MetricAction {
	return &MetricAction{
		MemStorage: memStorage,
		sender:     sender,
	}
}
