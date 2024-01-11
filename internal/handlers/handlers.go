package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	_ "net/http/pprof"
	"strconv"

	"github.com/NikitaBarysh/metrics_and_alertinc/internal/entity"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/interface/logger"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/interface/models"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

//go:generate mockgen -source ${GOFILE} -destination mocks_test.go -package ${GOPACKAGE}

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

	mType := entity.MType(metricType)

	metricName := chi.URLParam(r, "name")

	metricValue := chi.URLParam(r, "value")

	metricSlice := make([]entity.Metric, 0, 35)

	switch mType {
	case entity.Counter:
		value, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			http.Error(rw, "wrong counter type", http.StatusBadRequest)
			return
		}
		delta, errGet := h.storage.GetMetric(metricName)
		if errGet != nil || delta.ID == "" {
			h.logger.Error(fmt.Errorf("no metric PollCount yet: %w", err).Error())
		}
		delta.Delta += value
		metric := entity.Metric{ID: metricName, MType: mType, Delta: value}

		metricSlice = append(metricSlice, metric)
	case entity.Gauge:
		value, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			http.Error(rw, "wrong gauge type", http.StatusBadRequest)
			return
		}
		metric := entity.Metric{ID: metricName, MType: mType, Value: value}

		metricSlice = append(metricSlice, metric)

	default:
		http.Error(rw, "unknown metric type", http.StatusNotImplemented)
		return
	}
	err := h.storage.SetMetrics(metricSlice)
	if err != nil {
		h.logger.Error(fmt.Errorf("handlers: safe: SetMetric: %w", err).Error())
	}

	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	rw.WriteHeader(http.StatusOK)
}

func (h *Handler) SafeBatch(rw http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
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
	mType := entity.MType(metricType)

	metricName := chi.URLParam(r, "name")

	metricValueStruct, err := h.storage.GetMetric(metricName)

	if err != nil || errors.Is(err, models.ErrNotFound) {
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}

	switch mType {
	case entity.Gauge:
		metricValue := metricValueStruct.Value
		rw.WriteHeader(http.StatusOK)
		_, err := rw.Write([]byte(fmt.Sprintf("%v", metricValue)))
		if err != nil {
			h.logger.Error(fmt.Errorf("handler: get: write gauge metric: %w", err).Error())
		}
		return
	case entity.Counter:
		metricValue := metricValueStruct.Delta
		rw.WriteHeader(http.StatusOK)
		_, err := rw.Write([]byte(fmt.Sprintf("%v", metricValue)))
		if err != nil {
			h.logger.Error(fmt.Errorf("handler: get: write counter metric: %w", err).Error())
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
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()
	if err := h.storage.CheckPing(ctx); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		h.logger.Error(fmt.Errorf("handlers: CheckConnection: %w", err).Error())
		return
	}

	rw.WriteHeader(http.StatusOK)
}
