// Package service - содержит внутреннею логику приложения
package service

import (
	"context"
	"sync"
	"time"

	"github.com/NikitaBarysh/metrics_and_alertinc/internal/entity"
)

// SendMetricsToServer - планировщик получение и отправки метрик
func (m *MetricAction) SendMetricsToServer(ctx context.Context, reportInterval int64, flagRunAddr string, workers int) {
	metricsCh := make(chan []entity.Metric, 1)
	var wg sync.WaitGroup
	wg.Add(workers)

	sendTicker := time.NewTicker(time.Second * time.Duration(reportInterval))
	defer sendTicker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-sendTicker.C:
			metric, _ := m.storage.GetAllMetric()
			metricsCh <- metric
			for i := 0; i <= workers; i++ {
				go func() {
					defer wg.Done()
					m.WorkerPoll(ctx, flagRunAddr, metricsCh)
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
func (m *MetricAction) CollectPsutil(ctx context.Context, pollInterval int64) {
	collectTicker := time.NewTicker(time.Second * time.Duration(pollInterval))
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
func (m *MetricAction) CollectRuntimeMetric(ctx context.Context, pollInterval int64) {
	collectTicker := time.NewTicker(time.Second * time.Duration(pollInterval))
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
func (m *MetricAction) WorkerPoll(ctx context.Context, flagAddr string, ch <-chan []entity.Metric) {
	for met := range ch {
		m.SendMetric(ctx, met, flagAddr)
	}
}
