package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/NikitaBarysh/metrics_and_alertinc/internal/service/hasher"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/NikitaBarysh/metrics_and_alertinc/config/server"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/interface/logger"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/repository"
	_ "github.com/NikitaBarysh/metrics_and_alertinc/migrations"
	"go.uber.org/zap"

	"github.com/NikitaBarysh/metrics_and_alertinc/internal/handlers"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/router"
	"github.com/go-chi/chi/v5"
)

func main() {
	env, envErr := server.NewServer()
	if envErr != nil {
		log.Fatalf("config err: %s\n", envErr)
	}

	cfg, configError := server.NewConfig(env)
	if configError != nil {
		log.Fatalf("config err: %s\n", configError)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	termSig := make(chan os.Signal, 1)
	signal.Notify(termSig, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	loggingVar := logger.NewLoggingVar()
	loggerError := loggingVar.Initialize(cfg.LogLevel)
	if loggerError != nil {
		fmt.Println(fmt.Errorf("server: main: logger: %w", loggerError))
	}

	projectStorage, err := repository.New(ctx, cfg)
	if err != nil {
		panic(err)
	}

	handler := handlers.NewHandler(projectStorage, loggingVar)
	router := router.NewRouter(handler)
	chiRouter := chi.NewRouter()
	if cfg.Key != "" {
		hasher.Sign = hasher.NewHasher([]byte(cfg.Key))
		chiRouter.Use(hasher.Middleware)
	}
	chiRouter.Mount("/debug", middleware.Profiler())
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
