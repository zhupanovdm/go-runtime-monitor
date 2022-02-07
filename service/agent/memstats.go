package agent

import (
	"context"
	"math/rand"
	"runtime"
	"time"

	"github.com/zhupanovdm/go-runtime-monitor/config"
	"github.com/zhupanovdm/go-runtime-monitor/model/metric"
	"github.com/zhupanovdm/go-runtime-monitor/pkg/logging"
)

const memStatsName = "MemStats"

func memStats() Collector {
	var counter int64

	rand.Seed(time.Now().UnixNano())

	return func(ctx context.Context, reporter ReporterService) error {
		_, logger := logging.GetOrCreateLogger(ctx)

		counter++
		ctx = logging.SetLogger(ctx, logger.With().Int64(logging.PollCountKey, counter).Logger())

		stats := &runtime.MemStats{}
		runtime.ReadMemStats(stats)

		reporter.Publish(ctx, metric.NewGaugeMetric("Alloc", metric.Gauge(stats.Alloc)))
		reporter.Publish(ctx, metric.NewGaugeMetric("BuckHashSys", metric.Gauge(stats.BuckHashSys)))
		reporter.Publish(ctx, metric.NewGaugeMetric("GCCPUFraction", metric.Gauge(stats.GCCPUFraction)))
		reporter.Publish(ctx, metric.NewGaugeMetric("GCSys", metric.Gauge(stats.GCSys)))
		reporter.Publish(ctx, metric.NewGaugeMetric("HeapAlloc", metric.Gauge(stats.HeapAlloc)))
		reporter.Publish(ctx, metric.NewGaugeMetric("HeapIdle", metric.Gauge(stats.HeapIdle)))
		reporter.Publish(ctx, metric.NewGaugeMetric("HeapInuse", metric.Gauge(stats.HeapInuse)))
		reporter.Publish(ctx, metric.NewGaugeMetric("HeapObjects", metric.Gauge(stats.HeapObjects)))
		reporter.Publish(ctx, metric.NewGaugeMetric("HeapReleased", metric.Gauge(stats.HeapReleased)))
		reporter.Publish(ctx, metric.NewGaugeMetric("HeapSys", metric.Gauge(stats.HeapSys)))
		reporter.Publish(ctx, metric.NewGaugeMetric("LastGC", metric.Gauge(stats.LastGC)))
		reporter.Publish(ctx, metric.NewGaugeMetric("Lookups", metric.Gauge(stats.Lookups)))
		reporter.Publish(ctx, metric.NewGaugeMetric("MCacheInuse", metric.Gauge(stats.MCacheInuse)))
		reporter.Publish(ctx, metric.NewGaugeMetric("MCacheSys", metric.Gauge(stats.MCacheSys)))
		reporter.Publish(ctx, metric.NewGaugeMetric("MSpanInuse", metric.Gauge(stats.MSpanInuse)))
		reporter.Publish(ctx, metric.NewGaugeMetric("MSpanSys", metric.Gauge(stats.MSpanSys)))
		reporter.Publish(ctx, metric.NewGaugeMetric("Mallocs", metric.Gauge(stats.Mallocs)))
		reporter.Publish(ctx, metric.NewGaugeMetric("NextGC", metric.Gauge(stats.NextGC)))
		reporter.Publish(ctx, metric.NewGaugeMetric("NumForcedGC", metric.Gauge(stats.NumForcedGC)))
		reporter.Publish(ctx, metric.NewGaugeMetric("NumGC", metric.Gauge(stats.NumGC)))
		reporter.Publish(ctx, metric.NewGaugeMetric("OtherSys", metric.Gauge(stats.OtherSys)))
		reporter.Publish(ctx, metric.NewGaugeMetric("PauseTotalNs", metric.Gauge(stats.PauseTotalNs)))
		reporter.Publish(ctx, metric.NewGaugeMetric("StackInuse", metric.Gauge(stats.StackInuse)))
		reporter.Publish(ctx, metric.NewGaugeMetric("StackSys", metric.Gauge(stats.StackSys)))
		reporter.Publish(ctx, metric.NewGaugeMetric("Sys", metric.Gauge(stats.Sys)))

		reporter.Publish(ctx, metric.NewGaugeMetric("RandomValue", metric.Gauge(rand.Float64())))

		reporter.Publish(ctx, metric.NewGaugeMetric("Frees", metric.Gauge(stats.Frees)))
		reporter.Publish(ctx, metric.NewGaugeMetric("TotalAlloc", metric.Gauge(stats.TotalAlloc)))

		reporter.Publish(ctx, metric.NewCounterMetric("PollCount", metric.Counter(counter)))

		return nil
	}
}

func NewMemStatsCollector(cfg *config.Config, reporter ReporterService) CollectorService {
	return NewMetricsCollector(cfg, reporter, memStats(), memStatsName)
}
