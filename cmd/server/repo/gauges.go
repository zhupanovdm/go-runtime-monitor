package repo

import (
	"github.com/zhupanovdm/go-runtime-monitor/internal/model/metric"
)

type GaugeRepo interface {
	Save(string, metric.Gauge) error
	Get(string, *error) (metric.Gauge, bool)
	GetAll() (metric.List, error)
}

var _ GaugeRepo = (*gaugeStorage)(nil)

type gaugeStorage struct {
	data map[string]metric.Gauge
}

func (s *gaugeStorage) Save(id string, gauge metric.Gauge) error {
	s.data[id] = gauge
	return nil
}

func (s *gaugeStorage) Get(id string, _ *error) (gauge metric.Gauge, ok bool) {
	gauge, ok = s.data[id]
	return
}

func (s *gaugeStorage) GetAll() (list metric.List, _ error) {
	list = make(metric.List, 0, len(s.data))
	for id, gauge := range s.data {
		value := gauge
		list = append(list, &metric.Metric{
			Id:    id,
			Value: &value,
		})
	}
	return
}

func Gauges() GaugeRepo {
	return &gaugeStorage{
		data: make(map[string]metric.Gauge),
	}
}
