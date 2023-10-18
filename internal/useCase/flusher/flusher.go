package flusher

import (
	"context"
	"fmt"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/repository/storage"
	"time"
)

type Flusher struct {
	memory *storage.MemStorage
}

func NewFlusher(metric *storage.MemStorage) *Flusher {
	return &Flusher{
		memory: metric,
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
	data, err := f.memory.GetAllMetric()
	if err != nil {
		fmt.Println(fmt.Errorf("can't get all metrics: flusher: syncFlush: %w", err)) //TODO
	}
	_ = f.memory.FileEngine.SetMetrics(data)
}

func (f *Flusher) Restorer() error {
	data, err := f.memory.FileEngine.GetAllMetric()
	if err != nil {
		return fmt.Errorf("read file error: %w", err)
	}
	err = f.memory.SetMetrics(data)
	if err != nil {
		return fmt.Errorf("can't set metrics: flusher: restorer: %w", err)
	}
	return nil
}
