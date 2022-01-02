package trivial

import (
	"context"
	"sync"

	"github.com/zhupanovdm/go-runtime-monitor/model/metric"
	"github.com/zhupanovdm/go-runtime-monitor/storage"
)

var _ storage.GaugeStorage = (*trivialGaugeStorage)(nil)

type trivialGaugeStorage struct {
	sync.RWMutex
	data map[string]metric.Gauge
}

func (s *trivialGaugeStorage) Update(_ context.Context, id string, gauge metric.Gauge) error {
	s.Lock()
	defer s.Unlock()

	s.data[id] = gauge
	return nil
}

func (s *trivialGaugeStorage) Get(_ context.Context, id string) (*metric.Metric, error) {
	s.RLock()
	defer s.RUnlock()

	value, ok := s.data[id]
	if !ok {
		return nil, nil
	}
	return metric.NewGaugeMetric(id, value), nil
}

func (s *trivialGaugeStorage) GetAll(_ context.Context) (list []*metric.Metric, _ error) {
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
