package handlers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type storage interface {
	UpdateGaugeMetric(key string, value float64)
	UpdateCounterMetric(key string, value int64)
	ReadGaugeMetric() map[string]float64
	ReadCounterMetric() map[string]int64
	GetAllMetric() []string
}

type Handler struct {
	storage storage
}

func NewHandler(storage storage) *Handler {
	return &Handler{
		storage,
	}
}

func (h *Handler) Router() *chi.Mux {
	r := chi.NewRouter()
	r.Post("/{update}/{type}/{name}/{value}", h.Safe)
	r.Get("/{value}/{type}/{name}", h.Get)
	r.Get("/", h.GetAll)
	return r
}

func (h *Handler) Safe(rw http.ResponseWriter, r *http.Request) {
	res := strings.Split(r.URL.Path, "/")

	if len(res) < 5 {
		http.Error(rw, "wrong request", http.StatusNotFound)
		return
	}

	update := chi.URLParam(r, "update")
	if update != "update" {
		http.Error(rw, "not update", http.StatusNotFound)
		return
	}

	metricType := chi.URLParam(r, "type")

	metricName := chi.URLParam(r, "name")

	metricValue := chi.URLParam(r, "value")

	fmt.Println(metricName, metricValue)

	switch metricType {
	case "counter":
		value, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			http.Error(rw, "wrong counter type", http.StatusBadRequest)
			return
		}
		h.storage.UpdateCounterMetric(metricName, value)
	case "gauge":
		value, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			http.Error(rw, "wrong gauge type", http.StatusBadRequest)
			return
		}
		h.storage.UpdateGaugeMetric(metricName, value)
	default:
		http.Error(rw, "unknown metric type", http.StatusNotImplemented)
		return
	}

	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	rw.WriteHeader(http.StatusOK)
}

func (h *Handler) Get(rw http.ResponseWriter, r *http.Request) {
	res := strings.Split(r.URL.Path, "/")

	if len(res) < 3 {
		http.Error(rw, "wrong request", http.StatusBadRequest)
	}

	metricMethod := chi.URLParam(r, "value")
	if metricMethod != "value" {
		http.Error(rw, "unknown method", http.StatusNotFound)
	}

	metricType := chi.URLParam(r, "type")
	metricName := chi.URLParam(r, "name")

	switch metricType {
	case "gauge":
		metricValue := h.storage.ReadGaugeMetric()
		if value, ok := metricValue[metricName]; ok {
			rw.WriteHeader(http.StatusOK)
			rw.Write([]byte(fmt.Sprintf("%v", value)))
			return
		} else {
			http.Error(rw, "wrong type", http.StatusNotFound)
			return
		}
	case "counter":
		metricValue := h.storage.ReadCounterMetric()
		if value, ok := metricValue[metricName]; ok {
			rw.WriteHeader(http.StatusOK)
			rw.Write([]byte(fmt.Sprintf("%v", value)))
			return
		} else {
			http.Error(rw, "wrong type", http.StatusNotFound)
			return
		}
	default:
		http.Error(rw, "unknown metric type", http.StatusNotFound)
		return
	}
}

func (h *Handler) GetAll(rw http.ResponseWriter, _ *http.Request) {
	list := h.storage.GetAllMetric()
	io.WriteString(rw, strings.Join(list, ","))
}
