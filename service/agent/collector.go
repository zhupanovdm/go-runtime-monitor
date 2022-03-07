package agent

import (
	"context"
	"time"

	"github.com/zhupanovdm/go-runtime-monitor/config"
	"github.com/zhupanovdm/go-runtime-monitor/pkg/logging"
	"github.com/zhupanovdm/go-runtime-monitor/pkg/task"
)

var _ CollectorService = (*metricsCollector)(nil)

type (
	metricsCollector struct {
		reporter   ReporterService
		collectors []Collector
		interval   time.Duration
	}

	Collector func(ctx context.Context, reporter ReporterService) error
)

func (c *metricsCollector) Poll(ctx context.Context) {
	ctx, _ = logging.SetIfAbsentCID(ctx, logging.NewCID())
	ctx, logger := logging.GetOrCreateLogger(ctx, logging.WithService(c), logging.WithCID(ctx))
	logger.Info().Msg("polling metrics")

	for i, collector := range c.collectors {
		if err := collector(ctx, c.reporter); err != nil {
			logger.Err(err).Msgf("collector (%d) failed", i)
		}
	}
	logger.Info().Msg("poll completed")
}

func (c *metricsCollector) BackgroundTask() task.Task {
	return task.Task(c.Poll).With(task.PeriodicRun(c.interval))
}

func (c *metricsCollector) Name() string {
	return "Agent metrics collector"
}

func NewMetricsCollector(cfg *config.Config, reporter ReporterService, collectors ...Collector) CollectorService {
	return &metricsCollector{
		collectors: collectors,
		reporter:   reporter,
		interval:   cfg.PollInterval,
	}
}
