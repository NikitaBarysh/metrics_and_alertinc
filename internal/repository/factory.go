package repository

import (
	"context"
	"fmt"

	"github.com/NikitaBarysh/metrics_and_alertinc/config/server"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/entity"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/repository/filestorage"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/repository/memstorage"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/repository/postgres"
)

type Storage interface {
	GetAllMetric() ([]entity.Metric, error)
	GetMetric(key string) (entity.Metric, error)
	SetMetrics(metric []entity.Metric) error
	CheckPing(ctx context.Context) error
}

// New - выбираем куда будем складывать метрики и откуда получать
func New(ctx context.Context, cfg *server.Config) (Storage, error) {
	if cfg.DataBaseDSN != "" {
		return postgres.InitPostgres(cfg)
	} else if cfg.StorePath != "" {
		file, err := filestorage.NewFileEngine(cfg.StorePath)
		if err != nil {
			fmt.Println("memstorage-file error factory")
			return nil, err
		}
		memStorage, err := memstorage.NewMemStorage(ctx, cfg, file)
		if err != nil {
			return nil, err
		}
		return memStorage, nil
	} else {
		memStorage, err := memstorage.NewMemStorage(ctx, cfg, nil)
		if err != nil {
			fmt.Println("memstorage error factory")
			return nil, err
		}
		return memStorage, nil
	}
}
