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

func (c *client) Clear(context.Context) error {
	c.Lock()
	defer c.Unlock()
	c.gauges = make(map[string]metric.Gauge)
	c.counters = make(map[string]metric.Counter)
	return nil
}

func (c *client) IsPersistent() bool {
	return false
}

func (c *client) Init(context.Context) error {
	return nil
}

func (c *client) Update(ctx context.Context, id string, value metric.Value) error {
	ctx, _ = logging.SetIfAbsentCID(ctx, logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(trivialStorageName), logging.WithCID(ctx))

	if err := value.Type().Validate(); err != nil {
		err := fmt.Errorf("unknown metric type: %v", value.Type())
		logger.Err(err).Msg("update failed")
		return err
	}

	c.Lock()
	defer c.Unlock()
	c.save(id, value)

	logger.Trace().Msgf("metric [%s]: updated with [%v]", id, value)
	return nil
}

func (c *client) UpdateBulk(ctx context.Context, list metric.List) error {
	ctx, _ = logging.SetIfAbsentCID(ctx, logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(trivialStorageName), logging.WithCID(ctx))

	c.Lock()
	defer c.Unlock()
	for _, mtr := range list {
		if err := mtr.Type().Validate(); err != nil {
			err := fmt.Errorf("unknown metric type: %v", mtr.Type())
			logger.Err(err).Msg("update failed")
			return err
		}
		c.save(mtr.ID, mtr.Value)
	}

	logger.Trace().Msgf("%d records updated", len(list))
	return nil
}

func (c *client) Get(ctx context.Context, id string, typ metric.Type) (*metric.Metric, error) {
	ctx, _ = logging.SetIfAbsentCID(ctx, logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(trivialStorageName), logging.WithCID(ctx))

	c.RLock()
	defer c.RUnlock()

	switch typ {
	case metric.GaugeType:
		if value, ok := c.gauges[id]; ok {
			logger.Trace().Msgf("gauge [%s]: restored [%f]", id, value)
			return metric.NewGaugeMetric(id, value), nil
		}
	case metric.CounterType:
		if value, ok := c.counters[id]; ok {
			logger.Trace().Msgf("counter [%s]: restored [%d]", id, value)
			return metric.NewCounterMetric(id, value), nil
		}
	}

	logger.Trace().Msgf("counter [%s]: not found", id)
	return nil, nil
}

func (c *client) GetAll(ctx context.Context) (list metric.List, _ error) {
	ctx, _ = logging.SetIfAbsentCID(ctx, logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(trivialStorageName), logging.WithCID(ctx))

	list = make([]*metric.Metric, 0, len(c.gauges)+len(c.counters))

	c.RLock()
	defer c.RUnlock()
	for k, v := range c.gauges {
		list = append(list, metric.NewGaugeMetric(k, v))
	}
	for k, v := range c.counters {
		list = append(list, metric.NewCounterMetric(k, v))
	}

	logger.Trace().Msgf("counter: %d records read", len(list))
	return
}

func (c *client) Ping(context.Context) error {
	return nil
}

func (c *client) Close(context.Context) {}

func (c *client) save(id string, value metric.Value) {
	if m, ok := value.(*metric.Metric); ok {
		value = m.Value
	}
	switch value.Type() {
	case metric.GaugeType:
		c.gauges[id] = *value.(*metric.Gauge)
	case metric.CounterType:
		c.counters[id] += *value.(*metric.Counter)
	}
}

func New(*config.Config) storage.Storage {
	return &client{
		gauges:   make(map[string]metric.Gauge),
		counters: make(map[string]metric.Counter),
	}
}
