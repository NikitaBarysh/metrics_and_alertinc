package grpc

import (
	"context"
	"fmt"

	"github.com/NikitaBarysh/metrics_and_alertinc/internal/entity"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/repository"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type Service struct {
	storage repository.Storage
}

func NewService(storage repository.Storage) Service {
	return Service{storage: storage}
}

func (s *Service) mustEmbedUnimplementedSendMetricServer() {
}

func (s *Service) Update(ctx context.Context, req *UpdateMetric) (*emptypb.Empty, error) {
	metricsSlice := make([]entity.Metric, 0, len(req.Metric))

	for _, metric := range req.Metric {
		switch metric.Type {
		case MType_Gauge:
			gaugeMetric := entity.Metric{}
			gaugeMetric.SetMetric(metric.ID, entity.Gauge, 0, metric.Value)
			metricsSlice = append(metricsSlice, gaugeMetric)
		case MType_Counter:
			counterMetric := entity.Metric{}
			counterMetric.SetMetric(metric.ID, entity.Gauge, metric.Delta, 0)
			metricsSlice = append(metricsSlice, counterMetric)
		}
	}

	err := s.storage.SetMetrics(metricsSlice)
	if err != nil {
		return nil, fmt.Errorf("err to set metrics: %w", err)
	}

	return &emptypb.Empty{}, nil
}
