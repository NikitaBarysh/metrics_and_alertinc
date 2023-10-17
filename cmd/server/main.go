package main

import (
	"context"
	"fmt"
	"github.com/NikitaBarysh/metrics_and_alertinc/config/server"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/interface/logger"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/repository"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/repository/postgres"
	_ "github.com/NikitaBarysh/metrics_and_alertinc/internal/repository/postgres/migrations"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/repository/storage"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/service"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/useCase/flusher"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

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

	termSig := make(chan os.Signal, 1)
	signal.Notify(termSig, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	loggingVar := logger.NewLoggingVar()
	loggerError := loggingVar.Initialize(cfg.LogLevel)
	if loggerError != nil {
		fmt.Println(fmt.Errorf("server: main: logger: %w", loggerError))
	}

	projectStorage, err := repository.New(cfg)
	if err != nil {
		panic(err)
	}

	memStorage, err := storage.NewMemStorage()
	if err != nil {
		panic(err)
	}

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

	handler := handlers.NewHandler(projectStorage, loggingVar, db)
	router := router.NewRouter(handler)
	chiRouter := chi.NewRouter()
	chiRouter.Mount("/", router.Register())
	loggingVar.Log.Info("Running server", zap.String("address", cfg.RunAddr))
	go func() {
		err = http.ListenAndServe(cfg.RunAddr, chiRouter)
		if err != nil {
			panic(err)
		}
	}()

	sig := <-termSig
	fmt.Println("Server Graceful Shutdown", sig.String())
}
