package server

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/NikitaBarysh/metrics_and_alertinc/internal/storage"
)

const (
	poolInterval   = time.Second * 2
	reportInterval = time.Second * 10
)

type MemStorageAction struct {
	MemStorage *storage.MemStorage
}

func (m *MemStorageAction) Run(ctx context.Context) {
	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		m.CollectMetric(ctx)
	}()

	go func() {
		defer wg.Done()
		m.SendMetric(ctx)
	}()

	wg.Wait()
}

func (m *MemStorageAction) CollectMetric(ctx context.Context) {
	memStats := runtime.MemStats{}
	for {
		if ctx.Err() != nil {
			fmt.Println("err in func CollectMetric")
			return
		}
		runtime.ReadMemStats(&memStats)
		m.MemStorage.UpdateGaugeMetric("Alloc", float64(memStats.Alloc))
		m.MemStorage.UpdateGaugeMetric("BuckHashSys", float64(memStats.BuckHashSys))
		m.MemStorage.UpdateGaugeMetric("Frees", float64(memStats.Frees))
		m.MemStorage.UpdateGaugeMetric("GCCPUFraction", memStats.GCCPUFraction)
		m.MemStorage.UpdateGaugeMetric("GCSys", float64(memStats.GCSys))
		m.MemStorage.UpdateGaugeMetric("HeapAlloc", float64(memStats.HeapAlloc))
		m.MemStorage.UpdateGaugeMetric("HeapIdle", float64(memStats.HeapIdle))
		m.MemStorage.UpdateGaugeMetric("HeapInuse", float64(memStats.HeapInuse))
		m.MemStorage.UpdateGaugeMetric("HeapObjects", float64(memStats.HeapObjects))
		m.MemStorage.UpdateGaugeMetric("HeapReleased", float64(memStats.HeapReleased))
		m.MemStorage.UpdateGaugeMetric("HeapSys", float64(memStats.HeapSys))
		m.MemStorage.UpdateGaugeMetric("LastGC", float64(memStats.LastGC))
		m.MemStorage.UpdateGaugeMetric("Lookups", float64(memStats.Lookups))
		m.MemStorage.UpdateGaugeMetric("MCacheInuse", float64(memStats.MCacheInuse))
		m.MemStorage.UpdateGaugeMetric("MCacheSys", float64(memStats.MCacheSys))
		m.MemStorage.UpdateGaugeMetric("MSpanInuse", float64(memStats.MSpanInuse))
		m.MemStorage.UpdateGaugeMetric("MSpanSys", float64(memStats.MSpanSys))
		m.MemStorage.UpdateGaugeMetric("Mallocs", float64(memStats.Mallocs))
		m.MemStorage.UpdateGaugeMetric("NextGC", float64(memStats.NextGC))
		m.MemStorage.UpdateGaugeMetric("NumForcedGC", float64(memStats.NumForcedGC))
		m.MemStorage.UpdateGaugeMetric("NumGC", float64(memStats.NumGC))
		m.MemStorage.UpdateGaugeMetric("OtherSys", float64(memStats.OtherSys))
		m.MemStorage.UpdateGaugeMetric("PauseTotalNs", float64(memStats.PauseTotalNs))
		m.MemStorage.UpdateGaugeMetric("StackInuse", float64(memStats.StackInuse))
		m.MemStorage.UpdateGaugeMetric("StackSys", float64(memStats.StackSys))
		m.MemStorage.UpdateGaugeMetric("Sys", float64(memStats.Sys))
		m.MemStorage.UpdateGaugeMetric("TotalAlloc", float64(memStats.TotalAlloc))
		m.MemStorage.UpdateGaugeMetric("RandomValue", rand.Float64())
		m.MemStorage.UpdateCounterMetric("PollCount", int64(1))
		time.Sleep(poolInterval)
	}
}

func (m *MemStorageAction) SendMetric(ctx context.Context) {
	for {
		if ctx.Err() != nil {
			fmt.Println("err in func SendMetric")
			return
		}
		for metricName, metricValue := range m.MemStorage.ReadGaugeMetric() {
			url := "http://localhost:8088/update/gauge/" + metricName + "/" + fmt.Sprintf("%f", metricValue)
			//url := fmt.Sprintf("")
			request, err := http.NewRequest(http.MethodPost, url, nil)
			if err != nil {
				panic(err)
			}
			request.Header.Set(`Content-Type`, "text/plain")
			client := &http.Client{}
			_, err = client.Do(request)
			if err != nil {
				fmt.Println(err)
			}
		}
		for metricName, metricValue := range m.MemStorage.ReadCounterMetric() {
			url := "http://localhost:8088/update/counter/" + metricName + "/" + fmt.Sprintf("%d", metricValue)
			//url := fmt.Sprintf("")
			request, err := http.NewRequest(http.MethodPost, url, nil)
			if err != nil {
				panic(err)
			}
			request.Header.Set(`Content-Type`, "text/plain")
			client := &http.Client{}
			client.Do(request)
		}
		//for _, metricName := range m.MemStorage.GetMetric() {
		//	if metricValue, ok := m.MemStorage.ReadMetric(metricName); ok {
		//		url := "http://localhost:8080/update/gauge" + metricName + fmt.Sprintf("%f", metricValue)
		//		if metricName == "PollCount" {
		//			url = "http://localhost:8080/update/counter/" + metricName + fmt.Sprintf("%d", metricValue)
		//		}
		//		request, err := http.Post(url, "text/plain; charset=UTF-8", nil)
		//		if err != nil {
		//			panic(err)
		//		}
		//		request.Header.Set(`Content-Type`, "text/plain")
		//	}
		//}
		time.Sleep(reportInterval)
	}
}
