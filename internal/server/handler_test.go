package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_post(t *testing.T) {
	type want struct {
		code int
		url  string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "Test#1, not enough argument in url",
			want: want{
				code: http.StatusNotFound,
				url:  "http://localhost:8080/update/counter/someMetric",
			},
		},
		{
			name: "Test#2, check if all correct",
			want: want{
				code: http.StatusOK,
				url:  "http://localhost:8080/update/counter/someMetric/527",
			},
		}, {
			name: "Test#3, wrong metric",
			want: want{
				code: http.StatusNotImplemented,
				url:  "http://localhost:8080/update/anytype/someMetric/527",
			},
		},
		{
			name: "Test#4, wrong value for counter metric",
			want: want{
				code: http.StatusBadRequest,
				url:  "http://localhost:8080/update/counter/someMetric/wrong",
			},
		},
		{
			name: "Test#5, wrong value for gauge metric",
			want: want{
				code: http.StatusBadRequest,
				url:  "http://localhost:8080/update/gauge/someMetric/wrong",
			},
		},
		{
			name: "Test#6, not update in url",
			want: want{
				code: http.StatusNotFound,
				url:  "http://localhost:8080/notupdate/gauge/someMetric/527",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, tt.want.url, nil)
			rw := httptest.NewRecorder()
			post(rw, r)

			res := rw.Result()
			assert.Equal(t, res.StatusCode, tt.want.code)

		})
	}
}
