package router

//
//import (
//	"github.com/NikitaBarysh/metrics_and_alertinc/internal/handlers"
//	"github.com/go-chi/chi/v5"
//	"reflect"
//	"testing"
//)
//
//func TestRouter_Register(t *testing.T) {
//	type fields struct {
//		metricHandler *handlers.Handler
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		want   *chi.Mux
//	}{
//		{
//			name: "ok",
//			fields: fields{metricHandler: },
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			rt := &Router{
//				metricHandler: tt.fields.metricHandler,
//			}
//			if got := rt.Register(); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("Register() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
