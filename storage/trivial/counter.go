package trivial

import (
	"context"
	"sync"

	"github.com/zhupanovdm/go-runtime-monitor/model/metric"
	"github.com/zhupanovdm/go-runtime-monitor/pkg/logging"
	"github.com/zhupanovdm/go-runtime-monitor/storage"
)

const trivialCounterStorageName = "Trivial counter storage"

var _ storage.CounterStorage = (*trivialCounterStorage)(nil)

type trivialCounterStorage struct {
	sync.RWMutex
	data map[string]int64
}

func (s *trivialCounterStorage) Init(context.Context) error {
	return nil
}

func (s *trivialCounterStorage) Update(ctx context.Context, id string, counter metric.Counter) error {
	ctx, _ = logging.SetIfAbsentCID(ctx, logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(trivialCounterStorageName), logging.WithCID(ctx))

	inc := int64(counter)

	s.Lock()
	defer s.Unlock()
	s.data[id] += inc

	logger.Trace().Msgf("counter [%s]: increment by %d resulted in [%d]", id, inc, s.data[id])
	return nil
}

func (s *trivialCounterStorage) UpdateBulk(ctx context.Context, list metric.List) error {
	ctx, _ = logging.SetIfAbsentCID(ctx, logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(trivialCounterStorageName), logging.WithCID(ctx))

	s.Lock()
	defer s.Unlock()
	for _, m := range list {
		s.data[m.ID] += int64(*m.Value.(*metric.Counter))
	}

	logger.Trace().Msgf("counter: %d records updated", len(list))
	return nil
}

func (s *trivialCounterStorage) Get(ctx context.Context, id string) (*metric.Metric, error) {
	ctx, _ = logging.SetIfAbsentCID(ctx, logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(trivialCounterStorageName), logging.WithCID(ctx))

	s.RLock()
	defer s.RUnlock()

	value, ok := s.data[id]
	if !ok {
		logger.Trace().Msgf("counter [%s]: not found", id)
		return nil, nil
	}

	logger.Trace().Msgf("counter [%s]: restored [%d]", id, value)
	return metric.NewCounterMetric(id, metric.Counter(value)), nil
}

func (s *trivialCounterStorage) GetAll(ctx context.Context) (list metric.List, _ error) {
	ctx, _ = logging.SetIfAbsentCID(ctx, logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(trivialCounterStorageName), logging.WithCID(ctx))

	s.RLock()
	defer s.RUnlock()

	list = make([]*metric.Metric, 0, len(s.data))
	for k, v := range s.data {
		list = append(list, metric.NewCounterMetric(k, metric.Counter(v)))
	}

	logger.Trace().Msgf("counter: %d records read", len(list))
	return
}

func (s *trivialCounterStorage) Ping(context.Context) error {
	return nil
}

func (s *trivialCounterStorage) Close(context.Context) {}

func NewCounterStorage() storage.CounterStorage {
	return &trivialCounterStorage{
		data: make(map[string]int64),
	}
}
