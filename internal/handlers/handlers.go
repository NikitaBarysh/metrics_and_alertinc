package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/entity"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/interface/logger"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/interface/models"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strconv"
	"time"
)

type storage interface {
	GetAllMetric() ([]entity.Metric, error)
	GetMetric(key string) (entity.Metric, error)
	SetMetrics(metric []entity.Metric) error
	CheckPing(ctx context.Context) error
}

type Handler struct {
	storage storage
	logger  logger.LoggingVar
}

func NewHandler(storage storage, logger *logger.LoggingVar) *Handler {
	return &Handler{
		storage,
		*logger,
	}
}

func (h *Handler) Safe(rw http.ResponseWriter, r *http.Request) {

	metricType := chi.URLParam(r, "type")

	metricName := chi.URLParam(r, "name")

	metricValue := chi.URLParam(r, "value")

	metricSlice := make([]entity.Metric, 0, 35)

	switch metricType {
	case "counter":
		value, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			http.Error(rw, "wrong counter type", http.StatusBadRequest)
			return
		}
		delta, err := h.storage.GetMetric(metricName)
		if err != nil || delta.ID == "" {
			fmt.Println(fmt.Errorf("no metric PollCount yet: %w", err))
		}
		delta.Delta += value
		metric := entity.Metric{ID: metricName, MType: metricType, Delta: value}

		metricSlice = append(metricSlice, metric)
	case "gauge":
		value, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			http.Error(rw, "wrong gauge type", http.StatusBadRequest)
			return
		}
		metric := entity.Metric{ID: metricName, MType: metricType, Value: value}

		metricSlice = append(metricSlice, metric)

	default:
		http.Error(rw, "unknown metric type", http.StatusNotImplemented)
		return
	}
	err := h.storage.SetMetrics(metricSlice)
	if err != nil {
		fmt.Println(fmt.Errorf("handlers: safe: SetMetric: %w", err))
	}

	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	rw.WriteHeader(http.StatusOK)
}

func (h *Handler) SafeBatch(rw http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "`apllication/json" {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	metricSlice := make([]entity.Metric, 0, 30)
	if err = json.Unmarshal(body, &metricSlice); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = h.storage.SetMetrics(metricSlice); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)

}

func (h *Handler) Get(rw http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "type")

	metricName := chi.URLParam(r, "name")

	metricValueStruct, err := h.storage.GetMetric(metricName)

	if err != nil || errors.Is(err, models.ErrNotFound) {
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}

	switch metricType {
	case "gauge":
		metricValue := metricValueStruct.Value
		rw.WriteHeader(http.StatusOK)
		_, err := rw.Write([]byte(fmt.Sprintf("%v", metricValue)))
		if err != nil {
			fmt.Println(fmt.Errorf("handler: get: write gauge metric: %w", err))
		}
		return
	case "counter":
		metricValue := metricValueStruct.Delta
		rw.WriteHeader(http.StatusOK)
		_, err := rw.Write([]byte(fmt.Sprintf("%v", metricValue)))
		if err != nil {
			fmt.Println(fmt.Errorf("handler: get: write counter metric: %w", err))
		}
		return
	default:
		http.Error(rw, "unknown metric type", http.StatusNotFound)
		return
	}
}

func (h *Handler) GetAll(rw http.ResponseWriter, _ *http.Request) {
	list, err := h.storage.GetAllMetric()
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		fmt.Println(fmt.Errorf("handler: getAll: method getallmetric: %w", err))
	}
	rw.Header().Set("Content-Type", "text/html")
	for _, value := range list {
		_, err := rw.Write([]byte(fmt.Sprintf(value.ID, value.MType, value.Value, value.Delta)))
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			fmt.Println(fmt.Errorf("handler: getAll: write metrices: %w", err)) // TODO
		}
	}
}

func (h *Handler) GetJSON(rw http.ResponseWriter, r *http.Request) {
	var req models.Metrics
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Println("decoder: ", err)
		h.logger.Log.Fatal("error decode getJSON", zap.Error(err))
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	switch req.MType {
	case "gauge":
		metricValue, err := h.storage.GetMetric(req.ID)
		if err != nil {
			http.Error(rw, "get json gauge error", http.StatusNotFound)
		}
		req.NewMetricValue(metricValue.Value)
	case "counter":
		metricValue, err := h.storage.GetMetric(req.ID)
		if err != nil {
			http.Error(rw, "get json counter error", http.StatusNotFound)
		}
		req.NewMetricDelta(metricValue.Delta)
	default:
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(rw).Encode(req); err != nil {
		fmt.Println("encoder: ", err)
		h.logger.Log.Fatal("error encoding getJSON", zap.Error(err))
	}
}

func (h *Handler) SafeJSON(rw http.ResponseWriter, r *http.Request) {
	var req models.Metrics
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Log.Debug("error  decode safeJSON", zap.Error(err))
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	metricSlice := make([]entity.Metric, 0, 35)

	switch req.MType {
	case "gauge":
		if req.Value == nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		metric := entity.Metric{ID: req.ID, MType: "gauge", Value: *req.Value}

		metricSlice = append(metricSlice, metric)
	case "counter":
		if req.Delta == nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		metric := entity.Metric{ID: req.ID, MType: "counter", Delta: *req.Delta}
		metricSlice = append(metricSlice, metric)
	default:
		rw.WriteHeader(http.StatusNotImplemented)
		return
	}
	h.storage.SetMetrics(metricSlice)

	rw.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(rw).Encode(req); err != nil {
		h.logger.Log.Debug("error encoding safeJSON", zap.Error(err))
	}
}

func (h *Handler) CheckConnection(rw http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	if err := h.storage.CheckPing(ctx); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Println(fmt.Errorf("handlers: CheckConnection: %w", err))
		return
	}

	rw.WriteHeader(http.StatusOK)
}
