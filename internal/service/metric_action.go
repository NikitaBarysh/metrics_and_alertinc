package service

import (
	"github.com/NikitaBarysh/metrics_and_alertinc/config/agent"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/repository/memstorage"
)

type MetricAction struct {
	storage *memstorage.MemStorage
	sender  sender
	cfg     agent.Config
}

func NewMetricAction(storage *memstorage.MemStorage, sender sender, cfg *agent.Config) *MetricAction {
	return &MetricAction{
		storage: storage,
		sender:  sender,
		cfg:     *cfg,
	}
}
