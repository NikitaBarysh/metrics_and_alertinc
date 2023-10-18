package repository

import (
	"context"
	"fmt"
	"github.com/NikitaBarysh/metrics_and_alertinc/config/server"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/entity"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/repository/file_storage"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/repository/postgres"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/repository/storage"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/useCase/flusher"
)

type Storage interface {
	GetAllMetric() ([]entity.Metric, error)
	GetMetric(key string) (entity.Metric, error)
	SetMetrics(metric []entity.Metric) error
}

func New(cfg *server.Config) (Storage, error) {
	if cfg.DataBaseDSN != "" {
		return postgres.InitPostgres(cfg)
	} else if cfg.StorePath != "" {
		return file_storage.NewFileEngine(cfg.StorePath)
	} else {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		memStorage, err := storage.NewMemStorage(cfg)
		if err != nil {
			return nil, fmt.Errorf("can't create memStorage: %w", err)
		}
		flush := flusher.NewFlusher(memStorage)
		restorerError := flush.Restorer()
		if restorerError != nil {
			fmt.Println(fmt.Errorf("server: main: restorer: %w", restorerError))
		}

		if cfg.StoreInterval != 0 {
			go flush.Flush(ctx, cfg.StoreInterval)
		} else {
			memStorage.SetOnUpdate(flush.SyncFlush)
		}
		return memStorage, nil
	}
}
