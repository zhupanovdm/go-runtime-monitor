package sqldb

import (
	"github.com/zhupanovdm/go-runtime-monitor/model/metric"
)

type Scanner interface {
	Scan(...interface{}) error
}

type Metrics struct {
	ID    string
	typ   string
	value float64
	delta int64
}

func (m *Metrics) Scan(scanner Scanner) error {
	if err := scanner.Scan(&m.ID, &m.typ, &m.value, &m.delta); err != nil {
		return err
	}
	return nil
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

func FromCanonical(mtr *metric.Metric) *Metrics {
	switch mtr.Type() {
	case metric.GaugeType:
		return &Metrics{
			ID:    mtr.ID,
			typ:   string(mtr.Type()),
			value: float64(*mtr.Value.(*metric.Gauge)),
		}
	case metric.CounterType:
		return &Metrics{
			ID:    mtr.ID,
			typ:   string(mtr.Type()),
			delta: int64(*mtr.Value.(*metric.Counter)),
		}
	}
	return nil
}
