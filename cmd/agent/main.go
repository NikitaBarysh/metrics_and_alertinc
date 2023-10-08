package main

import (
	"context"
	"fmt"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/interface/config/agent"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/repository/storage"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/useCase/sender"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/NikitaBarysh/metrics_and_alertinc/internal/service"
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

	memStorage := storage.NewMemStorage()
	send := sender.NewSender()
	newMetricAction := service.NewMetricAction(memStorage, send)
	go newMetricAction.Run(ctx, cfg.PollInterval, cfg.ReportInterval, cfg.URL) // TODO

	sig := <-termSignal
	fmt.Println(sig.String())
}
