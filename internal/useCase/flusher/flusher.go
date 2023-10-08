package flusher

import (
	"context"
	"fmt"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/repository/storage"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/service"
	"time"
)

type Flusher struct {
	getMetric  *storage.MemStorage
	fileEngine *service.FileEngine
}

func NewFlusher(metric *storage.MemStorage, fileEngine *service.FileEngine) *Flusher {
	return &Flusher{
		getMetric:  metric,
		fileEngine: fileEngine,
	}
}

func (f *Flusher) Flush(ctx context.Context, interval uint64) {
	timeTicker := time.NewTicker(time.Second * time.Duration(interval))
	defer timeTicker.Stop()
	for {
		select {
		case <-timeTicker.C:
			f.SyncFlush()
		case <-ctx.Done():
			return
		}
	}
}

func (f *Flusher) SyncFlush() {
	data := f.getMetric.ReadMetric()
	_ = f.fileEngine.WriteFile(data)
}

func (f *Flusher) Restorer() error {
	data, err := f.fileEngine.ReadFile()
	if err != nil {
		return fmt.Errorf("read file error: %w", err)
	}
	f.getMetric.PutMetricMap(data)
	return nil
}
