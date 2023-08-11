package handlers

import (
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/storage/repositories"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

//func testRequest(t *testing.T, ts *httptest.Server, method, path string) *http.Response {
//	req, err := http.NewRequest(method, ts.URL+path, nil)
//	require.NoError(t, err)
//
//	resp, err := ts.Client().Do(req)
//	require.NoError(t, err)
//	defer resp.Body.Close()
//
//	return resp
//}

func Test_Safe(t *testing.T) {
	type want struct {
		code int
		url  string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			"Test#1, not enough argument in url", want{http.StatusNotFound, "/update/counter/someMetric"},
		},
		{
			"Test#2, check if all correct", want{http.StatusOK, "/update/counter/Alloc/527"},
		},
		{
			"Test#3, wrong metric", want{http.StatusNotImplemented, "/update/anytype/someMetric/527"},
		},
		{
			"Test#4, wrong value for counter metric", want{http.StatusBadRequest, "/update/counter/someMetric/wrong"},
		},
		{
			"Test#5, wrong value for gauge metric", want{http.StatusBadRequest, "/update/gauge/someMetric/wrong"},
		},
		{
			"Test#6, not update in url", want{http.StatusNotFound, "/notupdate/gauge/someMetric/527"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, tt.want.url, nil)
			rw := httptest.NewRecorder()
			handler := NewHandler(repositories.NewMemStorage())
			handler.Safe(rw, r)

			res := rw.Result()
			res.Body.Close()
			assert.Equal(t, res.StatusCode, tt.want.code)

		})
	}
}

//func testRequest(t *testing.T, ts *httptest.Server, method, path string) *http.Response {
//	req, err := http.NewRequest(method, ts.URL+path, nil)
//	require.NoError(t, err)
//
//	resp, err := ts.Client().Do(req)
//	require.NoError(t, err)
//	defer resp.Body.Close()
//
//	return resp
//}
//func Test_Safe(t *testing.T) {
//	storage := repositories.NewMemStorage()
//	handler := NewHandler(storage)
//	r := router2.NewRouter(handler)
//	ts := httptest.NewServer(handler)
//
//	defer ts.Close()
//
//	type want struct {
//		code int
//		url  string
//	}
//	tests := []struct {
//		name string
//		want want
//	}{
//		{
//			name: "Test#1, not enough argument in url",
//			want: want{
//				code: http.StatusNotFound,
//				url:  "/update/counter/someMetric",
//			},
//		},
//		{
//			name: "Test#2, check if all correct",
//			want: want{
//				code: http.StatusOK,
//				url:  "/update/counter/someMetric/527",
//			},
//		},
//		{
//			name: "Test#3, wrong metric",
//			want: want{
//				code: http.StatusNotImplemented,
//				url:  "/update/anytype/someMetric/527",
//			},
//		},
//		{
//			name: "Test#4, wrong value for counter metric",
//			want: want{
//				code: http.StatusBadRequest,
//				url:  "/update/counter/someMetric/wrong",
//			},
//		},
//		{
//			name: "Test#5, wrong value for gauge metric",
//			want: want{
//				code: http.StatusBadRequest,
//				url:  "/update/gauge/someMetric/wrong",
//			},
//		},
//		{
//			name: "Test#6, not update in url",
//			want: want{
//				code: http.StatusNotFound,
//				url:  "/notupdate/gauge/someMetric/527",
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			res := testRequest(t, ts, "POST", tt.want.url)
//			assert.Equal(t, tt.want.code, res.StatusCode, "Статус кода не соответсвует с ожидаемым")
//			res.Body.Close()
//		})
//	}
//}

//func Test_get(t *testing.T) {
//	storage := repositories.NewMemStorage()
//	handler := NewHandler(storage)
//	r := router2.NewRouter(handler)
//	ts := httptest.NewServer(r)
//	defer ts.Close()
//
//	type want struct {
//		code int
//		url  string
//	}
//	tests := []struct {
//		name string
//		want want
//	}{
//		{
//			name: "Test#1, not enough argument in url",
//			want: want{
//				code: http.StatusBadRequest,
//				url:  "/value",
//			},
//		},
//		{
//			name: "Test#2, check if all correct",
//			want: want{
//				code: http.StatusBadRequest,
//				url:  "/value/gauge/Alloc",
//			},
//		},
//		{
//			name: "Test#3, wrong metric",
//			want: want{
//				code: http.StatusNotImplemented,
//				url:  "/update/anytype/someMetric/527",
//			},
//		},
//		{
//			name: "Test#4, not value in url",
//			want: want{
//				code: http.StatusNotFound,
//				url:  "/notvalue/gauge/Alloc",
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			res := testRequest(t, ts, "GET", tt.want.url)
//			assert.Equal(t, tt.want.code, res.StatusCode, "Статус кода не соответсвует с ожидаемым")
//			res.Body.Close()
//		})
//	}
//}
