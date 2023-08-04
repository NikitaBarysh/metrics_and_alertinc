package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/NikitaBarysh/metrics_and_alertinc/internal/storage"

	"github.com/NikitaBarysh/metrics_and_alertinc/internal/server"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Обрабатывем сигналы от системы.
	termSignal := make(chan os.Signal, 1)
	signal.Notify(termSignal, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	CreateMemStorage := storage.CreateMemStorage()
	MemStorageAction := server.MemStorageAction{CreateMemStorage}
	go MemStorageAction.Run(ctx)
	
	sig := <-termSignal
	fmt.Println("end")
	fmt.Println(sig.String())
}
