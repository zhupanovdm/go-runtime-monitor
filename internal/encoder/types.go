package encoder

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type counter int64
type gauge float64
type metric struct {
	name  string
	value Encoder
}

var _ Encoder = (*counter)(nil)
var _ Encoder = (*gauge)(nil)
var _ Encoder = (*metric)(nil)

func (c *counter) Encode() string {
	return fmt.Sprintf("%d", *c)
}

func (c *counter) Decode(raw string) error {
	val, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return err
	}
	*c = counter(val)
	return nil
}

func (*counter) Type() string {
	return CounterType
}

func (g *gauge) Encode() string {
	return fmt.Sprintf("%f", float64(*g))
}

func (g *gauge) Decode(s string) error {
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}
	*g = gauge(val)
	return nil
}

func (*gauge) Type() string {
	return GaugeType
}

func (m *metric) Encode() string {
	return fmt.Sprintf("%s/%s/%s", m.value.Type(), m.name, m.value.Encode())
}

func (m *metric) Decode(raw string) error {
	chunks := strings.Split(raw, "/")
	if len(chunks) != 3 {
		return errors.New("expected 3 parts of metrics")
	}

	m.name = chunks[1]

	var val Encoder
	switch chunks[0] {
	case CounterType:
		val = new(counter)
	case GaugeType:
		val = new(gauge)
	default:
		return errors.New("unknown type")
	}

	err := val.Decode(chunks[3])
	if err != nil {
		return errors.New("can not read value")
	}

	return nil
}

func (m *metric) Type() string {
	return m.name
}
