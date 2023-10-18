package main

import (
	"context"
	"fmt"
	"github.com/NikitaBarysh/metrics_and_alertinc/config/agent"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/repository/storage"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/service"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/useCase/sender"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	cfg, err := agent.ParseAgentFlags()
	if err != nil {
		log.Fatalf("config err : %s\n", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	termSignal := make(chan os.Signal, 1)
	signal.Notify(termSignal, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	memStorage := storage.NewAgentStorage()

	send := sender.NewSender()
	newMetricAction := service.NewMetricAction(memStorage, send)

	go newMetricAction.Run(ctx, cfg.PollInterval, cfg.ReportInterval, cfg.URL) // TODO

	sig := <-termSignal
	fmt.Println("Agent Graceful Shutdown", sig.String())
}
