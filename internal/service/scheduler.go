package service

import (
	"context"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/entity"
	"sync"
	"time"
)

func (m *MetricAction) SendMetricsToServer(ctx context.Context, reportInterval int64, flagRunAddr string, workers int64) {
	metricsCh := make(chan []entity.Metric, 1)
	var wg sync.WaitGroup
	wg.Add(int(workers))

	sendTicker := time.NewTicker(time.Second * time.Duration(reportInterval))
	defer sendTicker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-sendTicker.C:
			metric, _ := m.storage.GetAllMetric()
			metricsCh <- metric
			for i := 0; i <= int(workers); i++ {
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

func (m *MetricAction) WorkerPoll(ctx context.Context, flagAddr string, ch <-chan []entity.Metric) {
	for met := range ch {
		m.SendMetric(ctx, met, flagAddr)
	}
}
