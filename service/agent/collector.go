package agent

import (
	"context"
	"fmt"
	"time"

	"github.com/zhupanovdm/go-runtime-monitor/config"
	"github.com/zhupanovdm/go-runtime-monitor/pkg/logging"
	"github.com/zhupanovdm/go-runtime-monitor/pkg/task"
)

var _ CollectorService = (*metricsCollector)(nil)

type (
	metricsCollector struct {
		name      string
		reporter  ReporterService
		collector Collector
		interval  time.Duration
	}

	Collector func(ctx context.Context, reporter ReporterService) error
)

func (c *metricsCollector) Poll(ctx context.Context) error {
	ctx, _ = logging.SetIfAbsentCID(ctx, logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithService(c), logging.WithCID(ctx))
	logger.Info().Msg("polling metrics")

	err := c.collector(logging.SetLogger(ctx, logger), c.reporter)

	logger.Info().Msg("poll completed")
	return err
}

func (c *metricsCollector) BackgroundTask() task.Task {
	return task.Task(func(ctx context.Context) { _ = c.Poll(ctx) }).With(task.PeriodicRun(c.interval))
}

func (c *metricsCollector) Name() string {
	return fmt.Sprintf("Agent metrics collector: %s", c.name)
}

func NewMetricsCollector(cfg *config.Config, reporter ReporterService, collector Collector, name string) CollectorService {
	return &metricsCollector{
		collector: collector,
		name:      name,
		reporter:  reporter,
		interval:  cfg.PollInterval,
	}
}
