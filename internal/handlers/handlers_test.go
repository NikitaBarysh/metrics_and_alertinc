package handlers

//import (
//	"github.com/NikitaBarysh/metrics_and_alertinc/internal/router"
//	"github.com/NikitaBarysh/metrics_and_alertinc/internal/storage/repositories"
//	"github.com/go-chi/chi/v5"
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/require"
//	"net/http"
//	"net/http/httptest"
//	"testing"
//)
//
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
//
//func Test_Safe(t *testing.T) {
//	storage := repositories.NewMemStorage()
//	handler := NewHandler(storage)
//	router := router.NewRouter(handler)
//	chiRouter := chi.NewRouter()
//	chiRouter.Mount("/")
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
//			"Test#1, not enough argument in url", want{http.StatusNotFound, "/update/counter/someMetric"},
//		},
//		{
//			"Test#2, check if all correct", want{http.StatusOK, "/update/counter/Alloc/527"},
//		},
//		{
//			"Test#3, wrong metric", want{http.StatusNotImplemented, "/update/anytype/someMetric/527"},
//		},
//		{
//			"Test#4, wrong value for counter metric", want{http.StatusBadRequest, "/update/counter/someMetric/wrong"},
//		},
//		{
//			"Test#5, wrong value for gauge metric", want{http.StatusBadRequest, "/update/gauge/someMetric/wrong"},
//		},
//		{
//			"Test#6, not update in url", want{http.StatusNotFound, "/notupdate/gauge/someMetric/527"},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			r := httptest.NewRequest(http.MethodPost, tt.want.url, nil)
//			rw := httptest.NewRecorder()
//			handler := NewHandler(repositories.NewMemStorage())
//			handler.Safe(rw, r)
//
//			res := rw.Result()
//			res.Body.Close()
//			assert.Equal(t, tt.want.code, res.StatusCode)
//
//		})
//	}
//}
