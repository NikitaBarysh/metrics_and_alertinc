package storage

import (
	"context"
	"fmt"
	"math/rand"
	"runtime"
	"time"

	"github.com/NikitaBarysh/metrics_and_alertinc/internal/storage/repositories"
)

type sender interface {
	SendPost(ctx context.Context, url string, storage repositories.MemStorageStruct)
}

func (m *MetricAction) Run(ctx context.Context, pollInterval int64, reportInterval int64, flagRunAddr string) error {

	collectTicker := time.NewTicker(time.Second * time.Duration(pollInterval))
	defer collectTicker.Stop()

	sendTicker := time.NewTicker(time.Second * time.Duration(reportInterval))
	defer sendTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-collectTicker.C:
			m.CollectMetric()
		case <-sendTicker.C:
			m.SendMetric(ctx, flagRunAddr)
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

func (m *MetricAction) SendMetric(ctx context.Context, flagRunAddr string) error {
	for metricName, metricValue := range m.MemStorage.ReadMetric() {
		metricType := m.MemStorage.MemStorageMap[metricName].MType
		switch metricType {
		case "gauge":
			url := fmt.Sprintf("http://%s/update/%s/%s/%.2f", flagRunAddr, metricType, metricName, metricValue.Value)
			m.sender.SendPost(ctx, url, metricValue)
		case "counter":
			url := fmt.Sprintf("http://%s/update/%s/%s/%d", flagRunAddr, metricType, metricName, metricValue.Delta)
			m.sender.SendPost(ctx, url, metricValue)
		}
	}
	return nil
}
