package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/NikitaBarysh/metrics_and_alertinc/internal/entity"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/interface/logger"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/interface/models"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/repository/postgres"
	"go.uber.org/zap"

	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type storage interface {
	GetAllMetric() []entity.Metric
	GetMetric(key string) (entity.Metric, error)
	SetMetric(metric entity.Metric)
}

type Handler struct {
	storage storage
	logger  logger.LoggingVar
	db      *postgres.Postgres
}

func NewHandler(storage storage, logger *logger.LoggingVar, db *postgres.Postgres) *Handler {
	return &Handler{
		storage,
		*logger,
		db,
	}
}

func (h *Handler) Safe(rw http.ResponseWriter, r *http.Request) {
	var metric entity.Metric

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
		metric = entity.NewMetric(metricName, metricType, value, 0)
		h.storage.SetMetric(metric)
	case "gauge":
		value, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			http.Error(rw, "wrong gauge type", http.StatusBadRequest)
			return
		}
		metric = entity.NewMetric(metricName, metricType, 0, value)
		h.storage.SetMetric(metric)
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

	metricValueStruct, err := h.storage.GetMetric(metricName)
	if err != nil {
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
	h.storage.GetAllMetric()
	rw.Header().Set("Content-Type", "text/html")

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
	var req entity.Metric
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Log.Debug("error  decode safeJSON", zap.Error(err))
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	//metric = entity.Metric{ID: metricName, MType: metricType, Value: value}

	switch req.MType {
	case "gauge":
		if req.Value == 0 {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		req = entity.Metric{ID: req.ID, MType: req.MType, Value: req.Value}
		h.storage.SetMetric(req)
	case "counter":
		if req.Delta == 0 {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		req = entity.Metric{ID: req.ID, MType: req.MType, Delta: req.Delta}
		h.storage.SetMetric(req)
	default:
		rw.WriteHeader(http.StatusNotImplemented)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(rw).Encode(req); err != nil {
		h.logger.Log.Debug("error encoding safeJSON", zap.Error(err))
	}
}

func (h *Handler) CheckConnection(rw http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	err := h.db.CheckPing(ctx)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	rw.WriteHeader(http.StatusOK)
}
