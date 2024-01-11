package filestorage

import (
	"fmt"
	"testing"

	"github.com/NikitaBarysh/metrics_and_alertinc/internal/entity"
)

func BenchmarkNewFileEngine(b *testing.B) {
	file, err := NewFileEngine("data.json")
	if err != nil {
		fmt.Println("err to create")
	}

	metric := map[string]entity.Metric{"Alloc": {ID: "Alloc", MType: "gauge", Delta: 0, Value: 527}}

	b.ResetTimer()

	b.Run("set metrics", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			err := file.SetMetrics(metric)
			if err != nil {
				fmt.Println("err to set metric")
			}
		}
	})

	b.Run("get metric", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			file.GetAllMetric()
		}
	})
}
