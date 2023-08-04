package main

import (
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/server"

	"net/http"
)

func main() {
	err := http.ListenAndServe(`:8080`, http.HandlerFunc(server.Router))
	if err != nil {
		panic(err)
	}
}
