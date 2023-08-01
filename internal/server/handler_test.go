package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_post(t *testing.T) {
	type want struct {
		code        int
		contentType string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "positive test #1",
			want: want{
				code:        200,
				contentType: "text/plain",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "http://localhost:8080/", nil)
			rw := httptest.NewRecorder()
			post(rw, r)

			res := rw.Result()
			assert.Equal(t, res.StatusCode, tt.want.code)
			assert.Equal(t, res.Header.Get("Content-Type"), tt.want.contentType)

		})
	}
}
