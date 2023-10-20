package router

import (
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/handlers"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/interface/compress"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/interface/logger"
	"github.com/go-chi/chi/v5"
)

type Router struct {
	metricHandler *handlers.Handler
}

func NewRouter(handler *handlers.Handler) *Router {
	return &Router{
		metricHandler: handler,
	}
}

func (rt *Router) Register() *chi.Mux {
	r := chi.NewRouter()
	r.Use(logger.WithLogging)
	r.Use(compress.GzipMiddleware)

	r.Post("/updates", rt.metricHandler.SafeBatch)
	r.Get("/ping", rt.metricHandler.CheckConnection)
	r.Post("/update/", rt.metricHandler.SafeJSON)
	r.Post("/update/{type}/{name}/{value}", rt.metricHandler.Safe)
	r.Get("/value/{type}/{name}", rt.metricHandler.Get)
	r.Get("/", rt.metricHandler.GetAll)
	r.Post("/value/", rt.metricHandler.GetJSON)

	return r
}
