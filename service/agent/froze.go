package agent

import (
	"sync"

	"github.com/zhupanovdm/go-runtime-monitor/model/metric"
)

// Froze object is used to hold collected metrics. Another point is to synchronize read and write access of several
// application components to operate consistent metrics data.
type Froze struct {
	sync.Mutex
	gauges   map[string]float64
	counters map[string]int64
}

// UpdateGauge updates single gauge metrics measure. Update will override previously stored value. Thread unsafe, should
// be locked before update.
func (f *Froze) UpdateGauge(id string, gauge float64) {
	f.gauges[id] = gauge
}

// UpdateCounter updates single gauge metrics measure. Update will increment previously stored value. Thread unsafe, should
// be locked before update.
func (f *Froze) UpdateCounter(id string, counter int64) {
	f.counters[id] += counter
}

// List entirely reads metrics measures copy into list. Thread unsafe, should be locked before read.
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

// NewFroze creates new Froze object.
func NewFroze() *Froze {
	return &Froze{
		gauges:   make(map[string]float64),
		counters: make(map[string]int64),
	}
}
