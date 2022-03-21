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
	expected := metric.List{
		metric.NewGaugeMetric("TotalMemory", metric.Gauge(0)),
		metric.NewGaugeMetric("FreeMemory", metric.Gauge(0)),
	}
	for i := 0; i < runtime.NumCPU(); i++ {
		expected = append(expected, metric.NewGaugeMetric(fmt.Sprintf("CPUutilization%d", i+1), metric.Gauge(0)))
	}

	froze := NewFroze()

	t.Run("Basic test", func(t *testing.T) {
		err := PS()(context.TODO(), froze)

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

func BenchmarkPS(b *testing.B) {
	b.Run("PS polling", func(b *testing.B) {
		b.StopTimer()
		froze := NewFroze()
		collector := PS()
		ctx := context.TODO()
		b.StartTimer()

		for i := 0; i < b.N; i++ {
			require.NoError(b, collector(ctx, froze))
		}
	})
}
