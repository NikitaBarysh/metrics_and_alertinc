package storage

import (
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/service"
)

type MetricAction struct {
	MemStorage *service.MemStorage
	sender     sender
}

func NewMetricAction(memStorage *service.MemStorage, sender sender) *MetricAction {
	return &MetricAction{
		MemStorage: memStorage,
		sender:     sender,
	}
}
