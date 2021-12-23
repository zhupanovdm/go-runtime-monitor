package storage

import "github.com/zhupanovdm/go-runtime-monitor/internal/measure"

type GaugesRepository interface {
	Save(string, measure.Gauge) error
	Get(string) (float64, bool)
}

type CountersRepository interface {
	Save(string, measure.Counter) error
	Get(string) (int64, bool)
}
