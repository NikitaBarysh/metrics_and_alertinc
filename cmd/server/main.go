package main

import (
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/handlers"
	"net/http"
)

func main() {
	err := http.ListenAndServe(`:8080`, http.HandlerFunc(handlers.Router))
	if err != nil {
		panic(err)
	}
}
