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
	froze    *Froze
	interval time.Duration
}

func (r *metricsReporter) Report(ctx context.Context) {
	ctx, _ = logging.SetIfAbsentCID(ctx, logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithService(r), logging.WithCID(ctx))
	logger.Info().Msg("reporting metrics to monitor")

	list := r.read()
	if err := r.UpdateBulk(ctx, list); err != nil {
		logger.Err(err).Msg("metrics not sent")
		return
	}
	logger.Info().Msgf("reporting completed (%d events)", len(list))
}

func (r metricsReporter) read() metric.List {
	r.froze.Lock()
	defer r.froze.Unlock()
	return r.froze.List()
}

func (r *metricsReporter) BackgroundTask() task.Task {
	return task.Task(r.Report).With(task.PeriodicRun(r.interval))
}

func (r *metricsReporter) Name() string {
	return "Agent metrics reporter"
}

// NewMetricsReporter creates new metrics reporting service. Each time the service is called to report it will read entirely
// metrics from Froze and send it to monitor.Provider.
func NewMetricsReporter(cfg *config.Config, froze *Froze, provider monitor.Provider) ReporterService {
	return &metricsReporter{
		froze:    froze,
		Provider: provider,
		interval: cfg.ReportInterval,
	}
}
