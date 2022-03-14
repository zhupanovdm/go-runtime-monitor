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
		froze      *Froze
		collectors []Collector
		interval   time.Duration
	}

	Collector func(ctx context.Context, froze *Froze) error
)

func (c *metricsCollector) Poll(ctx context.Context) {
	ctx, _ = logging.SetIfAbsentCID(ctx, logging.NewCID())
	ctx, logger := logging.GetOrCreateLogger(ctx, logging.WithService(c), logging.WithCID(ctx))
	logger.Info().Msg("polling metrics")

	c.froze.Lock()
	defer c.froze.Unlock()

	for i, collector := range c.collectors {
		if err := collector(ctx, c.froze); err != nil {
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

func NewMetricsCollector(cfg *config.Config, froze *Froze, collectors ...Collector) CollectorService {
	return &metricsCollector{
		collectors: collectors,
		froze:      froze,
		interval:   cfg.PollInterval,
	}
}
