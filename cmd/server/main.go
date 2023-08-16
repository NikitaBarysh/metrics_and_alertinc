package main

import (
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/logger"
	"go.uber.org/zap"
	"net/http"

	"github.com/NikitaBarysh/metrics_and_alertinc/internal/handlers"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/router"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/storage/repositories"
	"github.com/go-chi/chi/v5"
)

func main() {
	parseFlags()

	logger.Initialize(flagLogLevel)

	memStorage := repositories.NewMemStorage()
	handler := handlers.NewHandler(memStorage)
	router := router.NewRouter(handler)
	chiRouter := chi.NewRouter()
	chiRouter.Mount("/", router.Register())
	logger.Log.Info("Running server", zap.String("address", flagRunAddr))
	err := http.ListenAndServe(flagRunAddr, chiRouter)
	if err != nil {
		panic(err)
	}
}
