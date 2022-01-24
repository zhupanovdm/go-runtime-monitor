package trivial

import (
	"context"
	"sync"

	"github.com/zhupanovdm/go-runtime-monitor/model/metric"
	"github.com/zhupanovdm/go-runtime-monitor/pkg/logging"
	"github.com/zhupanovdm/go-runtime-monitor/storage"
)

const trivialGaugeStorageName = "Trivial gauge storage"

var _ storage.GaugeStorage = (*trivialGaugeStorage)(nil)

type trivialGaugeStorage struct {
	sync.RWMutex
	data map[string]float64
}

func (s *trivialGaugeStorage) Init(context.Context) error {
	return nil
}

func (s *trivialGaugeStorage) Update(ctx context.Context, id string, gauge metric.Gauge) error {
	ctx, _ = logging.SetIfAbsentCID(ctx, logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(trivialGaugeStorageName), logging.WithCID(ctx))

	value := float64(gauge)

	s.Lock()
	defer s.Unlock()
	s.data[id] = value

	logger.Trace().Msgf("gauge [%s]: updated with [%f]", id, value)
	return nil
}

func (s *trivialGaugeStorage) UpdateBulk(ctx context.Context, list metric.List) error {
	ctx, _ = logging.SetIfAbsentCID(ctx, logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(trivialCounterStorageName), logging.WithCID(ctx))

	s.Lock()
	defer s.Unlock()
	for _, m := range list {
		s.data[m.ID] = float64(*m.Value.(*metric.Gauge))
	}

	logger.Trace().Msgf("gauge: %d records updated", len(list))
	return nil
}

func (s *trivialGaugeStorage) Get(ctx context.Context, id string) (*metric.Metric, error) {
	ctx, _ = logging.SetIfAbsentCID(ctx, logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(trivialGaugeStorageName), logging.WithCID(ctx))

	s.RLock()
	defer s.RUnlock()

	value, ok := s.data[id]
	if !ok {
		logger.Trace().Msgf("gauge [%s]: not found", id)
		return nil, nil
	}

	logger.Trace().Msgf("gauge [%s]: restored [%f]", id, value)
	return metric.NewGaugeMetric(id, metric.Gauge(value)), nil
}

func (s *trivialGaugeStorage) GetAll(ctx context.Context) (list metric.List, _ error) {
	ctx, _ = logging.SetIfAbsentCID(ctx, logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(trivialGaugeStorageName), logging.WithCID(ctx))

	s.RLock()
	defer s.RUnlock()

	list = make([]*metric.Metric, 0, len(s.data))
	for k, v := range s.data {
		list = append(list, metric.NewGaugeMetric(k, metric.Gauge(v)))
	}

	logger.Trace().Msgf("gauge: %d records read", len(list))
	return
}

func (s *trivialGaugeStorage) Ping(context.Context) error {
	return nil
}

func (s *trivialGaugeStorage) Close(context.Context) {}

func NewGaugeStorage() storage.GaugeStorage {
	return &trivialGaugeStorage{
		data: make(map[string]float64),
	}
}
