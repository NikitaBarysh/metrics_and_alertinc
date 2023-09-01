package main

import (
	"context"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/config/serverConfig"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/flusher"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/logger"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/restorer"
	"go.uber.org/zap"
	"log"
	"net/http"

	"github.com/NikitaBarysh/metrics_and_alertinc/internal/handlers"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/router"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/storage/repositories"
	"github.com/go-chi/chi/v5"
)

func main() {
	cfg, configError := serverConfig.ParseServerConfig()
	if configError != nil {
		log.Fatalf("config err: %s\n", configError)
	}

	file := restorer.NewFileEngine(cfg.StorePath)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger.Initialize(cfg.LogLevel)

	memStorage := repositories.NewMemStorage()

	flush := flusher.NewFlusher(memStorage, file)
	flush.Restorer()

	if cfg.StoreInterval != 0 {
		go flush.Flush(ctx, cfg.StoreInterval)
	} else {
		memStorage.SetOnUpdate(flush.SyncFlush)
	}

	handler := handlers.NewHandler(memStorage)
	router := router.NewRouter(handler)
	chiRouter := chi.NewRouter()
	chiRouter.Mount("/", router.Register())
	logger.Log.Info("Running serverConfig", zap.String("address", cfg.RunAddr))
	err := http.ListenAndServe(cfg.RunAddr, chiRouter)
	if err != nil {
		panic(err)
	}
}
