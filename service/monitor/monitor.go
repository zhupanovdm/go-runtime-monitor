package monitor

import (
	"context"
	"fmt"
	"time"

	"github.com/zhupanovdm/go-runtime-monitor/config"
	"github.com/zhupanovdm/go-runtime-monitor/model/metric"
	"github.com/zhupanovdm/go-runtime-monitor/pkg/logging"
	"github.com/zhupanovdm/go-runtime-monitor/pkg/task"
	"github.com/zhupanovdm/go-runtime-monitor/storage"
)

var _ Monitor = (*monitor)(nil)

type monitor struct {
	interval time.Duration
	restore  bool
	dump     storage.Storage
	gauge    storage.GaugeStorage
	counter  storage.CounterStorage
}

func (m *monitor) Restore(ctx context.Context) error {
	ctx, _ = logging.SetIfAbsentCID(ctx, logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithService(m), logging.WithCID(ctx))

	logger.Info().Msg("serving [Restore]")

	if !m.restore {
		logger.Warn().Msg("restore: feature is disabled")
		return nil
	}
	if m.dump == nil {
		logger.Warn().Msg("restore: dump storage is not set")
		return nil
	}

	metrics, err := m.dump.GetAll(ctx)
	if err != nil {
		logger.Err(err).Msg("restore: failed to read from dump")
		return err
	}

	for typ, metrics := range metrics.ToMap() {
		switch typ {
		case metric.GaugeType:
			err = m.gauge.UpdateBulk(ctx, metrics)
		case metric.CounterType:
			err = m.counter.UpdateBulk(ctx, metrics)
		default:
			err = fmt.Errorf("unknown metric type %v", typ)
		}
		if err != nil {
			logger.Err(err).Msgf("restore: update %v storage failed", typ)
			return err
		}
	}
	return nil
}

func (m *monitor) Update(ctx context.Context, mtr *metric.Metric) (err error) {
	ctx, _ = logging.SetIfAbsentCID(ctx, logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithService(m), logging.WithCID(ctx))

	logger.UpdateContext(logging.LogCtxFrom(mtr))

	logger.Info().Msg("serving [Update]")

	ctx = logging.SetLogger(ctx, logger)
	switch value := mtr.Value.(type) {
	case *metric.Gauge:
		err = m.gauge.Update(ctx, mtr.ID, *value)
	case *metric.Counter:
		err = m.counter.Update(ctx, mtr.ID, *value)
	default:
		err = fmt.Errorf("type is not supported: %T", value)
		logger.Err(err).Msg("update: failed to update")
	}

	if err != nil && m.isSyncDump() {
		if err = m.store(ctx); err != nil {
			logger.Err(err).Msg("update: failed to sync dump")
		}
	}
	return
}

func (m *monitor) Get(ctx context.Context, id string, typ metric.Type) (value *metric.Metric, err error) {
	ctx, _ = logging.SetIfAbsentCID(ctx, logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithService(m), logging.WithCID(ctx))

	logger.UpdateContext(logging.LogCtxKeyStr(logging.MetricIDKey, id))
	logger.UpdateContext(logging.LogCtxFrom(typ))

	logger.Info().Msg("serving [Get]")

	ctx = logging.SetLogger(ctx, logger)
	switch typ {
	case metric.GaugeType:
		return m.gauge.Get(ctx, id)
	case metric.CounterType:
		return m.counter.Get(ctx, id)
	default:
		err = fmt.Errorf("unknown metric type: %v", typ)
		logger.Err(err).Msg("get: failed to read metric")
	}
	return
}

func (m *monitor) GetAll(ctx context.Context) (metric.List, error) {
	ctx, _ = logging.SetIfAbsentCID(ctx, logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithService(m), logging.WithCID(ctx))
	logger.Info().Msg("serving [GetAll]")

	return m.readAll(logging.SetLogger(ctx, logger))
}

func (m *monitor) Ping(ctx context.Context) error {
	ctx, _ = logging.SetIfAbsentCID(ctx, logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithService(m), logging.WithCID(ctx))
	logger.Info().Msg("serving [Ping]")

	if m.dump == nil {
		logger.Warn().Msg("ping: dump storage is not set")
		return nil
	}

	return m.dump.Ping(logging.SetLogger(ctx, logger))
}

func (m *monitor) store(ctx context.Context) error {
	ctx, _ = logging.SetIfAbsentCID(ctx, logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithService(m), logging.WithCID(ctx))
	logger.Info().Msg("serving [store]")

	if m.dump == nil {
		logger.Warn().Msg("store: dump storage is not set")
		return nil
	}

	metrics, err := m.readAll(logging.SetLogger(ctx, logger))
	if err != nil {
		logger.Err(err).Msg("store: failed to read metrics")
		return err
	}

	if err = m.dump.UpdateBulk(ctx, metrics); err != nil {
		logger.Err(err).Msg("store: dump update failed")
	}

	return nil
}

func (m *monitor) BackgroundTask() task.Task {
	if m.isSyncDump() {
		return task.VoidTask
	}
	return task.Task(func(ctx context.Context) { _ = m.store(ctx) }).With(task.PeriodicRun(m.interval))
}

func (m *monitor) Name() string {
	return "Monitor service"
}

func (m *monitor) isSyncDump() bool {
	return m.interval == 0
}

func (m *monitor) readAll(ctx context.Context) (list metric.List, err error) {
	_, logger := logging.GetOrCreateLogger(ctx)

	if list, err = m.gauge.GetAll(ctx); err != nil {
		logger.Err(err).Msg("read all: gauges read failed")
		return
	}

	var counters metric.List
	if counters, err = m.counter.GetAll(ctx); err != nil {
		logger.Err(err).Msg("read all: counters read failed")
		return
	}

	list = append(list, counters...)

	logger.Trace().Msgf("read all: got %d records total", len(list))

	return
}

func NewMonitor(cfg *config.Config, dump storage.Storage, gauges storage.GaugeStorage, counters storage.CounterStorage) Monitor {
	return &monitor{
		interval: cfg.StoreInterval,
		restore:  cfg.Restore,
		dump:     dump,
		gauge:    gauges,
		counter:  counters,
	}
}
