package agent

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zhupanovdm/go-runtime-monitor/model/metric"
)

func TestMemStats(t *testing.T) {
	froze := NewFroze()
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
		err := MemStats()(context.TODO(), froze)

		list := froze.List()
		list0 := make(metric.List, 0, len(list))
		for _, m := range froze.List() {
			v, err := m.Type().New()
			require.NoError(t, err)
			list0 = append(list0, &metric.Metric{ID: m.ID, Value: v})
		}
		if assert.NoError(t, err) {
			assert.ElementsMatch(t, list0, expected)
		}
	})
}

func BenchmarkMemStats(b *testing.B) {
	b.Run("Mem Stats polling", func(b *testing.B) {
		b.StopTimer()
		froze := NewFroze()
		collector := MemStats()
		ctx := context.TODO()
		b.StartTimer()

		for i := 0; i < b.N; i++ {
			require.NoError(b, collector(ctx, froze))
		}
	})
}
