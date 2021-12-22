package measure

import (
	"fmt"
)

type Type string

type Value interface {
	Type() Type
	Encode() string
	Decode(string) error
}

const (
	MetricType  Type = "metric"
	GaugeType   Type = "gauge"
	CounterType Type = "counter"
)

func (t Type) New() (Value, error) {
	return New(t)
}

func New(t Type) (Value, error) {
	switch t {
	case GaugeType:
		return new(Gauge), nil
	case CounterType:
		return new(Counter), nil
	case MetricType:
		return &Metric{}, nil
	}
	return nil, fmt.Errorf("type %v is not implemented", t)
}
