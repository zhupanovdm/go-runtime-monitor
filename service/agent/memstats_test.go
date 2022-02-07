package agent

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zhupanovdm/go-runtime-monitor/model/metric"
)

func Test_memStats(t *testing.T) {
	actual := make(metric.List, 0)
	stub := NewStubReporter(t, func(m *metric.Metric) { actual = append(actual, m) })
	expected := metric.List{
		metric.NewGaugeMetric("Alloc", metric.Gauge(0)),
		metric.NewGaugeMetric("BuckHashSys", metric.Gauge(0)),
		metric.NewGaugeMetric("GCCPUFraction", metric.Gauge(0)),
		metric.NewGaugeMetric("GCSys", metric.Gauge(0)),
		metric.NewGaugeMetric("HeapAlloc", metric.Gauge(0)),
		metric.NewGaugeMetric("HeapIdle", metric.Gauge(0)),
		metric.NewGaugeMetric("HeapInuse", metric.Gauge(0)),
		metric.NewGaugeMetric("HeapObjects", metric.Gauge(0)),
		metric.NewGaugeMetric("HeapReleased", metric.Gauge(0)),
		metric.NewGaugeMetric("HeapSys", metric.Gauge(0)),
		metric.NewGaugeMetric("LastGC", metric.Gauge(0)),
		metric.NewGaugeMetric("Lookups", metric.Gauge(0)),
		metric.NewGaugeMetric("MCacheInuse", metric.Gauge(0)),
		metric.NewGaugeMetric("MCacheSys", metric.Gauge(0)),
		metric.NewGaugeMetric("MSpanInuse", metric.Gauge(0)),
		metric.NewGaugeMetric("MSpanSys", metric.Gauge(0)),
		metric.NewGaugeMetric("Mallocs", metric.Gauge(0)),
		metric.NewGaugeMetric("NextGC", metric.Gauge(0)),
		metric.NewGaugeMetric("NumForcedGC", metric.Gauge(0)),
		metric.NewGaugeMetric("NumGC", metric.Gauge(0)),
		metric.NewGaugeMetric("OtherSys", metric.Gauge(0)),
		metric.NewGaugeMetric("PauseTotalNs", metric.Gauge(0)),
		metric.NewGaugeMetric("StackInuse", metric.Gauge(0)),
		metric.NewGaugeMetric("StackSys", metric.Gauge(0)),
		metric.NewGaugeMetric("Sys", metric.Gauge(0)),
		metric.NewGaugeMetric("RandomValue", metric.Gauge(0)),
		metric.NewGaugeMetric("Frees", metric.Gauge(0)),
		metric.NewGaugeMetric("TotalAlloc", metric.Gauge(0)),
		metric.NewCounterMetric("PollCount", metric.Counter(0)),
	}

	t.Run("Basic test", func(t *testing.T) {
		err := memStats()(context.TODO(), stub)
		if assert.NoError(t, err) {
			assert.ElementsMatch(t, actual, expected)
		}
	})
}
