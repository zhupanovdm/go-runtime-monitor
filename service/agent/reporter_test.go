package agent

import (
	"context"
	"fmt"
	"testing"

	"github.com/rs/zerolog"

	"github.com/zhupanovdm/go-runtime-monitor/config"
	"github.com/zhupanovdm/go-runtime-monitor/providers/monitor/stub"
)

func BenchmarkMetricsReporter(b *testing.B) {
	zerolog.SetGlobalLevel(zerolog.Disabled)

	count := 1024
	ctx := context.TODO()

	b.Run("Reporter publishing", func(b *testing.B) {
		froze := NewFroze()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			publish(froze, count)
		}
	})

	b.Run("Report bulk", func(b *testing.B) {
		froze := NewFroze()
		publish(froze, count)
		rep := NewMetricsReporter(&config.Config{}, froze, stub.New())
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			rep.Report(ctx)
		}
	})

}

func publish(froze *Froze, count int) {
	for j := 0; j < count; j++ {
		froze.UpdateGauge(fmt.Sprintf("foo%d", j+1), float64(j))
	}
}
