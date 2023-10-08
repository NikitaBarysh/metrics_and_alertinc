package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/entity"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/interface/logger"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/interface/models"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/repository/postgres"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
)

type storage interface {
	UpdateGaugeMetric(key string, value float64)
	UpdateCounterMetric(key string, value int64)
	ReadMetric() map[string]entity.Metric
	GetAllMetric() []string
	ReadDefinitelyMetric(key string) (entity.Metric, error)
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

	metricValueStruct, err := h.storage.ReadDefinitelyMetric(metricName)
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
	list := h.storage.GetAllMetric()
	rw.Header().Set("Content-Type", "text/html")
	_, err := io.WriteString(rw, strings.Join(list, ","))
	if err != nil {
		fmt.Println(fmt.Errorf("handler: getAll: write metrices: %w", err))
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
		metricValue, err := h.storage.ReadDefinitelyMetric(req.ID)
		if err != nil {
			http.Error(rw, "get json gauge error", http.StatusNotFound)
		}
		req.NewMetricValue(metricValue.Value)
	case "counter":
		metricValue, err := h.storage.ReadDefinitelyMetric(req.ID)
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
		h.logger.Log.Debug("error encoding safeJSON", zap.Error(err))
	}
}

func (h *Handler) CheckConnection(rw http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	if err := h.db.CheckPing(ctx); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Println(fmt.Errorf("handlers: CheckConnection: %w", err))
		return
	}

	rw.WriteHeader(http.StatusOK)
}
