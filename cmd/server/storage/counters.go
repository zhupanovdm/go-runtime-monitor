package storage

import "github.com/zhupanovdm/go-runtime-monitor/internal/measure"

var _ CountersRepository = (*countersStorage)(nil)

type countersStorage struct {
	data map[string]int64
}

func NewCounters() CountersRepository {
	return &countersStorage{
		data: make(map[string]int64),
	}
}

func (r *countersStorage) Save(key string, c measure.Counter) error {
	r.data[key] += int64(c)
	return nil
}

func (r *countersStorage) Get(key string) (d int64, ok bool) {
	d, ok = r.data[key]
	return
}
