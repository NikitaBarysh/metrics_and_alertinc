package main

import (
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/config/serverConfig"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/logger"
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

	logger.Initialize(cfg.LogLevel)

	memStorage := repositories.NewMemStorage()
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
