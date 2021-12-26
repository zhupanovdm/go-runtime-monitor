package service

import (
	"fmt"
	"sync"

	"github.com/zhupanovdm/go-runtime-monitor/cmd/server/repo"
	"github.com/zhupanovdm/go-runtime-monitor/internal/model/metric"
)

type Metrics interface {
	Save(metric metric.Metric) error
	Get(id string, typ metric.Type) (metric.Value, error)
	GetAll() (metric.List, error)
}

var _ Metrics = (*metricsService)(nil)

type metricsService struct {
	gaugesMutex sync.RWMutex
	gauges      repo.GaugeRepo

	countersMutex sync.RWMutex
	counters      repo.CounterRepo
}

func (s *metricsService) Save(m metric.Metric) (err error) {
	switch value := m.Value.(type) {
	case *metric.Gauge:
		s.gaugesMutex.Lock()
		defer s.gaugesMutex.Unlock()
		err = s.gauges.Save(m.Id, *value)

	case *metric.Counter:
		s.countersMutex.Lock()
		defer s.countersMutex.Unlock()
		counter, ok := s.counters.Get(m.Id, &err)
		if ok {
			*value += counter
		}
		if err != nil {
			return
		}
		err = s.counters.Save(m.Id, *value)

	default:
		err = fmt.Errorf("type is not supported yet: %T", value)

	}
	return
}

func (s *metricsService) Get(id string, typ metric.Type) (value metric.Value, err error) {
	switch typ {
	case metric.GaugeType:
		s.gaugesMutex.RLock()
		defer s.gaugesMutex.RUnlock()
		if gauge, ok := s.gauges.Get(id, &err); ok {
			value = &gauge
		}

	case metric.CounterType:
		s.countersMutex.RLock()
		defer s.countersMutex.RUnlock()
		if counter, ok := s.counters.Get(id, &err); ok {
			value = &counter
		}

	default:
		err = fmt.Errorf("unknown metric type: %v", typ)
	}
	return
}

func (s *metricsService) GetAll() (list metric.List, err error) {
	var counters metric.List

	withLock(s.gaugesMutex.RLocker(), func() {
		list, err = s.gauges.GetAll()
		if err != nil {
			return
		}
	})

	withLock(s.countersMutex.RLocker(), func() {
		counters, err = s.counters.GetAll()
		if err != nil {
			return
		}

		list = append(list, counters...)
	})

	return
}

func NewMetrics(gauges repo.GaugeRepo, counters repo.CounterRepo) Metrics {
	return &metricsService{
		gauges:   gauges,
		counters: counters,
	}
}

func withLock(lock sync.Locker, do func()) {
	lock.Lock()
	defer lock.Unlock()
	do()
}
