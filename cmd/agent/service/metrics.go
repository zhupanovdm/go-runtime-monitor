package service

import (
	"github.com/zhupanovdm/go-runtime-monitor/internal/measure"
	"math/rand"
	"runtime"
)

func pollRuntimeMetrics(subscriber chan<- string, PollCount int64) {
	stats := runtime.MemStats{}
	runtime.ReadMemStats(&stats)

	subscriber <- gaugeu("Alloc", stats.Alloc)
	subscriber <- gaugeu("BuckHashSys", stats.BuckHashSys)
	subscriber <- gauge("GCCPUFraction", stats.GCCPUFraction)
	subscriber <- gaugeu("GCSys", stats.GCSys)
	subscriber <- gaugeu("HeapAlloc", stats.HeapAlloc)
	subscriber <- gaugeu("HeapIdle", stats.HeapIdle)
	subscriber <- gaugeu("HeapInuse", stats.HeapInuse)
	subscriber <- gaugeu("HeapObjects", stats.HeapObjects)
	subscriber <- gaugeu("HeapReleased", stats.HeapReleased)
	subscriber <- gaugeu("HeapSys", stats.HeapSys)
	subscriber <- gaugeu("LastGC", stats.LastGC)
	subscriber <- gaugeu("Lookups", stats.Lookups)
	subscriber <- gaugeu("MCacheInuse", stats.MCacheInuse)
	subscriber <- gaugeu("MCacheSys", stats.MCacheSys)
	subscriber <- gaugeu("MSpanInuse", stats.MSpanInuse)
	subscriber <- gaugeu("MSpanSys", stats.MSpanSys)
	subscriber <- gaugeu("Mallocs", stats.Mallocs)
	subscriber <- gaugeu("NextGC", stats.NextGC)
	subscriber <- gaugeu("NumForcedGC", uint64(stats.NumForcedGC))
	subscriber <- gaugeu("NumGC", uint64(stats.NumGC))
	subscriber <- gaugeu("OtherSys", stats.OtherSys)
	subscriber <- gaugeu("PauseTotalNs", stats.PauseTotalNs)
	subscriber <- gaugeu("StackInuse", stats.StackInuse)
	subscriber <- gaugeu("StackSys", stats.StackSys)
	subscriber <- gaugeu("Sys", stats.Sys)
	subscriber <- counter("PollCount", PollCount)
	subscriber <- gauge("RandomValue", rand.Float64())
}

func gauge(name string, value float64) string {
	g := measure.Gauge(value)
	return (&measure.Metric{
		Name:  name,
		Value: &g,
	}).Encode()
}

func gaugeu(name string, value uint64) string {
	g := measure.Gauge(value)
	return (&measure.Metric{
		Name:  name,
		Value: &g,
	}).Encode()
}

func counter(name string, value int64) string {
	c := measure.Counter(value)
	return (&measure.Metric{
		Name:  name,
		Value: &c,
	}).Encode()
}
