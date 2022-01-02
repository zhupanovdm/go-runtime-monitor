package metric

import (
	"fmt"
)

const (
	GaugeType   Type = "gauge"
	CounterType Type = "counter"
)

type Type string

func (t Type) String() string {
	return string(t)
}

func (t Type) Validate() error {
	switch t {
	case GaugeType, CounterType:
		return nil
	default:
		return fmt.Errorf("unkown metric type: %v", t)
	}
}

func (t Type) New() (value Value, err error) {
	switch t {
	case GaugeType:
		value = new(Gauge)
	case CounterType:
		value = new(Counter)
	default:
		err = fmt.Errorf("type %v creation is not supported", t)
	}
	return
}

func (t Type) Parse(s string) (value Value, err error) {
	if value, err = t.New(); err != nil {
		return
	}
	if err = value.Parse(s); err != nil {
		value = nil
	}
	return
}
