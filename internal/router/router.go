package router

import (
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/handlers"
	"github.com/go-chi/chi/v5"
)

type Router struct {
	metricHandler *handlers.Handler // TODO interface metricHandler
}

func NewRouter(handler *handlers.Handler) *Router {
	return &Router{
		metricHandler: handler,
	}
}

func (rt *Router) Register() *chi.Mux {
	r := chi.NewRouter()
	r.Post("/{update}/{type}/{name}/{value}", rt.metricHandler.Safe)
	r.Get("/{value}/{type}/{name}", rt.metricHandler.Get)
	r.Get("/", rt.metricHandler.GetAll)
	return r
}