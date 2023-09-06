package storage

import "github.com/NikitaBarysh/metrics_and_alertinc/internal/storage/repositories"

type MetricAction struct {
	MemStorage *repositories.MemStorage
	sender     sender
}

func NewMetricAction(memStorage *repositories.MemStorage, sender sender) *MetricAction {
	return &MetricAction{
		MemStorage: memStorage,
		sender:     sender,
	}
}
