// Package service - содержит внутреннею логику приложения
package service

import (
	"context"
	"fmt"

	"github.com/NikitaBarysh/metrics_and_alertinc/internal/entity"
)

type sender interface {
	SendPostCompressJSON(ctx context.Context, url string, storage entity.Metric, ip string)
}

// SendMetric - подготовка для отправки метрик на сервер
func (m *MetricAction) SendMetric(ctx context.Context, allMetric []entity.Metric, flagRunAddr string, ip string) error {
	for _, value := range allMetric {
		metricType := value.MType
		switch metricType {
		case entity.Gauge:
			url := fmt.Sprintf("http://%s/update/%s/%s/%.2f", flagRunAddr, value.MType, value.ID, value.Value)
			m.sender.SendPostCompressJSON(ctx, url, value, ip)
		case entity.Counter:
			url := fmt.Sprintf("http://%s/update/%s/%s/%d", flagRunAddr, value.MType, value.ID, value.Delta)
			m.sender.SendPostCompressJSON(ctx, url, value, ip)
		}
	}
	return nil
}
