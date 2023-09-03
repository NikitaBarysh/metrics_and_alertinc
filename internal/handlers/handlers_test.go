package handlers

import (
	"context"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/logger"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NikitaBarysh/metrics_and_alertinc/internal/storage/repositories"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestHandler_Safe(t *testing.T) {
	type want struct {
		code int
		url  string
	}
	tests := []struct {
		name  string
		want  want
		param map[string]any
	}{
		{
			"Test#1, not enough argument in url",
			want{http.StatusBadRequest, "/update/counter/someMetric"},
			map[string]any{"update": "update", "type": "counter", "name": "someMetric"},
		},
		{"Test#2, check if all correct",
			want{http.StatusOK, "/update/gauge/Alloc/527"},
			map[string]any{"update": "update", "type": "gauge", "name": "Alloc", "value": "527"},
		},
		{"Test#3, wrong metric",
			want{http.StatusNotImplemented, "/update/anytype/someMetric/527"},
			map[string]any{"update": "update", "type": "anytype", "name": "someMetric", "value": "527"},
		},
		{"Test#4, wrong value for counter metric",
			want{http.StatusBadRequest, "/update/counter/someMetric/wrong"},
			map[string]any{"update": "update", "type": "counter", "name": "someMetric", "value": "wrong"},
		},
		{"Test#5, wrong value for gauge metric",
			want{http.StatusBadRequest, "/update/gauge/someMetric/wrong"},
			map[string]any{"update": "update", "type": "gauge", "name": "someMetric", "value": "wrong"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, tt.want.url, nil)
			rctx := chi.NewRouteContext()
			for k, v := range tt.param {
				strVal := v.(string)
				rctx.URLParams.Add(k, strVal)
			}

			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
			rw := httptest.NewRecorder()
			handler := NewHandler(repositories.NewMemStorage(), logger.NewLoggingVar())
			handler.Safe(rw, r)

			res := rw.Result()
			res.Body.Close()
			assert.Equal(t, tt.want.code, res.StatusCode)

		})
	}
}

func TestHandler_GetGaugeMetric(t *testing.T) {
	type want struct {
		code int
		url  string
	}
	tests := []struct {
		name    string
		want    want
		param   map[string]any
		metrics map[string]float64
	}{
		{
			"Test#1, correct test with gauge metric",
			want{http.StatusOK, "/value/gauge/Alloc"},
			map[string]any{"value": "value", "type": "gauge", "name": "Alloc"},
			map[string]float64{"Alloc": 456},
		},
		{
			"Test#2, unknown metric type",
			want{http.StatusNotFound, "/value/any/PollCounter"},
			map[string]any{"value": "value", "type": "any", "name": "Alloc"},
			map[string]float64{"PollCounter": 456},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testStorage := repositories.NewMemStorage()
			for k, v := range tt.metrics {
				testStorage.UpdateGaugeMetric(k, v)
			}
			r := httptest.NewRequest(http.MethodGet, tt.want.url, nil)
			rctx := chi.NewRouteContext()
			for k, v := range tt.param {
				strVal := v.(string)
				rctx.URLParams.Add(k, strVal)
			}

			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
			rw := httptest.NewRecorder()
			handler := NewHandler(testStorage, logger.NewLoggingVar())
			handler.Get(rw, r)

			res := rw.Result()
			res.Body.Close()
			assert.Equal(t, tt.want.code, res.StatusCode)

		})
	}
}

func TestHandler_GetCounterMetric(t *testing.T) {
	type want struct {
		code int
		url  string
	}
	tests := []struct {
		name    string
		want    want
		param   map[string]any
		metrics map[string]int64
	}{
		{
			"Test#1, correct test with gauge metric",
			want{http.StatusOK, "/value/counter/PollCount"},
			map[string]any{"value": "value", "type": "counter", "name": "PollCount"},
			map[string]int64{"PollCount": 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testStorage := repositories.NewMemStorage()
			for k, v := range tt.metrics {
				testStorage.UpdateCounterMetric(k, v)
			}
			r := httptest.NewRequest(http.MethodGet, tt.want.url, nil)
			rctx := chi.NewRouteContext()
			for k, v := range tt.param {
				strVal := v.(string)
				rctx.URLParams.Add(k, strVal)
			}

			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
			rw := httptest.NewRecorder()
			handler := NewHandler(testStorage, logger.NewLoggingVar())
			handler.Get(rw, r)

			res := rw.Result()
			res.Body.Close()
			assert.Equal(t, tt.want.code, res.StatusCode)

		})
	}
}

func TestHandler_GetAll(t *testing.T) {
	type want struct {
		code int
		url  string
	}
	tests := []struct {
		name    string
		want    want
		param   map[string]any
		metrics map[string]float64
	}{
		{
			"Test#1, correct test ",
			want{http.StatusOK, "/value/gauge/Alloc"},
			map[string]any{"value": "value", "type": "gauge", "name": "Alloc"},
			map[string]float64{
				"Alloc":       456.34,
				"Frees":       354.35,
				"BuckHashSys": 645.64,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testStorage := repositories.NewMemStorage()
			for k, v := range tt.metrics {
				testStorage.UpdateGaugeMetric(k, v)
			}
			listMetric := testStorage.GetAllMetric()
			r := httptest.NewRequest(http.MethodGet, tt.want.url, nil)
			rctx := chi.NewRouteContext()
			for _, v := range listMetric {
				rctx.URLParams.Add("AllMetric", v)
			}

			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
			rw := httptest.NewRecorder()
			handler := NewHandler(testStorage, logger.NewLoggingVar())
			handler.GetAll(rw, r)

			res := rw.Result()
			res.Body.Close()
			assert.Equal(t, tt.want.code, res.StatusCode)

		})
	}
}
