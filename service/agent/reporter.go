package agent

import (
	"context"
	"time"

	"github.com/zhupanovdm/go-runtime-monitor/config"
	"github.com/zhupanovdm/go-runtime-monitor/model/metric"
	"github.com/zhupanovdm/go-runtime-monitor/pkg/logging"
	"github.com/zhupanovdm/go-runtime-monitor/pkg/task"
	"github.com/zhupanovdm/go-runtime-monitor/providers/monitor"
)

var _ ReporterService = (*metricsReporter)(nil)

type metricsReporter struct {
	monitor.Provider
	interval time.Duration
	events   chan metricEvent
}

type metricEvent struct {
	*metric.Metric
	CorrelationID string
}

func (r *metricsReporter) Publish(ctx context.Context, mtr *metric.Metric) {
	var cid string
	ctx, cid = logging.SetIfAbsentCID(ctx, logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithService(r), logging.WithCID(ctx))

	logger.UpdateContext(logging.LogCtxFrom(mtr))
	logger.Trace().Msg("publishing metric")

	r.events <- metricEvent{
		Metric:        mtr,
		CorrelationID: cid,
	}
}

func (r *metricsReporter) report(ctx context.Context) error {
	ctx, logger := logging.GetOrCreateLogger(ctx, logging.WithService(r))
	logger.Info().Msg("reporting metrics to monitor")

	for cnt := len(r.events); cnt > 0; cnt-- {
		event := <-r.events

		ctx, _ := logging.SetCID(ctx, event.CorrelationID)
		_, logger := logging.GetOrCreateLogger(ctx, logging.WithCID(ctx))

		logger.UpdateContext(logging.LogCtxFrom(event.Metric))
		logger.Trace().Msg("transporting metric")

		ctx = logging.SetLogger(ctx, logger)
		if err := r.Update(ctx, event.Metric); err != nil {
			logger.Err(err).Msg("metric not sent")
			return err
		}
	}

	logger.Info().Msg("reporting completed")
	return nil
}

func (r *metricsReporter) BackgroundTask() task.Task {
	return task.Task(func(ctx context.Context) { _ = r.report(ctx) }).With(task.PeriodicRun(r.interval))
}

func (r *metricsReporter) Name() string {
	return "Agent metrics reporter"
}

func NewMetricsReporter(cfg *config.Config, provider monitor.Provider) ReporterService {
	return &metricsReporter{
		events:   make(chan metricEvent, cfg.ReportBuffer),
		Provider: provider,
		interval: cfg.ReportInterval,
	}
}
