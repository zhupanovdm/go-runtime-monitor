package agent

import (
	"context"
	"fmt"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"

	"github.com/zhupanovdm/go-runtime-monitor/config"
	"github.com/zhupanovdm/go-runtime-monitor/model/metric"
	"github.com/zhupanovdm/go-runtime-monitor/pkg/logging"
)

const psName = "ps"

// Collects memory and cpu utilisation
func ps() Collector {
	return func(ctx context.Context, reporter ReporterService) error {
		_, logger := logging.GetOrCreateLogger(ctx)

		v, err := mem.VirtualMemoryWithContext(ctx)
		if err != nil {
			logger.Err(err).Msg("failed to collect metrics")
			return err
		}

		reporter.Publish(ctx, metric.NewGaugeMetric("TotalMemory", metric.Gauge(v.Total)))
		reporter.Publish(ctx, metric.NewGaugeMetric("FreeMemory", metric.Gauge(v.Free)))

		usage, err := cpu.PercentWithContext(ctx, 0, true)
		if err != nil {
			logger.Err(err).Msg("failed to collect metrics")
			return err
		}
		for i := 0; i < len(usage); i++ {
			name := fmt.Sprintf("CPUutilization%d", i+1)
			reporter.Publish(ctx, metric.NewGaugeMetric(name, metric.Gauge(usage[i])))
		}
		return nil
	}
}

func NewPsCollector(cfg *config.Config, reporter ReporterService) CollectorService {
	return NewMetricsCollector(cfg, reporter, ps(), psName)
}
