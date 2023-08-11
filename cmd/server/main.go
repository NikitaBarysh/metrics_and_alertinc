package main

import (
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/storage/repositories"
	"github.com/go-chi/chi/v5"
	"net/http"

	"github.com/NikitaBarysh/metrics_and_alertinc/internal/handlers"
)

func main() {
	parseFlag()

	memStorage := repositories.NewMemStorage()
	handler := handlers.NewHandler(memStorage)
	//router := router.NewRouter(handler)
	chiRouter := chi.NewRouter()
	chiRouter.Mount("/", handler.Router())
	err := http.ListenAndServe(url, chiRouter)
	if err != nil {
		panic(err)
	}
}
