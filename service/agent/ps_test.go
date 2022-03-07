package agent

import (
	"context"
	"fmt"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zhupanovdm/go-runtime-monitor/model/metric"
)

func TestPS(t *testing.T) {
	actual := make(metric.List, 0)
	stub := NewStubReporter(t, func(m *metric.Metric) { actual = append(actual, m) })
	expected := metric.List{
		metric.NewGaugeMetric("TotalMemory", metric.Gauge(0)),
		metric.NewGaugeMetric("FreeMemory", metric.Gauge(0)),
	}
	for i := 0; i < runtime.NumCPU(); i++ {
		expected = append(expected, metric.NewGaugeMetric(fmt.Sprintf("CPUutilization%d", i+1), metric.Gauge(0)))
	}

	t.Run("Basic test", func(t *testing.T) {
		err := PS()(context.TODO(), stub)
		if assert.NoError(t, err) {
			assert.ElementsMatch(t, actual, expected)
		}
	})
}

func BenchmarkPS(b *testing.B) {
	b.Run("PS polling", func(b *testing.B) {
		b.StopTimer()
		stub := NewStubReporter(b, func(*metric.Metric) {})
		collector := PS()
		ctx := context.TODO()
		b.StartTimer()

		for i := 0; i < b.N; i++ {
			require.NoError(b, collector(ctx, stub))
		}
	})
}
