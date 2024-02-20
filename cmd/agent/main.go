package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/NikitaBarysh/metrics_and_alertinc/config/agent"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/repository/memstorage"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/service"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/usecase"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/usecase/sender"
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

	cfg, err := agent.NewAgent()
	if err != nil {
		log.Fatalf("config err : %s\n", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	termSignal := make(chan os.Signal, 1)
	signal.Notify(termSignal, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	hash := usecase.WithHash(cfg)

	storage := memstorage.NewAgentStorage()

	send := sender.NewSender(hash)
	newMetricAction := service.NewMetricAction(storage, send, cfg)

	go newMetricAction.CollectPsutil(ctx)

	go newMetricAction.CollectRuntimeMetric(ctx)

	go newMetricAction.SendMetricsToServer(ctx, cfg)

	sig := <-termSignal
	fmt.Println("Agent Graceful Shutdown ", sig.String())
}
