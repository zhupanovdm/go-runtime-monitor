package monitor

import (
	"context"
	"fmt"

	"github.com/zhupanovdm/go-runtime-monitor/model/metric"
	"github.com/zhupanovdm/go-runtime-monitor/pkg/logging"
	"github.com/zhupanovdm/go-runtime-monitor/storage"
)

var _ MetricsMonitorService = (*monitorSvc)(nil)

type monitorSvc struct {
	gauge   storage.GaugeStorage
	counter storage.CounterStorage
}

func (m *monitorSvc) Update(ctx context.Context, mtr *metric.Metric) (err error) {
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
		err = fmt.Errorf("type is not supported yet: %T", value)
		logger.Err(err).Msg("metric is not updated")
	}
	return
}

func (m *monitorSvc) Get(ctx context.Context, id string, typ metric.Type) (value *metric.Metric, err error) {
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
		logger.Err(err).Msg("metric is not read")
	}
	return
}

func (m *monitorSvc) GetAll(ctx context.Context) (list []*metric.Metric, err error) {
	ctx, _ = logging.SetIfAbsentCID(ctx, logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithService(m), logging.WithCID(ctx))
	logger.Info().Msg("serving [GetAll]")

	if list, err = m.gauge.GetAll(ctx); err != nil {
		logger.Err(err).Msg("gauges read failed")
		return
	}
	logger.Trace().Msgf("got %d gauges", len(list))

	var counters []*metric.Metric
	if counters, err = m.counter.GetAll(ctx); err != nil {
		logger.Err(err).Msg("counters read failed")
		return
	}
	logger.Trace().Msgf("got %d counters", len(counters))

	list = append(list, counters...)

	logger.Trace().Msgf("got %d records total", len(list))
	return
}

func (m *monitorSvc) Name() string {
	return "Monitor service"
}

func NewMetricsMonitor(gaugeStorage storage.GaugeStorage, counterStorage storage.CounterStorage) MetricsMonitorService {
	return &monitorSvc{
		gauge:   gaugeStorage,
		counter: counterStorage,
	}
}
