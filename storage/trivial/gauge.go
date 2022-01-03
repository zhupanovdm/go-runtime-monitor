package trivial

import (
	"context"
	"sync"

	"github.com/zhupanovdm/go-runtime-monitor/model/metric"
	"github.com/zhupanovdm/go-runtime-monitor/pkg/logging"
	"github.com/zhupanovdm/go-runtime-monitor/storage"
)

const trivialGaugeStorageName = "Trivial storage of Gauge"

var _ storage.GaugeStorage = (*trivialGaugeStorage)(nil)

type trivialGaugeStorage struct {
	sync.RWMutex
	data map[string]metric.Gauge
}

func (s *trivialGaugeStorage) Update(ctx context.Context, id string, gauge metric.Gauge) error {
	ctx, _ = logging.SetIfAbsentCID(ctx, logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(trivialGaugeStorageName), logging.WithCID(ctx))
	defer func() { logger.Info().Msg("Update query executed") }()

	s.Lock()
	defer s.Unlock()

	s.data[id] = gauge
	return nil
}

func (s *trivialGaugeStorage) Get(ctx context.Context, id string) (*metric.Metric, error) {
	ctx, _ = logging.SetIfAbsentCID(ctx, logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(trivialGaugeStorageName), logging.WithCID(ctx))
	defer func() { logger.Info().Msg("Get query executed") }()

	s.RLock()
	defer s.RUnlock()

	value, ok := s.data[id]
	if !ok {
		return nil, nil
	}
	return metric.NewGaugeMetric(id, value), nil
}

func (s *trivialGaugeStorage) GetAll(ctx context.Context) (list []*metric.Metric, _ error) {
	ctx, _ = logging.SetIfAbsentCID(ctx, logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(trivialGaugeStorageName), logging.WithCID(ctx))
	defer func() { logger.Info().Msg("Get query executed") }()

	s.RLock()
	defer s.RUnlock()

	list = make([]*metric.Metric, 0, len(s.data))
	for k, v := range s.data {
		list = append(list, metric.NewGaugeMetric(k, v))
	}
	return
}

func NewGaugeStorage() storage.GaugeStorage {
	return &trivialGaugeStorage{
		data: make(map[string]metric.Gauge),
	}
}
