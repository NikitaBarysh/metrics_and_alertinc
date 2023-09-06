package main

import (
	"context"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/logger"
	"go.uber.org/zap"
	"log"
	"net/http"
	"time"

	"github.com/NikitaBarysh/metrics_and_alertinc/internal/config/server"

	"github.com/NikitaBarysh/metrics_and_alertinc/internal/restorer"

	"github.com/NikitaBarysh/metrics_and_alertinc/internal/handlers"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/router"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/storage/repositories"

	"github.com/go-chi/chi/v5"
)

func main() {
	cfg, configError := server.ParseServerConfig()
	if configError != nil {
		log.Fatalf("config err: %s\n", configError)
	}

	file := restorer.NewFileEngine(cfg.StorePath)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	loggingVar := logger.NewLoggingVar()
	loggingVar.Initialize(cfg.LogLevel)

	memStorage := repositories.NewMemStorage(file)

	go TimeTicker(ctx, cfg.StoreInterval, memStorage)

	handler := handlers.NewHandler(memStorage, loggingVar)
	router := router.NewRouter(handler)
	chiRouter := chi.NewRouter()
	chiRouter.Mount("/", router.Register())
	loggingVar.Log.Info("Running server", zap.String("address", cfg.RunAddr))
	err := http.ListenAndServe(cfg.RunAddr, chiRouter)
	if err != nil {
		panic(err)
	}
}

func TimeTicker(ctx context.Context, interval uint64, storage *repositories.MemStorage) {
	if interval < 1 {
		interval = 1
	}
	ticker := time.NewTicker(time.Second * time.Duration(interval))
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			storage.SaveData()
		case <-ctx.Done():
			return
		}
	}
}
