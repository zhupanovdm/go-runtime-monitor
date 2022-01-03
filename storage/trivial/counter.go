package trivial

import (
	"context"
	"sync"

	"github.com/zhupanovdm/go-runtime-monitor/model/metric"
	"github.com/zhupanovdm/go-runtime-monitor/pkg/logging"
	"github.com/zhupanovdm/go-runtime-monitor/storage"
)

const trivialCounterStorageName = "Trivial storage of Counter"

var _ storage.CounterStorage = (*trivialCounterStorage)(nil)

type trivialCounterStorage struct {
	sync.RWMutex
	data map[string]metric.Counter
}

func (s *trivialCounterStorage) Update(ctx context.Context, id string, counter metric.Counter) error {
	ctx, _ = logging.SetIfAbsentCID(ctx, logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(trivialCounterStorageName), logging.WithCID(ctx))
	defer func() { logger.Info().Msg("Update query executed") }()

	s.Lock()
	defer s.Unlock()

	s.data[id] += counter
	return nil
}

func (s *trivialCounterStorage) Get(ctx context.Context, id string) (*metric.Metric, error) {
	ctx, _ = logging.SetIfAbsentCID(ctx, logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(trivialCounterStorageName), logging.WithCID(ctx))
	defer func() { logger.Info().Msg("Get query executed") }()

	s.RLock()
	defer s.RUnlock()

	value, ok := s.data[id]
	if !ok {
		return nil, nil
	}
	return metric.NewCounterMetric(id, value), nil
}

func (s *trivialCounterStorage) GetAll(ctx context.Context) (list []*metric.Metric, _ error) {
	ctx, _ = logging.SetIfAbsentCID(ctx, logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(trivialCounterStorageName), logging.WithCID(ctx))
	defer func() { logger.Info().Msg("Get all query executed") }()

	s.RLock()
	defer s.RUnlock()

	list = make([]*metric.Metric, 0, len(s.data))
	for k, v := range s.data {
		list = append(list, metric.NewCounterMetric(k, v))
	}
	return
}

func NewCounterStorage() storage.CounterStorage {
	return &trivialCounterStorage{
		data: make(map[string]metric.Counter),
	}
}
