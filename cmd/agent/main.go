package main

import (
	"context"
	"fmt"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/config/agentConfig"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/NikitaBarysh/metrics_and_alertinc/internal/sender"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/storage"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/storage/repositories"
)

func main() {

	cfg, err := agentConfig.ParseAgentFlags()
	if err != nil {
		log.Fatalf("config err : %s\n", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	termSignal := make(chan os.Signal, 1)
	signal.Notify(termSignal, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	memStorage := repositories.NewMemStorage()
	sender := sender.NewSender()
	newMetricAction := storage.NewMetricAction(memStorage, sender)
	go newMetricAction.Run(ctx, cfg.PollInterval, cfg.ReportInterval, cfg.URL)

	sig := <-termSignal
	fmt.Println(sig.String())
}
