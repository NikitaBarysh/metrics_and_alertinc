package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/NikitaBarysh/metrics_and_alertinc/config/server"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/interface/logger"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/repository"
	_ "github.com/NikitaBarysh/metrics_and_alertinc/internal/repository/postgres/migrations"
	"go.uber.org/zap"

	"github.com/NikitaBarysh/metrics_and_alertinc/internal/handlers"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/router"
	"github.com/go-chi/chi/v5"
)

func main() {
	cfg, configError := server.ParseServerConfig()
	if configError != nil {
		log.Fatalf("config err: %s\n", configError)
	}

	//ctx, cancel := context.WithCancel(context.Background())
	//defer cancel()

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

	handler := handlers.NewHandler(projectStorage, loggingVar)
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
