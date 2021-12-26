package metric

import (
	"fmt"
	"sort"
)

const (
	GaugeType   Type = "gauge"
	CounterType Type = "counter"
)

type List []*Metric

var _ sort.Interface = (List)(nil)

func (m List) Len() int      { return len(m) }
func (m List) Swap(i, j int) { m[i], m[j] = m[j], m[i] }
func (m List) Less(i, j int) bool {
	return m[i].String() < m[j].String()
}

type Metric struct {
	ID    string
	Value Value
}

type Value interface {
	Parser
	fmt.Stringer
	Type() Type
}

type Parser interface {
	Parse(string) error
}

type Type string

var _ fmt.Stringer = (*Metric)(nil)

func (t Type) NewValue() (value Value, ok bool) {
	ok = true
	switch t {
	case GaugeType:
		value = new(Gauge)
	case CounterType:
		value = new(Counter)
	default:
		ok = false
	}
	return
}

func (m Metric) String() string {
	return fmt.Sprintf("%s/%s/%v", m.Value.Type(), m.ID, m.Value)
}

func NewGauge(id string, value float64) *Metric {
	gauge := Gauge(value)
	return &Metric{
		ID:    id,
		Value: &gauge,
	}
}

func NewGaugeFromUInt(id string, value uint64) *Metric {
	gauge := Gauge(value)
	return &Metric{
		ID:    id,
		Value: &gauge,
	}
}

func NewCounter(id string, value int64) *Metric {
	counter := Counter(value)
	return &Metric{
		ID:    id,
		Value: &counter,
	}
}
