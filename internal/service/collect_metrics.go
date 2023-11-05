package service

import (
	"math/rand"
	"runtime"
)

func (m *MetricAction) CollectMetric() {
	memStats := runtime.MemStats{}
	runtime.ReadMemStats(&memStats)
	m.storage.UpdateGaugeMetric("Alloc", float64(memStats.Alloc))
	m.storage.UpdateGaugeMetric("BuckHashSys", float64(memStats.BuckHashSys))
	m.storage.UpdateGaugeMetric("Frees", float64(memStats.Frees))
	m.storage.UpdateGaugeMetric("GCCPUFraction", memStats.GCCPUFraction)
	m.storage.UpdateGaugeMetric("GCSys", float64(memStats.GCSys))
	m.storage.UpdateGaugeMetric("HeapAlloc", float64(memStats.HeapAlloc))
	m.storage.UpdateGaugeMetric("HeapIdle", float64(memStats.HeapIdle))
	m.storage.UpdateGaugeMetric("HeapInuse", float64(memStats.HeapInuse))
	m.storage.UpdateGaugeMetric("HeapObjects", float64(memStats.HeapObjects))
	m.storage.UpdateGaugeMetric("HeapReleased", float64(memStats.HeapReleased))
	m.storage.UpdateGaugeMetric("HeapSys", float64(memStats.HeapSys))
	m.storage.UpdateGaugeMetric("LastGC", float64(memStats.LastGC))
	m.storage.UpdateGaugeMetric("Lookups", float64(memStats.Lookups))
	m.storage.UpdateGaugeMetric("MCacheInuse", float64(memStats.MCacheInuse))
	m.storage.UpdateGaugeMetric("MCacheSys", float64(memStats.MCacheSys))
	m.storage.UpdateGaugeMetric("MSpanInuse", float64(memStats.MSpanInuse))
	m.storage.UpdateGaugeMetric("MSpanSys", float64(memStats.MSpanSys))
	m.storage.UpdateGaugeMetric("Mallocs", float64(memStats.Mallocs))
	m.storage.UpdateGaugeMetric("NextGC", float64(memStats.NextGC))
	m.storage.UpdateGaugeMetric("NumForcedGC", float64(memStats.NumForcedGC))
	m.storage.UpdateGaugeMetric("NumGC", float64(memStats.NumGC))
	m.storage.UpdateGaugeMetric("OtherSys", float64(memStats.OtherSys))
	m.storage.UpdateGaugeMetric("PauseTotalNs", float64(memStats.PauseTotalNs))
	m.storage.UpdateGaugeMetric("StackInuse", float64(memStats.StackInuse))
	m.storage.UpdateGaugeMetric("StackSys", float64(memStats.StackSys))
	m.storage.UpdateGaugeMetric("Sys", float64(memStats.Sys))
	m.storage.UpdateGaugeMetric("TotalAlloc", float64(memStats.TotalAlloc))
	m.storage.UpdateGaugeMetric("RandomValue", rand.Float64())
	m.storage.UpdateCounterMetric("PollCount", int64(1))
}
