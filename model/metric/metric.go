package metric

import (
	"fmt"
	"sort"
)

type Metric struct {
	ID string
	Value
}

var _ sort.Interface = (ByString)(nil)

type ByString []*Metric

func (m ByString) Len() int           { return len(m) }
func (m ByString) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }
func (m ByString) Less(i, j int) bool { return m[i].String() < m[j].String() }

func (m Metric) String() string {
	return fmt.Sprintf("%s/%s/%v", m.Value.Type(), m.ID, m.Value)
}

func NewGaugeMetric(id string, gauge Gauge) *Metric {
	return &Metric{
		ID:    id,
		Value: &gauge,
	}
}

func NewCounterMetric(id string, counter Counter) *Metric {
	return &Metric{
		ID:    id,
		Value: &counter,
	}
}
