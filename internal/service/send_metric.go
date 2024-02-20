// Package service - содержит внутреннею логику приложения
package service

import (
	"context"
	"fmt"

	"github.com/NikitaBarysh/metrics_and_alertinc/internal/entity"
	grpc2 "github.com/NikitaBarysh/metrics_and_alertinc/internal/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type sender interface {
	SendPostCompressJSON(ctx context.Context, url string, storage entity.Metric, ip string)
	SendGRPC(metrics []entity.Metric, ip string, grpcClient grpc2.SendMetricClient)
}

// SendMetric - подготовка для отправки метрик на сервер
func (m *MetricAction) SendMetric(ctx context.Context, allMetric []entity.Metric) error {
	if m.cfg.ServiceType == "http" {
		for _, value := range allMetric {
			metricType := value.MType
			switch metricType {
			case entity.Gauge:
				url := fmt.Sprintf("http://%s/update/%s/%s/%.2f", m.cfg.URL, value.MType, value.ID, value.Value)
				m.sender.SendPostCompressJSON(ctx, url, value, m.cfg.IP)
			case entity.Counter:
				url := fmt.Sprintf("http://%s/update/%s/%s/%d", m.cfg.URL, value.MType, value.ID, value.Delta)
				m.sender.SendPostCompressJSON(ctx, url, value, m.cfg.IP)
			}
		}
	} else if m.cfg.ServiceType == "grpc" {
		conn, err := grpc.Dial(m.cfg.URL, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return fmt.Errorf("err to dial grpc: %w", err)
		}
		c := grpc2.NewSendMetricClient(conn)

		m.sender.SendGRPC(allMetric, m.cfg.IP, c)
	}

	return nil
}
