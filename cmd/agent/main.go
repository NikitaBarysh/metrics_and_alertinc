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
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
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
	newMetricAction := service.NewMetricAction(storage, send)

	go newMetricAction.CollectPsutil(ctx, cfg.PollInterval)

	go newMetricAction.CollectRuntimeMetric(ctx, cfg.PollInterval)

	go newMetricAction.SendMetricsToServer(ctx, cfg.ReportInterval, cfg.URL, cfg.Limit)

	sig := <-termSignal
	fmt.Println("Agent Graceful Shutdown", sig.String())
}
