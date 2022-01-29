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

func (c *client) Clear(ctx context.Context) error {
	ctx, _ = logging.SetIfAbsentCID(ctx, logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(trivialStorageName), logging.WithCID(ctx))

	c.Lock()
	defer c.Unlock()
	c.gauges = make(map[string]metric.Gauge)
	c.counters = make(map[string]metric.Counter)

	logger.Info().Msg("cleared")
	return nil
}

func (c *client) IsPersistent() bool {
	return false
}

func (c *client) Init(ctx context.Context) error {
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(trivialStorageName))
	logger.Info().Msg("initialized")
	return nil
}

func (c *client) Update(ctx context.Context, mtr *metric.Metric) error {
	ctx, _ = logging.SetIfAbsentCID(ctx, logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(trivialStorageName), logging.WithCID(ctx))

	c.Lock()
	defer c.Unlock()
	return c.update(logging.SetLogger(ctx, logger), mtr)
}

func (c *client) UpdateBulk(ctx context.Context, list metric.List) error {
	ctx, _ = logging.SetIfAbsentCID(ctx, logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(trivialStorageName), logging.WithCID(ctx))

	c.Lock()
	defer c.Unlock()
	for _, mtr := range list {
		if err := c.update(logging.SetLogger(ctx, logger), mtr); err != nil {
			return err
		}
	}
	logger.Trace().Msgf("%d records updated", len(list))
	return nil
}

func (c *client) Get(ctx context.Context, id string, typ metric.Type) (*metric.Metric, error) {
	ctx, _ = logging.SetIfAbsentCID(ctx, logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(trivialStorageName), logging.WithCID(ctx))
	logger.UpdateContext(logging.LogCtxKeyStr(logging.MetricIDKey, id))
	logger.UpdateContext(logging.LogCtxFrom(typ))

	c.RLock()
	defer c.RUnlock()
	switch typ {
	case metric.GaugeType:
		if value, ok := c.gauges[id]; ok {
			mtr := metric.NewGaugeMetric(id, value)
			logger.UpdateContext(logging.LogCtxFrom(mtr))
			logger.Trace().Msg("read")
			return mtr, nil
		}
	case metric.CounterType:
		if delta, ok := c.counters[id]; ok {
			mtr := metric.NewCounterMetric(id, delta)
			logger.UpdateContext(logging.LogCtxFrom(mtr))
			logger.Trace().Msg("read")
			return mtr, nil
		}
	default:
		err := fmt.Errorf("unknown metric %v", typ)
		logger.Err(err).Msg("read failed")
		return nil, err
	}

	logger.Trace().Msg("not found")
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

	logger.Trace().Msgf("%d records read", len(list))
	return
}

func (c *client) Ping(ctx context.Context) error {
	ctx, _ = logging.SetIfAbsentCID(ctx, logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(trivialStorageName), logging.WithCID(ctx))
	logger.Trace().Msg("storage is online")
	return nil
}

func (c *client) Close(ctx context.Context) {
	ctx, _ = logging.SetIfAbsentCID(ctx, logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(trivialStorageName), logging.WithCID(ctx))
	logger.Info().Msg("closed")
}

func (c *client) update(ctx context.Context, mtr *metric.Metric) error {
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(trivialStorageName), logging.WithCID(ctx))
	logger.UpdateContext(logging.LogCtxFrom(mtr))

	if err := mtr.Type().Validate(); err != nil {
		logger.Err(err).Msg("update failed")
		return err
	}

	switch mtr.Type() {
	case metric.GaugeType:
		c.gauges[mtr.ID] = *mtr.Value.(*metric.Gauge)
	case metric.CounterType:
		c.counters[mtr.ID] += *mtr.Value.(*metric.Counter)
	default:
		err := fmt.Errorf("unknown metric %v", mtr.Type())
		logger.Err(err).Msg("update failed")
		return err
	}

	logger.Trace().Msg("updated")
	return nil
}

func New(*config.Config) storage.Storage {
	return &client{
		gauges:   make(map[string]metric.Gauge),
		counters: make(map[string]metric.Counter),
	}
}
