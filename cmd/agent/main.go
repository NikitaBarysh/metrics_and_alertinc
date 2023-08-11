package main

import (
	"context"
	"fmt"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/sender"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/storage"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/storage/repositories"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	parseFlags()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	termSignal := make(chan os.Signal, 1)
	signal.Notify(termSignal, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	memStorage := repositories.NewMemStorage()
	sender := sender.NewSender()
	newMetricAction := storage.NewMetricAction(memStorage, sender)
	go newMetricAction.Run(ctx)

	sig := <-termSignal
	fmt.Println(sig.String())
}
