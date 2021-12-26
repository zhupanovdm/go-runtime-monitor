package metric

import (
	"fmt"
	"strconv"
)

type Gauge float64

var _ Value = (*Gauge)(nil)

func (g Gauge) String() string {
	return fmt.Sprintf("%.3f", g)
}

func (g *Gauge) Parse(s string) error {
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return fmt.Errorf("can't parse gauge from '%s': %v", s, err)
	}
	*g = Gauge(val)
	return nil
}

func (Gauge) Type() Type {
	return GaugeType
}
