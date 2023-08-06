package main

import (
	"context"
	"fmt"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/storage"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NikitaBarysh/metrics_and_alertinc/internal/server"
)

func main() {
	time.Sleep(time.Second * 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	termSignal := make(chan os.Signal, 1)
	signal.Notify(termSignal, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	createMemStorage := storage.CreateMemStorage()
	memStorageAction := server.MemStorageAction{MemStorage: createMemStorage}
	go memStorageAction.Run(ctx)

	sig := <-termSignal
	fmt.Println(sig.String())
}
