package trivial

import (
	"context"
	"fmt"
	"sync"

	"github.com/zhupanovdm/go-runtime-monitor/config"
	"github.com/zhupanovdm/go-runtime-monitor/model/metric"
	"github.com/zhupanovdm/go-runtime-monitor/pkg/logging"
	"github.com/zhupanovdm/go-runtime-monitor/storage"
)

const trivialStorageName = "Trivial storage"

var _ storage.Storage = (*client)(nil)

type client struct {
	sync.RWMutex
	gauges   map[string]metric.Gauge
	counters map[string]metric.Counter
}

func (s *client) IsPersistent() bool {
	return false
}

func (s *client) Init(context.Context) error {
	return nil
}

func (s *client) Update(ctx context.Context, id string, value metric.Value) error {
	ctx, _ = logging.SetIfAbsentCID(ctx, logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(trivialStorageName), logging.WithCID(ctx))

	if err := value.Type().Validate(); err != nil {
		err := fmt.Errorf("unknown metric type: %v", value.Type())
		logger.Err(err).Msg("update failed")
		return err
	}

	s.Lock()
	defer s.Unlock()
	s.save(id, value)

	logger.Trace().Msgf("metric [%s]: updated with [%v]", id, value)
	return nil
}

func (s *client) UpdateBulk(ctx context.Context, list metric.List) error {
	ctx, _ = logging.SetIfAbsentCID(ctx, logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(trivialStorageName), logging.WithCID(ctx))

	s.Lock()
	defer s.Unlock()
	for _, mtr := range list {
		if err := mtr.Type().Validate(); err != nil {
			err := fmt.Errorf("unknown metric type: %v", mtr.Type())
			logger.Err(err).Msg("update failed")
			return err
		}
		s.save(mtr.ID, mtr.Value)
	}

	logger.Trace().Msgf("%d records updated", len(list))
	return nil
}

func (s *client) Get(ctx context.Context, id string, typ metric.Type) (*metric.Metric, error) {
	ctx, _ = logging.SetIfAbsentCID(ctx, logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(trivialStorageName), logging.WithCID(ctx))

	s.RLock()
	defer s.RUnlock()

	switch typ {
	case metric.GaugeType:
		if value, ok := s.gauges[id]; ok {
			logger.Trace().Msgf("gauge [%s]: restored [%f]", id, value)
			return metric.NewGaugeMetric(id, value), nil
		}
	case metric.CounterType:
		if value, ok := s.counters[id]; ok {
			logger.Trace().Msgf("counter [%s]: restored [%d]", id, value)
			return metric.NewCounterMetric(id, value), nil
		}
	}

	logger.Trace().Msgf("counter [%s]: not found", id)
	return nil, nil
}

func (s *client) GetAll(ctx context.Context) (list metric.List, _ error) {
	ctx, _ = logging.SetIfAbsentCID(ctx, logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(trivialStorageName), logging.WithCID(ctx))

	list = make([]*metric.Metric, 0, len(s.gauges)+len(s.counters))

	s.RLock()
	defer s.RUnlock()
	for k, v := range s.gauges {
		list = append(list, metric.NewGaugeMetric(k, v))
	}
	for k, v := range s.counters {
		list = append(list, metric.NewCounterMetric(k, v))
	}

	logger.Trace().Msgf("counter: %d records read", len(list))
	return
}

func (s *client) Ping(context.Context) error {
	return nil
}

func (s *client) Close(context.Context) {}

func (s *client) save(id string, value metric.Value) {
	if m, ok := value.(*metric.Metric); ok {
		value = m.Value
	}
	switch value.Type() {
	case metric.GaugeType:
		s.gauges[id] = *value.(*metric.Gauge)
	case metric.CounterType:
		s.counters[id] += *value.(*metric.Counter)
	}
}

func New(*config.Config) storage.Storage {
	return &client{
		gauges:   make(map[string]metric.Gauge),
		counters: make(map[string]metric.Counter),
	}
}
