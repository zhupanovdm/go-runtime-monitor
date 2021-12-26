package repo

import (
	"github.com/zhupanovdm/go-runtime-monitor/internal/model/metric"
)

type CounterRepo interface {
	Save(string, metric.Counter) error
	Get(string, *error) (metric.Counter, bool)
	GetAll() (metric.List, error)
}

var _ CounterRepo = (*counterStorage)(nil)

type counterStorage struct {
	data map[string]metric.Counter
}

func (r *counterStorage) Save(id string, counter metric.Counter) error {
	r.data[id] = counter
	return nil
}

func (r *counterStorage) Get(id string, _ *error) (counter metric.Counter, ok bool) {
	counter, ok = r.data[id]
	return
}

func (r *counterStorage) GetAll() (list metric.List, _ error) {
	list = make(metric.List, 0, len(r.data))
	for id, counter := range r.data {
		value := counter
		list = append(list, &metric.Metric{
			ID:    id,
			Value: &value,
		})
	}
	return
}

func Counters() CounterRepo {
	return &counterStorage{
		data: make(map[string]metric.Counter),
	}
}
