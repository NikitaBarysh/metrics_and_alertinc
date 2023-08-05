package main

import (
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/storage"

	"github.com/NikitaBarysh/metrics_and_alertinc/internal/server"
)

func main() {
	CreateMemStorage := storage.CreateMemStorage()
	MemStorageAction := server.MemStorageAction{CreateMemStorage}
	go MemStorageAction.Run()

}
