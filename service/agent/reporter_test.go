package agent

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zhupanovdm/go-runtime-monitor/config"
	"github.com/zhupanovdm/go-runtime-monitor/model/metric"
	"github.com/zhupanovdm/go-runtime-monitor/providers/monitor/stub"
)

func Benchmark_metricsReporter(b *testing.B) {
	count := 1024

	b.Run("Reporter publishing", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			rep := NewMetricsReporter(&config.Config{ReportBuffer: count + 1}, stub.New())
			ctx := context.TODO()
			b.StartTimer()

			publish(ctx, rep, count)
		}
	})

	b.Run("Report bulk", func(b *testing.B) {
		b.StopTimer()
		rep := NewMetricsReporter(&config.Config{ReportBuffer: count + 1}, stub.New())
		ctx := context.TODO()
		publish(ctx, rep, count)
		b.StartTimer()

		for i := 0; i < b.N; i++ {
			require.NoError(b, rep.ReportBulk(ctx))
		}
	})

}

func publish(ctx context.Context, rep ReporterService, count int) {
	for j := 0; j < count; j++ {
		rep.Publish(ctx, metric.NewGaugeMetric("foo", metric.Gauge(0)))
	}
}
