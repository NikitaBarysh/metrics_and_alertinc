package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/NikitaBarysh/metrics_and_alertinc/internal/encrypt"
	grpc2 "github.com/NikitaBarysh/metrics_and_alertinc/internal/grpc"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/service/hasher"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/usecase"
	"github.com/go-chi/chi/v5/middleware"
	"google.golang.org/grpc"

	"github.com/NikitaBarysh/metrics_and_alertinc/config/server"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/interface/logger"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/repository"
	_ "github.com/NikitaBarysh/metrics_and_alertinc/migrations"
	"go.uber.org/zap"

	"github.com/NikitaBarysh/metrics_and_alertinc/internal/handlers"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/router"
	"github.com/go-chi/chi/v5"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)

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
	if cfg.CryptoKey != "" {
		if err := encrypt.InitializeDecryptor(cfg.CryptoKey); err != nil {
			loggingVar.Error("err to create encryptor")
		}
		chiRouter.Use(encrypt.Middleware)
	}
	if cfg.TrustedSubnet != "" {
		if _, err := usecase.InitIPChecker(cfg.TrustedSubnet); err != nil {
			loggingVar.Error("err to init trusted subnet")
		}
		chiRouter.Use(usecase.Middleware)
	}

	chiRouter.Mount("/debug", middleware.Profiler())
	chiRouter.Mount("/", router.Register())
	loggingVar.Log.Info("Running server", zap.String("address", cfg.RunAddr))
	if cfg.ServerType == "http" {
		go func() {
			err = http.ListenAndServe(cfg.RunAddr, chiRouter)
			if err != nil {
				panic(err)
			}
		}()
	} else if cfg.ServerType == "grpc" {
		s := grpc.NewServer()
		service := grpc2.NewService(projectStorage)
		grpc2.RegisterSendMetricServer(s, &service)
		listen, err := net.Listen("tcp", cfg.RunAddr)
		if err != nil {
			log.Fatalf("err to start grpc server: %w", err)
		}
		if err = s.Serve(listen); err != nil {
			panic(err)
		}

	}

	sig := <-termSig
	loggingVar.Log.Info("Server Graceful Shutdown", zap.String("-", sig.String()))
}
