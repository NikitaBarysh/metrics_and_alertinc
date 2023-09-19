package main

import (
	"context"
	"fmt"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/config/agent"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/service"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/useCase/sender"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/NikitaBarysh/metrics_and_alertinc/internal/storage"
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

	memStorage := service.NewMemStorage()
	sender := sender.NewSender()
	newMetricAction := storage.NewMetricAction(memStorage, sender)
	go newMetricAction.Run(ctx, cfg.PollInterval, cfg.ReportInterval, cfg.URL)

	sig := <-termSignal
	fmt.Println(sig.String())
}
