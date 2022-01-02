package trivial

import (
	"context"
	"sync"

	"github.com/zhupanovdm/go-runtime-monitor/model/metric"
	"github.com/zhupanovdm/go-runtime-monitor/storage"
)

var _ storage.CounterStorage = (*trivialCounterStorage)(nil)

type trivialCounterStorage struct {
	sync.RWMutex
	data map[string]metric.Counter
}

func (s *trivialCounterStorage) Update(_ context.Context, id string, counter metric.Counter) error {
	s.Lock()
	defer s.Unlock()

	s.data[id] += counter
	return nil
}

func (s *trivialCounterStorage) Get(_ context.Context, id string) (*metric.Metric, error) {
	s.RLock()
	defer s.RUnlock()

	value, ok := s.data[id]
	if !ok {
		return nil, nil
	}
	return metric.NewCounterMetric(id, value), nil
}

func (s *trivialCounterStorage) GetAll(_ context.Context) (list []*metric.Metric, _ error) {
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
