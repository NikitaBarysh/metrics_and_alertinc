package service

import (
	"context"
	"fmt"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/entity"
)

type sender interface {
	SendPostCompressJSON(ctx context.Context, url string, storage entity.Metric)
}

func (m *MetricAction) SendMetric(ctx context.Context, flagRunAddr string) error {
	allMetric, err := m.storage.GetAllMetric()
	if err != nil {
		return fmt.Errorf("can't get all metrics: %w", err)
	}
	for _, value := range allMetric {
		metricType := value.MType
		switch metricType {
		case entity.Gauge:
			url := fmt.Sprintf("http://%s/update/%s/%s/%.2f", flagRunAddr, value.MType, value.ID, value.Value)
			m.sender.SendPostCompressJSON(ctx, url, value)
		case entity.Counter:
			url := fmt.Sprintf("http://%s/update/%s/%s/%d", flagRunAddr, value.MType, value.ID, value.Delta)
			m.sender.SendPostCompressJSON(ctx, url, value)
		}
	}
	return nil
}
