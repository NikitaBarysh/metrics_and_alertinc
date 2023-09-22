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
	for metricName, metricValue := range m.MemStorage.ReadMetric() {
		metricType := m.MemStorage.MetricMap[metricName].MType
		switch metricType {
		case "gauge":
			url := fmt.Sprintf("http://%s/update/%s/%s/%.2f", flagRunAddr, metricType, metricName, metricValue.Value)
			m.sender.SendPost(ctx, url, metricValue)
		case "counter":
			url := fmt.Sprintf("http://%s/update/%s/%s/%d", flagRunAddr, metricType, metricName, metricValue.Delta)
			m.sender.SendPost(ctx, url, metricValue)
		}
	}
	return nil
}
