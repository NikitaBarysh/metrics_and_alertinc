package main

import (
	"fmt"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/handlers"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/router"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/storage/repositories"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func main() {
	parseFlags()

	memStorage := repositories.NewMemStorage()
	handler := handlers.NewHandler(memStorage)
	router := router.NewRouter(handler)
	chiRouter := chi.NewRouter()
	chiRouter.Mount("/", router.Register())
	err := http.ListenAndServe(fmt.Sprintf("localhost%s", flagRunAddr), chiRouter)
	if err != nil {
		panic(err)
	}
}
