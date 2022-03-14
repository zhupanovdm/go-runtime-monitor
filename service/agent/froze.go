package agent

import (
	"sync"

	"github.com/zhupanovdm/go-runtime-monitor/model/metric"
)

type Froze struct {
	sync.Mutex
	gauges   map[string]float64
	counters map[string]int64
}

func (f *Froze) UpdateGauge(id string, gauge float64) {
	f.gauges[id] = gauge
}

func (f *Froze) UpdateCounter(id string, counter int64) {
	f.counters[id] += counter
}

func (f *Froze) List() metric.List {
	list := make(metric.List, 0, len(f.gauges)+len(f.counters))
	for id, gauge := range f.gauges {
		list = append(list, metric.NewGaugeMetric(id, metric.Gauge(gauge)))
	}
	for id, counter := range f.counters {
		list = append(list, metric.NewCounterMetric(id, metric.Counter(counter)))
	}
	return list
}

func NewFroze() *Froze {
	return &Froze{
		gauges:   make(map[string]float64),
		counters: make(map[string]int64),
	}
}
