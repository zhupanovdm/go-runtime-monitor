package sqldb

import (
	"github.com/zhupanovdm/go-runtime-monitor/model/metric"
)

type Metrics struct {
	ID    string
	typ   string
	value float64
	delta int64
}

func (m Metrics) ToCanonical() *metric.Metric {
	switch metric.Type(m.typ) {
	case metric.GaugeType:
		return metric.NewGaugeMetric(m.ID, metric.Gauge(m.value))
	case metric.CounterType:
		return metric.NewCounterMetric(m.ID, metric.Counter(m.delta))
	}
	return nil
}

func toPrimitive(value metric.Value) (v float64, d int64) {
	if m, ok := value.(*metric.Metric); ok {
		value = m.Value
	}
	switch value.Type() {
	case metric.GaugeType:
		v = float64(*value.(*metric.Gauge))
	case metric.CounterType:
		d = int64(*value.(*metric.Counter))
	}
	return
}
