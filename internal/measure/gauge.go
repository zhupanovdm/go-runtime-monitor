package measure

import (
	"fmt"
	"strconv"
)

type Gauge float64

var _ Value = (*Gauge)(nil)

func (g Gauge) Encode() string {
	return fmt.Sprintf("%f", g)
}

func (g *Gauge) Decode(s string) error {
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}
	*g = Gauge(val)
	return nil
}

func (Gauge) Type() Type {
	return GaugeType
}
