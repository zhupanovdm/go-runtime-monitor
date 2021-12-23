package storage

import "github.com/zhupanovdm/go-runtime-monitor/internal/measure"

var _ GaugesRepository = (*gaugesStorage)(nil)

type gaugesStorage struct {
	data map[string]float64
}

func NewGauges() GaugesRepository {
	return &gaugesStorage{
		data: make(map[string]float64),
	}
}

func (s *gaugesStorage) Save(key string, g measure.Gauge) error {
	s.data[key] = float64(g)
	return nil
}

func (s *gaugesStorage) Get(key string) (d float64, ok bool) {
	d, ok = s.data[key]
	return
}
