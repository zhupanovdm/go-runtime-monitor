package agent

import (
	"context"
	"math/rand"
	"runtime"
	"time"

	"github.com/rs/zerolog"

	"github.com/zhupanovdm/go-runtime-monitor/config"
	"github.com/zhupanovdm/go-runtime-monitor/model/metric"
	"github.com/zhupanovdm/go-runtime-monitor/pkg/logging"
	"github.com/zhupanovdm/go-runtime-monitor/pkg/task"
)

var _ CollectorService = (*runtimeMetricsCollector)(nil)
var _ logging.LogCtxProvider = (*runtimeMetricsCollector)(nil)

type runtimeMetricsCollector struct {
	reporter ReporterService
	counter  int64
	interval time.Duration
}

func (c *runtimeMetricsCollector) poll(ctx context.Context) error {
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithService(c))

	c.counter++

	logger.UpdateContext(logging.LogCtxFrom(c))
	logger.Info().Msg("polling runtime metrics")

	ctx = logging.SetLogger(ctx, logger)

	stats := &runtime.MemStats{}
	runtime.ReadMemStats(stats)

	c.reporter.Publish(ctx, metric.NewGaugeMetric("Alloc", metric.Gauge(stats.Alloc)))
	c.reporter.Publish(ctx, metric.NewGaugeMetric("BuckHashSys", metric.Gauge(stats.BuckHashSys)))
	c.reporter.Publish(ctx, metric.NewGaugeMetric("GCCPUFraction", metric.Gauge(stats.GCCPUFraction)))
	c.reporter.Publish(ctx, metric.NewGaugeMetric("GCSys", metric.Gauge(stats.GCSys)))
	c.reporter.Publish(ctx, metric.NewGaugeMetric("HeapAlloc", metric.Gauge(stats.HeapAlloc)))
	c.reporter.Publish(ctx, metric.NewGaugeMetric("HeapIdle", metric.Gauge(stats.HeapIdle)))
	c.reporter.Publish(ctx, metric.NewGaugeMetric("HeapInuse", metric.Gauge(stats.HeapInuse)))
	c.reporter.Publish(ctx, metric.NewGaugeMetric("HeapObjects", metric.Gauge(stats.HeapObjects)))
	c.reporter.Publish(ctx, metric.NewGaugeMetric("HeapReleased", metric.Gauge(stats.HeapReleased)))
	c.reporter.Publish(ctx, metric.NewGaugeMetric("HeapSys", metric.Gauge(stats.HeapSys)))
	c.reporter.Publish(ctx, metric.NewGaugeMetric("LastGC", metric.Gauge(stats.LastGC)))
	c.reporter.Publish(ctx, metric.NewGaugeMetric("Lookups", metric.Gauge(stats.Lookups)))
	c.reporter.Publish(ctx, metric.NewGaugeMetric("MCacheInuse", metric.Gauge(stats.MCacheInuse)))
	c.reporter.Publish(ctx, metric.NewGaugeMetric("MCacheSys", metric.Gauge(stats.MCacheSys)))
	c.reporter.Publish(ctx, metric.NewGaugeMetric("MSpanInuse", metric.Gauge(stats.MSpanInuse)))
	c.reporter.Publish(ctx, metric.NewGaugeMetric("MSpanSys", metric.Gauge(stats.MSpanSys)))
	c.reporter.Publish(ctx, metric.NewGaugeMetric("Mallocs", metric.Gauge(stats.Mallocs)))
	c.reporter.Publish(ctx, metric.NewGaugeMetric("NextGC", metric.Gauge(stats.NextGC)))
	c.reporter.Publish(ctx, metric.NewGaugeMetric("NumForcedGC", metric.Gauge(stats.NumForcedGC)))
	c.reporter.Publish(ctx, metric.NewGaugeMetric("NumGC", metric.Gauge(stats.NumGC)))
	c.reporter.Publish(ctx, metric.NewGaugeMetric("OtherSys", metric.Gauge(stats.OtherSys)))
	c.reporter.Publish(ctx, metric.NewGaugeMetric("PauseTotalNs", metric.Gauge(stats.PauseTotalNs)))
	c.reporter.Publish(ctx, metric.NewGaugeMetric("StackInuse", metric.Gauge(stats.StackInuse)))
	c.reporter.Publish(ctx, metric.NewGaugeMetric("StackSys", metric.Gauge(stats.StackSys)))
	c.reporter.Publish(ctx, metric.NewGaugeMetric("Sys", metric.Gauge(stats.Sys)))

	c.reporter.Publish(ctx, metric.NewGaugeMetric("RandomValue", metric.Gauge(rand.Float64())))

	c.reporter.Publish(ctx, metric.NewGaugeMetric("Frees", metric.Gauge(stats.Frees)))
	c.reporter.Publish(ctx, metric.NewGaugeMetric("TotalAlloc", metric.Gauge(stats.TotalAlloc)))

	c.reporter.Publish(ctx, metric.NewCounterMetric("PollCount", metric.Counter(c.counter)))

	logger.Info().Msg("poll completed")

	return nil
}

func (c *runtimeMetricsCollector) BackgroundTask() task.Task {
	return task.Task(func(ctx context.Context) { _ = c.poll(ctx) }).With(task.PeriodicRun(c.interval))
}

func (c *runtimeMetricsCollector) Name() string {
	return "Agent metrics collector"
}

func (c *runtimeMetricsCollector) LoggerCtx(ctx zerolog.Context) zerolog.Context {
	return ctx.Int64(logging.PollCountKey, c.counter)
}

func NewRuntimeMetricsCollector(cfg *config.Config, reporter ReporterService) CollectorService {
	rand.Seed(time.Now().UnixNano())
	return &runtimeMetricsCollector{
		reporter: reporter,
		interval: cfg.PollInterval,
	}
}
