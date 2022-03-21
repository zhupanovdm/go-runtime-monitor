package agent

import (
	"context"
	"fmt"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"

	"github.com/zhupanovdm/go-runtime-monitor/pkg/logging"
)

func PS() Collector {
	return func(ctx context.Context, froze *Froze) error {
		_, logger := logging.GetOrCreateLogger(ctx)

		v, err := mem.VirtualMemoryWithContext(ctx)
		if err != nil {
			logger.Err(err).Msg("failed to collect metrics")
			return err
		}

		froze.UpdateGauge("TotalMemory", float64(v.Total))
		froze.UpdateGauge("FreeMemory", float64(v.Free))

		usage, err := cpu.PercentWithContext(ctx, 0, true)
		if err != nil {
			logger.Err(err).Msg("failed to collect metrics")
			return err
		}
		for i := 0; i < len(usage); i++ {
			froze.UpdateGauge(fmt.Sprintf("CPUutilization%d", i+1), usage[i])
		}
		return nil
	}
}
