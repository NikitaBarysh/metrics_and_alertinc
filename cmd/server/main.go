package main

import (
	"context"
	"fmt"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/interface/config/server"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/interface/logger"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/repository/postgres"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/repository/storage"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/service"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/useCase/flusher"
	"go.uber.org/zap"
	"log"
	"net/http"

	"github.com/NikitaBarysh/metrics_and_alertinc/internal/handlers"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/router"
	"github.com/go-chi/chi/v5"
)

func main() {
	cfg, configError := server.ParseServerConfig()
	if configError != nil {
		log.Fatalf("config err: %s\n", configError)
	}

	file := service.NewFileEngine(cfg.StorePath)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	loggingVar := logger.NewLoggingVar()
	loggerError := loggingVar.Initialize(cfg.LogLevel)
	if loggerError != nil {
		fmt.Println(fmt.Errorf("server: main: logger: %w", loggerError))
	}

	memStorage := storage.NewMemStorage()

	flush := flusher.NewFlusher(memStorage, file)
	restorerError := flush.Restorer()
	if restorerError != nil {
		fmt.Println(fmt.Errorf("server: main: restorer: %w", restorerError))
	}

	if cfg.StoreInterval != 0 {
		go flush.Flush(ctx, cfg.StoreInterval)
	} else {
		memStorage.SetOnUpdate(flush.SyncFlush)
	}
	db, err := postgres.InitPostgres(cfg)
	if err != nil {
		fmt.Println(fmt.Errorf("can't connect: %w", err))
	}

	handler := handlers.NewHandler(memStorage, loggingVar, db)
	router := router.NewRouter(handler)
	chiRouter := chi.NewRouter()
	chiRouter.Mount("/", router.Register())
	loggingVar.Log.Info("Running server", zap.String("address", cfg.RunAddr))
	err = http.ListenAndServe(cfg.RunAddr, chiRouter)
	if err != nil {
		panic(err)
	}
}
