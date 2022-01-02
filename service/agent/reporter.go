package agent

import (
	"context"
	"fmt"
	"github.com/zhupanovdm/go-runtime-monitor/config"
	"github.com/zhupanovdm/go-runtime-monitor/model/metric"
	"github.com/zhupanovdm/go-runtime-monitor/pkg/task"
	"time"

	"github.com/zhupanovdm/go-runtime-monitor/providers/monitor"
)

var _ ReporterService = (*metricsReporter)(nil)

type metricsReporter struct {
	monitor.Provider
	interval time.Duration
	pipe     chan *metric.Metric
}

func (r *metricsReporter) Publish(_ context.Context, mtr *metric.Metric) {
	r.pipe <- mtr
}

func (r *metricsReporter) report(ctx context.Context) error {
	for cnt := len(r.pipe); cnt > 0; cnt-- {
		if err := r.Update(ctx, <-r.pipe); err != nil {
			return fmt.Errorf("error while reporting metrics to monitor: %w", err)
		}
	}
	return nil
}

func (r *metricsReporter) BackgroundTask() task.Task {
	return task.Task(func(ctx context.Context) { _ = r.report(ctx) }).With(task.PeriodicRun(r.interval))
}

func NewMetricsReporter(cfg *config.Config, provider monitor.Provider) ReporterService {
	return &metricsReporter{
		pipe:     make(chan *metric.Metric, cfg.ReporterBufferSize),
		Provider: provider,
		interval: cfg.ReportInterval,
	}
}
