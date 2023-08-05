package server

//import (
//	"github.com/stretchr/testify/assert"
//	"net/http"
//	"net/http/httptest"
//	"testing"
//)
//
//func Test_SendGauge(t *testing.T) {
//	type want struct {
//		code int
//		url  string
//	}
//
//	tests := []struct {
//		name string
//		want want
//	}{
//		{
//			name: "Test#1, not enough argument in url",
//			want: want{
//				code: http.StatusOk,
//				url:  "http://localhost:8080/update/counter/someMetric",
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			r := httptest.NewRequest(http.MethodPost, tt.want.url, nil)
//			rw := httptest.NewRecorder()
//
//			assert.Equal(t, res.StatusCode, tt.want.code)
//		})
//	}
//}
