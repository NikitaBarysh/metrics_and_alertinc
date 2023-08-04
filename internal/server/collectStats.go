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

var counter = map[string]int64{
	"PollCount": 64,
}

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
		m.Collect(ctx)
	}()

	go func() {
		defer wg.Done()
		m.SendGauge(ctx)
	}()
	wg.Wait()

	go func() {
		defer wg.Done()
		SendCounter(ctx)
	}()
}

func (m *MemStorageAction) Collect(ctx context.Context) {
	memStats := runtime.MemStats{}
	for {
		if ctx.Err() != nil {
			fmt.Println("exit #1")
			return
		}
		runtime.ReadMemStats(&memStats)
		m.MemStorage.Put("Alloc", float64(memStats.Alloc))
		m.MemStorage.Put("BuckHashSys", float64(memStats.BuckHashSys))
		m.MemStorage.Put("Frees", float64(memStats.Frees))
		m.MemStorage.Put("GCCPUFraction", memStats.GCCPUFraction)
		m.MemStorage.Put("GCSys", float64(memStats.GCSys))
		m.MemStorage.Put("HeapAlloc", float64(memStats.HeapAlloc))
		m.MemStorage.Put("HeapIdle", float64(memStats.HeapIdle))
		m.MemStorage.Put("HeapInuse", float64(memStats.HeapInuse))
		m.MemStorage.Put("HeapObjects", float64(memStats.HeapObjects))
		m.MemStorage.Put("HeapReleased", float64(memStats.HeapReleased))
		m.MemStorage.Put("HeapSys", float64(memStats.HeapSys))
		m.MemStorage.Put("LastGC", float64(memStats.LastGC))
		m.MemStorage.Put("Lookups", float64(memStats.Lookups))
		m.MemStorage.Put("MCacheInuse", float64(memStats.MCacheInuse))
		m.MemStorage.Put("MCacheSys", float64(memStats.MCacheSys))
		m.MemStorage.Put("MSpanInuse", float64(memStats.MSpanInuse))
		m.MemStorage.Put("MSpanSys", float64(memStats.MSpanSys))
		m.MemStorage.Put("Mallocs", float64(memStats.Mallocs))
		m.MemStorage.Put("NextGC", float64(memStats.NextGC))
		m.MemStorage.Put("NumForcedGC", float64(memStats.NumForcedGC))
		m.MemStorage.Put("NumGC", float64(memStats.NumGC))
		m.MemStorage.Put("OtherSys", float64(memStats.OtherSys))
		m.MemStorage.Put("PauseTotalNs", float64(memStats.PauseTotalNs))
		m.MemStorage.Put("StackInuse", float64(memStats.StackInuse))
		m.MemStorage.Put("StackSys", float64(memStats.StackSys))
		m.MemStorage.Put("Sys", float64(memStats.Sys))
		m.MemStorage.Put("TotalAlloc", float64(memStats.TotalAlloc))
		m.MemStorage.Put("RandomValue", rand.Float64())
		time.Sleep(poolInterval)
		counter["PollCounter"]++
	}
}

func (m *MemStorageAction) SendGauge(ctx context.Context) {
	for {
		if ctx.Err() != nil {
			fmt.Println("exit #2")
			return
		}
		for _, metricName := range m.MemStorage.Get() {
			if metricValue, ok := m.MemStorage.Read(metricName); ok {
				fmt.Println(metricName, metricValue)
				url := "http://localhost:8080/update/gauge" + metricName + fmt.Sprintf("%f", metricValue)
				request, err := http.NewRequest(http.MethodPost, url, nil)
				if err != nil {
					panic(err)
				}
				request.Header.Set(`Content-Type`, "text/plain")
			}
		}
		time.Sleep(reportInterval)
	}
}

func SendCounter(ctx context.Context) {
	for {
		if ctx.Err() != nil {
			fmt.Println("exit #2")
			return
		}
		for metricName, metricValue := range counter {
			fmt.Println(metricName, metricValue)
			url := "http://localhost:8080/update/counter/" + metricName + fmt.Sprintf("%d", metricValue)
			request, err := http.NewRequest(http.MethodPost, url, nil)
			if err != nil {
				panic(err)
			}
			request.Header.Set(`Content-Type`, "text/plain")
		}
		time.Sleep(reportInterval)
	}
}
