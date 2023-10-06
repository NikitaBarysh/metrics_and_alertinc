package service

import (
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/repository/memory"
)

type MetricAction struct {
	MemStorage *memory.MemStorage
	sender     sender
}

func NewMetricAction(memStorage *memory.MemStorage, sender sender) *MetricAction {
	return &MetricAction{
		MemStorage: memStorage,
		sender:     sender,
	}
}
