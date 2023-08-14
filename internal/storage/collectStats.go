package storage

import (
	"context"
	"fmt"
	"math/rand"
	"runtime"
	"time"
)

const (
	poolInterval   = time.Second * 2
	reportInterval = time.Second * 10
)

type sender interface {
	SendPost(url string)
}

func (m *MetricAction) Run(ctx context.Context) {

	collectTicker := time.NewTicker(poolInterval)
	defer collectTicker.Stop()

	sendTicker := time.NewTicker(reportInterval)
	defer sendTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-collectTicker.C:
			m.CollectMetric()
		case <-sendTicker.C:
			m.SendMetric()
		}
	}
}

func (m *MetricAction) CollectMetric() {
	memStats := runtime.MemStats{}
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
}

func (m *MetricAction) SendMetric() {
	for metricName, metricValue := range m.MemStorage.ReadGaugeMetric() {
		url := fmt.Sprintf("http://localhost:8080/update/gauge/%s/%2f", metricName, metricValue)
		m.sender.SendPost(url)
	}
	for metricName, metricValue := range m.MemStorage.ReadCounterMetric() {
		url := fmt.Sprintf("http://localhost:8080/update/counter/%s/%d", metricName, metricValue)
		m.sender.SendPost(url)
	}
}
