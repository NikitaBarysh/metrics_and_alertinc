package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/logger"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/models"
	"go.uber.org/zap"
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

func (h *Handler) GetJSON(rw http.ResponseWriter, r *http.Request) {
	var req models.Metrics
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log.Debug("error decode getJSON", zap.Error(err))
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	switch req.MType {
	case "gauge":
		metricValue := h.storage.ReadGaugeMetric()
		if value, ok := metricValue[req.ID]; ok {
			req.NewMetricValue(value)
		} else {
			rw.WriteHeader(http.StatusNotFound)
		}
	case "counter":
		metricValue := h.storage.ReadCounterMetric()
		if value, ok := metricValue[req.ID]; ok {
			req.NewMetricDelta(value)
		} else {
			rw.WriteHeader(http.StatusNotFound)
		}
	default:
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(rw).Encode(req); err != nil {
		logger.Log.Debug("error encoding getJSON", zap.Error(err))
	}
}

func (h *Handler) SafeJSON(rw http.ResponseWriter, r *http.Request) {
	var req models.Metrics
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log.Debug("error  decode safeJSON", zap.Error(err))
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	switch req.MType {
	case "gauge":
		if req.Value == nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		h.storage.UpdateGaugeMetric(req.ID, *req.Value)
	case "counter":
		if req.Delta == nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		h.storage.UpdateCounterMetric(req.ID, *req.Delta)
	default:
		rw.WriteHeader(http.StatusNotImplemented)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(rw).Encode(req); err != nil {
		logger.Log.Debug("error encoding safeJSON", zap.Error(err))
	}
}
