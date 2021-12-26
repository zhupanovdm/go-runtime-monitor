package service

import (
	"math/rand"
	"runtime"

	"github.com/zhupanovdm/go-runtime-monitor/internal/model/metric"
)

func pollRuntimeMetrics(subscriber chan<- *metric.Metric, PollCount int64) {
	stats := runtime.MemStats{}
	runtime.ReadMemStats(&stats)

	subscriber <- metric.NewGaugeFromUInt("Alloc", stats.Alloc)
	subscriber <- metric.NewGaugeFromUInt("BuckHashSys", stats.BuckHashSys)
	subscriber <- metric.NewGauge("GCCPUFraction", stats.GCCPUFraction)
	subscriber <- metric.NewGaugeFromUInt("GCSys", stats.GCSys)
	subscriber <- metric.NewGaugeFromUInt("HeapAlloc", stats.HeapAlloc)
	subscriber <- metric.NewGaugeFromUInt("HeapIdle", stats.HeapIdle)
	subscriber <- metric.NewGaugeFromUInt("HeapInuse", stats.HeapInuse)
	subscriber <- metric.NewGaugeFromUInt("HeapObjects", stats.HeapObjects)
	subscriber <- metric.NewGaugeFromUInt("HeapReleased", stats.HeapReleased)
	subscriber <- metric.NewGaugeFromUInt("HeapSys", stats.HeapSys)
	subscriber <- metric.NewGaugeFromUInt("LastGC", stats.LastGC)
	subscriber <- metric.NewGaugeFromUInt("Lookups", stats.Lookups)
	subscriber <- metric.NewGaugeFromUInt("MCacheInuse", stats.MCacheInuse)
	subscriber <- metric.NewGaugeFromUInt("MCacheSys", stats.MCacheSys)
	subscriber <- metric.NewGaugeFromUInt("MSpanInuse", stats.MSpanInuse)
	subscriber <- metric.NewGaugeFromUInt("MSpanSys", stats.MSpanSys)
	subscriber <- metric.NewGaugeFromUInt("Mallocs", stats.Mallocs)
	subscriber <- metric.NewGaugeFromUInt("NextGC", stats.NextGC)
	subscriber <- metric.NewGaugeFromUInt("NumForcedGC", uint64(stats.NumForcedGC))
	subscriber <- metric.NewGaugeFromUInt("NumGC", uint64(stats.NumGC))
	subscriber <- metric.NewGaugeFromUInt("OtherSys", stats.OtherSys)
	subscriber <- metric.NewGaugeFromUInt("PauseTotalNs", stats.PauseTotalNs)
	subscriber <- metric.NewGaugeFromUInt("StackInuse", stats.StackInuse)
	subscriber <- metric.NewGaugeFromUInt("StackSys", stats.StackSys)
	subscriber <- metric.NewGaugeFromUInt("Sys", stats.Sys)
	subscriber <- metric.NewCounter("PollCount", PollCount)
	subscriber <- metric.NewGauge("RandomValue", rand.Float64())
}
