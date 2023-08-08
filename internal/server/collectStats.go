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

	//go func() {
	//	defer wg.Done()
	//	m.SendGauge(ctx)
	//}()
	//
	//go func() {
	//	defer wg.Done()
	//	SendCounter(ctx)
	//}()

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
		m.MemStorage.PutMetric("Alloc", float64(memStats.Alloc))
		m.MemStorage.PutMetric("BuckHashSys", float64(memStats.BuckHashSys))
		m.MemStorage.PutMetric("Frees", float64(memStats.Frees))
		m.MemStorage.PutMetric("GCCPUFraction", memStats.GCCPUFraction)
		m.MemStorage.PutMetric("GCSys", float64(memStats.GCSys))
		m.MemStorage.PutMetric("HeapAlloc", float64(memStats.HeapAlloc))
		m.MemStorage.PutMetric("HeapIdle", float64(memStats.HeapIdle))
		m.MemStorage.PutMetric("HeapInuse", float64(memStats.HeapInuse))
		m.MemStorage.PutMetric("HeapObjects", float64(memStats.HeapObjects))
		m.MemStorage.PutMetric("HeapReleased", float64(memStats.HeapReleased))
		m.MemStorage.PutMetric("HeapSys", float64(memStats.HeapSys))
		m.MemStorage.PutMetric("LastGC", float64(memStats.LastGC))
		m.MemStorage.PutMetric("Lookups", float64(memStats.Lookups))
		m.MemStorage.PutMetric("MCacheInuse", float64(memStats.MCacheInuse))
		m.MemStorage.PutMetric("MCacheSys", float64(memStats.MCacheSys))
		m.MemStorage.PutMetric("MSpanInuse", float64(memStats.MSpanInuse))
		m.MemStorage.PutMetric("MSpanSys", float64(memStats.MSpanSys))
		m.MemStorage.PutMetric("Mallocs", float64(memStats.Mallocs))
		m.MemStorage.PutMetric("NextGC", float64(memStats.NextGC))
		m.MemStorage.PutMetric("NumForcedGC", float64(memStats.NumForcedGC))
		m.MemStorage.PutMetric("NumGC", float64(memStats.NumGC))
		m.MemStorage.PutMetric("OtherSys", float64(memStats.OtherSys))
		m.MemStorage.PutMetric("PauseTotalNs", float64(memStats.PauseTotalNs))
		m.MemStorage.PutMetric("StackInuse", float64(memStats.StackInuse))
		m.MemStorage.PutMetric("StackSys", float64(memStats.StackSys))
		m.MemStorage.PutMetric("Sys", float64(memStats.Sys))
		m.MemStorage.PutMetric("TotalAlloc", float64(memStats.TotalAlloc))
		m.MemStorage.PutMetric("RandomValue", rand.Float64())
		time.Sleep(poolInterval)
		m.MemStorage.PutMetric("PollCount", int64(+1))
	}
}

//func (m *MemStorageAction) SendGauge(ctx context.Context) {
//	for {
//		for _, metricName := range m.MemStorage.GetMetric() {
//			if metricValue, ok := m.MemStorage.ReadMetric(metricName); ok {
//				fmt.Println(metricValue, metricName)
//				url := "http://localhost:8080/update/gauge" + metricName + fmt.Sprintf("%f", metricValue)
//				request, err := http.NewRequest(http.MethodPost, url, nil)
//				if err != nil {
//					panic(err)
//				}
//				request.Header.Set(`Content-Type`, "text/plain")
//			}
//		}
//		time.Sleep(reportInterval)
//	}
//}
//
//func SendCounter(ctx context.Context) {
//	for {
//		for metricName, metricValue := range counter {
//			fmt.Println(metricValue, metricName)
//			url := "http://localhost:8080/update/counter/" + metricName + fmt.Sprintf("%d", metricValue)
//			request, err := http.NewRequest(http.MethodPost, url, nil)
//			if err != nil {
//				panic(err)
//			}
//			request.Header.Set(`Content-Type`, "text/plain")
//		}
//		time.Sleep(reportInterval)
//	}
//}

func (m *MemStorageAction) SendMetric(ctx context.Context) {
	for {
		if ctx.Err() != nil {
			fmt.Println("err in func SendMetric")
			return
		}
		for _, metricName := range m.MemStorage.GetMetric() {
			if metricValue, ok := m.MemStorage.ReadMetric(metricName); ok {
				if metricName == "PollCount" {
					url := "http://localhost:8080/update/counter/" + metricName + fmt.Sprintf("%d", metricValue)
					request, err := http.NewRequest(http.MethodPost, url, nil)
					if err != nil {
						panic(err)
					}
					request.Header.Set(`Content-Type`, "text/plain")
				} else {
					url := "http://localhost:8080/update/gauge" + metricName + fmt.Sprintf("%f", metricValue)
					request, err := http.NewRequest(http.MethodPost, url, nil)
					if err != nil {
						panic(err)
					}
					request.Header.Set(`Content-Type`, "text/plain")
				}
			}
		}
		time.Sleep(reportInterval)
	}
}
