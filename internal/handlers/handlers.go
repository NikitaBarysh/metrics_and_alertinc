package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
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

func (h *Handler) Safe(rw http.ResponseWriter, r *http.Request) {

	update := chi.URLParam(r, "update")
	if update != "update" {
		http.Error(rw, "not update", http.StatusNotFound)
		return
	}

	metricType := chi.URLParam(r, "type")

	metricName := chi.URLParam(r, "name")

	metricValue := chi.URLParam(r, "value")

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
		}
		http.Error(rw, "wrong type", http.StatusNotFound)
	case "counter":
		metricValue := h.storage.ReadCounterMetric()
		if value, ok := metricValue[metricName]; ok {
			rw.WriteHeader(http.StatusOK)
			rw.Write([]byte(fmt.Sprintf("%v", value)))
			return
		}
		http.Error(rw, "wrong type", http.StatusNotFound)
	default:
		http.Error(rw, "unknown metric type", http.StatusNotFound)
		return
	}
}

func (h *Handler) GetAll(rw http.ResponseWriter, _ *http.Request) {
	list := h.storage.GetAllMetric()
	io.WriteString(rw, strings.Join(list, ","))
}
