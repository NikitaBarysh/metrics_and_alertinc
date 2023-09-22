package service

import (
	"context"
	"time"
)

func (m *MetricAction) Run(ctx context.Context, pollInterval int64, reportInterval int64, flagRunAddr string) error {

	collectTicker := time.NewTicker(time.Second * time.Duration(pollInterval))
	defer collectTicker.Stop()

	sendTicker := time.NewTicker(time.Second * time.Duration(reportInterval))
	defer sendTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-collectTicker.C:
			m.CollectMetric()
		case <-sendTicker.C:
			m.SendMetric(ctx, flagRunAddr) // TODO
		}
	}
}
