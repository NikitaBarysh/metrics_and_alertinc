package memstorage

import (
	"context"
	"log"
	"testing"

	"github.com/NikitaBarysh/metrics_and_alertinc/config/server"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/entity"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/repository/filestorage"
)

func BenchmarkMemStorage(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var file *filestorage.FileEngine

	env, envErr := server.NewServer()
	if envErr != nil {
		log.Fatalf("config err: %s\n", envErr)
	}

	cfg, configError := server.NewConfig(env)
	if configError != nil {
		log.Fatalf("config err: %s\n", configError)
	}

	metricSliceGauge := []entity.Metric{{ID: "Alloc", MType: "gauge", Delta: 0, Value: 527}}
	metricSliceCounter := []entity.Metric{{ID: "PollCount", MType: "counter", Delta: 12, Value: 0}}

	storage, err := NewMemStorage(ctx, cfg, file)
	if err != nil {
		log.Println("err to create")
	}
	b.ResetTimer()

	b.Run("update gauge metric", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			storage.UpdateGaugeMetric("Alloc", 123.23)
		}
	})

	b.Run("update counter metric", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			storage.UpdateCounterMetric("PollCount", 1)
		}
	})

	b.Run("get metric", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			storage.GetMetric("Alloc")
		}
	})

	b.Run("get all metric", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			storage.GetAllMetric()
		}
	})

	b.Run("set gauge metric", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			storage.SetMetrics(metricSliceGauge)
		}
	})

	b.Run("set counter metric", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			storage.SetMetrics(metricSliceCounter)
		}
	})
}
