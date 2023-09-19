package main

import (
	"context"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/logger"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/service"
	"go.uber.org/zap"
	"log"
	"net/http"

	"github.com/NikitaBarysh/metrics_and_alertinc/internal/config/server"

	"github.com/NikitaBarysh/metrics_and_alertinc/internal/flusher"

	"github.com/NikitaBarysh/metrics_and_alertinc/internal/restorer"

	"github.com/NikitaBarysh/metrics_and_alertinc/internal/handlers"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/router"
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

	memStorage := service.NewMemStorage()

	flush := flusher.NewFlusher(memStorage, file)
	flush.Restorer()

	if cfg.StoreInterval != 0 {
		go flush.Flush(ctx, cfg.StoreInterval)
	} else {
		memStorage.SetOnUpdate(flush.SyncFlush)
	}

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
