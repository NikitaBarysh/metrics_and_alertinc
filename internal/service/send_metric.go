package service

import (
	"context"
	"fmt"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/entity"
)

type sender interface {
	SendPost(ctx context.Context, url string, storage entity.Metric)
}

func (m *MetricAction) SendMetric(ctx context.Context, flagRunAddr string) error {
	for _, value := range m.MemStorage.GetAllMetric() {
		metricType := value.MType
		switch metricType {
		case "gauge":
			url := fmt.Sprintf("http://%s/update/%s/%s/%.2f", flagRunAddr, value.MType, value.ID, value.Value)
			m.sender.SendPost(ctx, url, value)
		case "counter":
			url := fmt.Sprintf("http://%s/update/%s/%s/%d", flagRunAddr, value.MType, value.ID, value.Delta)
			m.sender.SendPost(ctx, url, value)
		}
	}
	return nil
}
