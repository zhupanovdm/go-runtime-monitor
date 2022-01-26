package monitor

import (
	"context"
	"time"

	"github.com/zhupanovdm/go-runtime-monitor/config"
	"github.com/zhupanovdm/go-runtime-monitor/model/metric"
	"github.com/zhupanovdm/go-runtime-monitor/pkg/logging"
	"github.com/zhupanovdm/go-runtime-monitor/pkg/task"
	"github.com/zhupanovdm/go-runtime-monitor/storage"
)

var _ Monitor = (*monitor)(nil)

type monitor struct {
	interval      time.Duration
	restore       bool
	dumpStorage   storage.Storage
	metricStorage storage.Storage
}

func (m *monitor) Restore(ctx context.Context) error {
	ctx, _ = logging.SetIfAbsentCID(ctx, logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithService(m), logging.WithCID(ctx))
	logger.Info().Msg("serving [Restore]")

	if !m.restore {
		logger.Warn().Msg("restore: feature is disabled")
		return nil
	}
	if m.dumpStorage == nil {
		logger.Warn().Msg("restore: dump storage is not set")
		return nil
	}
	metrics, err := m.dumpStorage.GetAll(ctx)
	if err != nil {
		logger.Err(err).Msg("restore: failed to read from dump")
		return err
	}
	if m.metricStorage.IsPersistent() {
		if err := m.metricStorage.Clear(ctx); err != nil {
			logger.Err(err).Msg("restore: failed to clear metrics storage")
			return err
		}
	}
	if err := m.metricStorage.UpdateBulk(ctx, metrics); err != nil {
		logger.Err(err).Msg("restore: update storage failed")
		return err
	}
	return nil
}

func (m *monitor) Update(ctx context.Context, mtr *metric.Metric) error {
	ctx, _ = logging.SetIfAbsentCID(ctx, logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithService(m), logging.WithCID(ctx))
	logger.UpdateContext(logging.LogCtxFrom(mtr))
	logger.Info().Msg("serving [Update]")

	ctx = logging.SetLogger(ctx, logger)
	if err := m.metricStorage.Update(ctx, mtr.ID, mtr); err != nil {
		logger.Err(err).Msg("update: failed to update storage")
		return err
	}
	if m.isSyncDump() {
		if err := m.Dump(ctx); err != nil {
			logger.Err(err).Msg("update: failed to sync dump")
			return err
		}
	}
	return nil
}

func (m *monitor) Get(ctx context.Context, id string, typ metric.Type) (*metric.Metric, error) {
	ctx, _ = logging.SetIfAbsentCID(ctx, logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithService(m), logging.WithCID(ctx))
	logger.UpdateContext(logging.LogCtxKeyStr(logging.MetricIDKey, id))
	logger.UpdateContext(logging.LogCtxFrom(typ))
	logger.Info().Msg("serving [Get]")

	return m.metricStorage.Get(logging.SetLogger(ctx, logger), id, typ)
}

func (m *monitor) GetAll(ctx context.Context) (metric.List, error) {
	ctx, _ = logging.SetIfAbsentCID(ctx, logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithService(m), logging.WithCID(ctx))
	logger.Info().Msg("serving [GetAll]")

	return m.metricStorage.GetAll(logging.SetLogger(ctx, logger))
}

func (m *monitor) Ping(ctx context.Context) error {
	ctx, _ = logging.SetIfAbsentCID(ctx, logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithService(m), logging.WithCID(ctx))
	logger.Info().Msg("serving [Ping]")

	return m.metricStorage.Ping(logging.SetLogger(ctx, logger))
}

func (m *monitor) Dump(ctx context.Context) error {
	ctx, _ = logging.SetIfAbsentCID(ctx, logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithService(m), logging.WithCID(ctx))
	ctx = logging.SetLogger(ctx, logger)
	logger.Info().Msg("serving [store]")

	if m.dumpStorage == nil {
		logger.Warn().Msg("store: dump storage is not set")
		return nil
	}

	metrics, err := m.metricStorage.GetAll(ctx)
	if err != nil {
		logger.Err(err).Msg("store: failed to read metrics")
		return err
	}
	if err = m.dumpStorage.UpdateBulk(ctx, metrics); err != nil {
		logger.Err(err).Msg("store: dump update failed")
	}
	return nil
}

func (m *monitor) BackgroundTask() task.Task {
	if m.isSyncDump() {
		return task.VoidTask
	}
	return task.Task(func(ctx context.Context) { _ = m.Dump(ctx) }).With(task.PeriodicRun(m.interval))
}

func (m *monitor) isSyncDump() bool {
	return m.interval == 0
}

func (m *monitor) Name() string {
	return "Monitor service"
}

func NewMonitor(cfg *config.Config, dumpStorage storage.Storage, metricStorage storage.Storage) Monitor {
	return &monitor{
		interval:      cfg.StoreInterval,
		restore:       cfg.Restore,
		dumpStorage:   dumpStorage,
		metricStorage: metricStorage,
	}
}
