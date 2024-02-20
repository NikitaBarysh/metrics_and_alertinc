// Package service - содержит внутреннею логику приложения
package service

import (
	"context"
	"sync"
	"time"

	"github.com/NikitaBarysh/metrics_and_alertinc/config/agent"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/entity"
)

// SendMetricsToServer - планировщик получение и отправки метрик
func (m *MetricAction) SendMetricsToServer(ctx context.Context, cfg *agent.Config) {
	metricsCh := make(chan []entity.Metric, 1)
	var wg sync.WaitGroup
	wg.Add(cfg.Limit)

	sendTicker := time.NewTicker(time.Second * time.Duration(cfg.ReportInterval))
	defer sendTicker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-sendTicker.C:
			metric, _ := m.storage.GetAllMetric()
			metricsCh <- metric
			for i := 0; i <= cfg.Limit; i++ {
				go func() {
					defer wg.Done()
					m.WorkerPoll(ctx, metricsCh)
				}()
			}
			go func() {
				wg.Wait()
				close(metricsCh)
			}()
		}
	}
}

// CollectPsutil - планировщик сбора psutil метрик
func (m *MetricAction) CollectPsutil(ctx context.Context) {
	collectTicker := time.NewTicker(time.Second * time.Duration(m.cfg.PollInterval))
	defer collectTicker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-collectTicker.C:
			m.CollectPsutilMetrics()
		}
	}
}

// CollectRuntimeMetric - планировщик сбора runtime метрик
func (m *MetricAction) CollectRuntimeMetric(ctx context.Context) {
	collectTicker := time.NewTicker(time.Second * time.Duration(m.cfg.PollInterval))
	defer collectTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-collectTicker.C:
			m.CollectMetric()
		}
	}
}

// WorkerPoll - воркер, который отправляет метрики
func (m *MetricAction) WorkerPoll(ctx context.Context, ch <-chan []entity.Metric) {
	for met := range ch {
		m.SendMetric(ctx, met)
	}
}
