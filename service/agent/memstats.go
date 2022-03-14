package agent

import (
	"context"
	"math/rand"
	"runtime"
	"time"
)

func MemStats() Collector {
	var counter int64
	var stats runtime.MemStats

	rand.Seed(time.Now().UnixNano())

	return func(ctx context.Context, froze *Froze) error {
		counter++
		runtime.ReadMemStats(&stats)

		froze.UpdateGauge("Alloc", float64(stats.Alloc))
		froze.UpdateGauge("BuckHashSys", float64(stats.BuckHashSys))
		froze.UpdateGauge("GCCPUFraction", stats.GCCPUFraction)
		froze.UpdateGauge("GCSys", float64(stats.GCSys))
		froze.UpdateGauge("HeapAlloc", float64(stats.HeapAlloc))
		froze.UpdateGauge("HeapIdle", float64(stats.HeapIdle))
		froze.UpdateGauge("HeapInuse", float64(stats.HeapInuse))
		froze.UpdateGauge("HeapObjects", float64(stats.HeapObjects))
		froze.UpdateGauge("HeapSys", float64(stats.HeapSys))
		froze.UpdateGauge("HeapReleased", float64(stats.HeapReleased))
		froze.UpdateGauge("LastGC", float64(stats.LastGC))
		froze.UpdateGauge("Lookups", float64(stats.Lookups))
		froze.UpdateGauge("MCacheSys", float64(stats.MCacheSys))
		froze.UpdateGauge("MCacheInuse", float64(stats.MCacheInuse))
		froze.UpdateGauge("MSpanInuse", float64(stats.MSpanInuse))
		froze.UpdateGauge("MSpanSys", float64(stats.MSpanSys))
		froze.UpdateGauge("Mallocs", float64(stats.Mallocs))
		froze.UpdateGauge("NextGC", float64(stats.NextGC))
		froze.UpdateGauge("NumForcedGC", float64(stats.NumForcedGC))
		froze.UpdateGauge("NumGC", float64(stats.NumGC))
		froze.UpdateGauge("OtherSys", float64(stats.OtherSys))
		froze.UpdateGauge("PauseTotalNs", float64(stats.PauseTotalNs))
		froze.UpdateGauge("StackInuse", float64(stats.StackInuse))
		froze.UpdateGauge("StackSys", float64(stats.StackSys))
		froze.UpdateGauge("Sys", float64(stats.Sys))

		froze.UpdateGauge("RandomValue", rand.Float64())

		froze.UpdateGauge("Frees", float64(stats.Frees))
		froze.UpdateGauge("TotalAlloc", float64(stats.TotalAlloc))

		froze.UpdateCounter("PollCount", counter)
		return nil
	}
}
