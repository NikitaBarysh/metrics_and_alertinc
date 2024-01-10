package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NikitaBarysh/metrics_and_alertinc/internal/entity"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/interface/logger"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHandler_Safe(t *testing.T) {
	type mockBehaviour func(s *Mockstorage, metric entity.Metric, metricSlice []entity.Metric, key string)
	type want struct {
		code int
		url  string
	}
	tests := []struct {
		name          string
		metric        entity.Metric
		mockBehaviour mockBehaviour
		metricSlice   []entity.Metric
		want          want
		param         map[string]any
	}{
		{
			name: "Test#1, check if all correct for counter",
			metric: entity.Metric{
				ID:    "PollCount",
				MType: "counter",
				Delta: 10,
				Value: 0,
			},
			metricSlice: []entity.Metric{{ID: "PollCount", MType: "counter", Delta: 10, Value: 0}},
			mockBehaviour: func(s *Mockstorage, metric entity.Metric, metricSlice []entity.Metric, key string) {
				s.EXPECT().GetMetric(key).Return(metric, nil)
				s.EXPECT().SetMetrics(metricSlice).Return(nil)
			},
			want:  want{http.StatusOK, "/update/counter/PollCount/10"},
			param: map[string]any{"update": "update", "type": "counter", "name": "PollCount", "value": "10"},
		},
		{
			name: "Test#2, check if all correct for gauge",
			metric: entity.Metric{
				ID:    "Alloc",
				MType: "gauge",
				Delta: 0,
				Value: 527,
			},
			metricSlice: []entity.Metric{{ID: "Alloc", MType: "gauge", Delta: 0, Value: 527}},
			mockBehaviour: func(s *Mockstorage, metric entity.Metric, metricSlice []entity.Metric, key string) {
				s.EXPECT().SetMetrics(metricSlice).Return(nil)
			},
			want:  want{http.StatusOK, "/update/gauge/Alloc/527"},
			param: map[string]any{"update": "update", "type": "gauge", "name": "Alloc", "value": "527"},
		},
		{
			name: "Test#3, err get counter metric",
			metric: entity.Metric{
				ID:    "PollCount",
				MType: "counter",
				Delta: 10,
				Value: 0,
			},
			metricSlice: []entity.Metric{{ID: "PollCount", MType: "counter", Delta: 10, Value: 0}},
			mockBehaviour: func(s *Mockstorage, metric entity.Metric, metricSlice []entity.Metric, key string) {
				s.EXPECT().GetMetric(key).Return(entity.Metric{}, errors.New("no metric PollCount yet"))
				s.EXPECT().SetMetrics(metricSlice).Return(nil)
			},
			want:  want{http.StatusOK, "/update/counter/PollCount/10"},
			param: map[string]any{"update": "update", "type": "counter", "name": "PollCount", "value": "10"},
		},
		{
			name: "Test#4, err set metric",
			metric: entity.Metric{
				ID:    "PollCount",
				MType: "counter",
				Delta: 10,
				Value: 0,
			},
			metricSlice: []entity.Metric{{ID: "PollCount", MType: "counter", Delta: 10, Value: 0}},
			mockBehaviour: func(s *Mockstorage, metric entity.Metric, metricSlice []entity.Metric, key string) {
				s.EXPECT().GetMetric(key).Return(metric, nil)
				s.EXPECT().SetMetrics(metricSlice).Return(errors.New("err to set metric"))
			},
			want:  want{http.StatusOK, "/update/counter/PollCount/10"},
			param: map[string]any{"update": "update", "type": "counter", "name": "PollCount", "value": "10"},
		},
		{
			name: "Test#5, wrong metric type",
			metric: entity.Metric{
				ID:    "PollCount",
				MType: "counter",
				Delta: 10,
				Value: 0,
			},
			metricSlice: []entity.Metric{{ID: "PollCount", MType: "counter", Delta: 10, Value: 0}},
			mockBehaviour: func(s *Mockstorage, metric entity.Metric, metricSlice []entity.Metric, key string) {
			},
			want:  want{http.StatusNotImplemented, "/update/count/PollCount/10"},
			param: map[string]any{"update": "update", "type": "count", "name": "PollCount", "value": "10"},
		},
		{
			name: "Test#6, wrong counter value",
			metric: entity.Metric{
				ID:    "PollCount",
				MType: "counter",
				Delta: 10,
				Value: 0,
			},
			metricSlice: []entity.Metric{{ID: "PollCount", MType: "counter", Delta: 10, Value: 0}},
			mockBehaviour: func(s *Mockstorage, metric entity.Metric, metricSlice []entity.Metric, key string) {
			},
			want:  want{http.StatusBadRequest, "/update/counter/PollCount/val"},
			param: map[string]any{"update": "update", "type": "counter", "name": "PollCount", "value": "val"},
		},
		{
			name: "Test#7, wrong gauge value",
			metric: entity.Metric{
				ID:    "Alloc",
				MType: "gauge",
				Delta: 0,
				Value: 527,
			},
			metricSlice: []entity.Metric{{ID: "Alloc", MType: "gauge", Delta: 0, Value: 527}},
			mockBehaviour: func(s *Mockstorage, metric entity.Metric, metricSlice []entity.Metric, key string) {
			},
			want:  want{http.StatusBadRequest, "/update/gauge/Alloc/val"},
			param: map[string]any{"update": "update", "type": "gauge", "name": "Alloc", "value": "val"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			r := httptest.NewRequest(http.MethodPost, tt.want.url, nil)
			rctx := chi.NewRouteContext()
			for k, v := range tt.param {
				strVal := v.(string)
				rctx.URLParams.Add(k, strVal)
			}

			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
			rw := httptest.NewRecorder()

			storageMock := NewMockstorage(c)
			tt.mockBehaviour(storageMock, tt.metric, tt.metricSlice, tt.param["name"].(string))
			handler := NewHandler(storageMock, logger.NewLoggingVar())
			handler.Safe(rw, r)

			res := rw.Result()
			res.Body.Close()
			assert.Equal(t, tt.want.code, res.StatusCode)
		})
	}
}

func TestHandler_Get(t *testing.T) {
	type mockBehaviour func(s *Mockstorage, metric entity.Metric, key string)
	type want struct {
		code int
		url  string
	}
	tests := []struct {
		name          string
		mockBehaviour mockBehaviour
		metric        entity.Metric
		want          want
		param         map[string]any
	}{
		{
			name: "Test#1, correct test with counter metric",
			metric: entity.Metric{
				ID:    "PollCount",
				MType: "counter",
				Delta: 10,
				Value: 0,
			},
			mockBehaviour: func(s *Mockstorage, metric entity.Metric, key string) {
				s.EXPECT().GetMetric(key).Return(metric, nil)
			},
			want:  want{http.StatusOK, "/value/counter/PollCount"},
			param: map[string]any{"value": "value", "type": "counter", "name": "PollCount"},
		},
		{
			name: "Test#2, err to find metric",
			metric: entity.Metric{
				ID:    "PollCount",
				MType: "counter",
				Delta: 10,
				Value: 0,
			},
			mockBehaviour: func(s *Mockstorage, metric entity.Metric, key string) {
				s.EXPECT().GetMetric(key).Return(entity.Metric{}, errors.New("err to find metric"))
			},
			want:  want{http.StatusNotFound, "/value/counter/PollCount"},
			param: map[string]any{"value": "value", "type": "counter", "name": "PollCount"},
		},
		{
			name: "Test#3, unknown type",
			metric: entity.Metric{
				ID:    "PollCount",
				MType: "counter",
				Delta: 10,
				Value: 0,
			},
			mockBehaviour: func(s *Mockstorage, metric entity.Metric, key string) {
				s.EXPECT().GetMetric(key).Return(metric, nil)
			},
			want:  want{http.StatusNotFound, "/value/count/PollCount"},
			param: map[string]any{"value": "value", "type": "count", "name": "PollCount"},
		},
		{
			name: "Test#4, correct test with gauge metric",
			metric: entity.Metric{
				ID:    "Alloc",
				MType: "gauge",
				Delta: 0,
				Value: 213,
			},
			mockBehaviour: func(s *Mockstorage, metric entity.Metric, key string) {
				s.EXPECT().GetMetric(key).Return(metric, nil)
			},
			want:  want{http.StatusOK, "/value/gauge/Alloc"},
			param: map[string]any{"value": "value", "type": "gauge", "name": "Alloc"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			testStorage := NewMockstorage(c)

			r := httptest.NewRequest(http.MethodGet, tt.want.url, nil)
			rctx := chi.NewRouteContext()
			for k, v := range tt.param {
				strVal := v.(string)
				rctx.URLParams.Add(k, strVal)
			}

			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
			rw := httptest.NewRecorder()
			tt.mockBehaviour(testStorage, tt.metric, tt.param["name"].(string))
			handler := NewHandler(testStorage, logger.NewLoggingVar())
			handler.Get(rw, r)

			res := rw.Result()
			res.Body.Close()
			assert.Equal(t, tt.want.code, res.StatusCode)

		})
	}
}

func TestHandler_GetAll(t *testing.T) {
	type mockBehaviour func(s *Mockstorage, metricSlice []entity.Metric)
	type want struct {
		code int
		url  string
	}
	tests := []struct {
		name          string
		metricSlice   []entity.Metric
		mockBehaviour mockBehaviour
		want          want
	}{
		{
			name: "Test#1, correct test ",
			metricSlice: []entity.Metric{
				{ID: "PollCount", MType: "counter", Delta: 10, Value: 0},
				{ID: "Alloc", MType: "gauge", Delta: 0, Value: 527},
			},
			mockBehaviour: func(s *Mockstorage, metricSlice []entity.Metric) {
				s.EXPECT().GetAllMetric().Return(metricSlice, nil)
			},
			want: want{http.StatusOK, "/"},
		},
		{
			name: "Test#2, err to get all metric",
			metricSlice: []entity.Metric{
				{ID: "PollCount", MType: "counter", Delta: 10, Value: 0},
				{ID: "Alloc", MType: "gauge", Delta: 0, Value: 527},
			},
			mockBehaviour: func(s *Mockstorage, metricSlice []entity.Metric) {
				s.EXPECT().GetAllMetric().Return(nil, errors.New("err to get"))
			},
			want: want{http.StatusBadRequest, "/"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			testStorage := NewMockstorage(c)

			r := httptest.NewRequest(http.MethodGet, tt.want.url, nil)
			rctx := chi.NewRouteContext()

			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
			rw := httptest.NewRecorder()
			tt.mockBehaviour(testStorage, tt.metricSlice)
			handler := NewHandler(testStorage, logger.NewLoggingVar())
			handler.GetAll(rw, r)

			res := rw.Result()
			res.Body.Close()
			assert.Equal(t, tt.want.code, res.StatusCode)

		})
	}
}

func TestHandler_SafeBatch(t *testing.T) {
	type mockBehaviour func(s *Mockstorage, metricSlice []entity.Metric)
	type want struct {
		code int
		url  string
	}
	tests := []struct {
		name          string
		mockBehaviour mockBehaviour
		metricSlice   []entity.Metric
		want          want
		param         map[string]any
	}{
		{
			name:        "Test#1, check if all correct ",
			metricSlice: []entity.Metric{{ID: "PollCount", MType: "counter", Delta: 10, Value: 0}},
			mockBehaviour: func(s *Mockstorage, metricSlice []entity.Metric) {
				s.EXPECT().SetMetrics(metricSlice).Return(nil)
			},
			want:  want{http.StatusOK, "/updates"},
			param: map[string]any{"update": "updates"},
		},
		{
			name:        "Test#2, err to safe",
			metricSlice: []entity.Metric{{ID: "PollCount", MType: "counter", Delta: 10, Value: 0}},
			mockBehaviour: func(s *Mockstorage, metricSlice []entity.Metric) {
				s.EXPECT().SetMetrics(metricSlice).Return(errors.New("error to safe"))
			},
			want:  want{http.StatusInternalServerError, "/updates"},
			param: map[string]any{"update": "updates"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			body, err := json.Marshal(tt.metricSlice)
			if err != nil {
				fmt.Println("err to marshal")
			}

			r := httptest.NewRequest(http.MethodPost, tt.want.url,
				bytes.NewBuffer(body))
			rctx := chi.NewRouteContext()
			for k, v := range tt.param {
				strVal := v.(string)
				rctx.URLParams.Add(k, strVal)
			}

			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
			rw := httptest.NewRecorder()
			r.Header.Set("Content-Type", "application/json")

			storageMock := NewMockstorage(c)
			tt.mockBehaviour(storageMock, tt.metricSlice)
			handler := NewHandler(storageMock, logger.NewLoggingVar())
			handler.SafeBatch(rw, r)

			res := rw.Result()
			res.Body.Close()
			assert.Equal(t, tt.want.code, res.StatusCode)
		})
	}
}

func TestHandler_CheckConnection(t *testing.T) {
	type mockBehaviour func(s *Mockstorage, ctx context.Context)
	type want struct {
		code int
		url  string
	}
	tests := []struct {
		name          string
		mockBehaviour mockBehaviour
		want          want
		param         map[string]any
	}{
		{
			name: "Test#1, check if all correct ",
			mockBehaviour: func(s *Mockstorage, ctx context.Context) {
				s.EXPECT().CheckPing(ctx).Return(nil)
				//s.EXPECT().CheckPing(ctx).Return(errors.New("err to check")
			},
			want:  want{http.StatusOK, "/ping"},
			param: map[string]any{"ping": "ping"},
		},
		{
			name: "Test#1, check if all correct ",
			mockBehaviour: func(s *Mockstorage, ctx context.Context) {
				s.EXPECT().CheckPing(ctx).Return(errors.New("err to check"))
			},
			want:  want{http.StatusInternalServerError, "/ping"},
			param: map[string]any{"ping": "ping"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			r := httptest.NewRequest(http.MethodGet, tt.want.url, nil)
			rctx := chi.NewRouteContext()
			for k, v := range tt.param {
				strVal := v.(string)
				rctx.URLParams.Add(k, strVal)
			}

			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
			rw := httptest.NewRecorder()

			ctx, cancel := context.WithCancel(r.Context())
			defer cancel()

			storageMock := NewMockstorage(c)

			tt.mockBehaviour(storageMock, ctx)
			handler := NewHandler(storageMock, logger.NewLoggingVar())
			handler.CheckConnection(rw, r)

			res := rw.Result()
			res.Body.Close()
			assert.Equal(t, tt.want.code, res.StatusCode)
		})
	}
}
