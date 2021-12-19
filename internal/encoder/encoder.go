package encoder

type Encoder interface {
	Encode() string
	Decode(string) error
	Type() string
}

const (
	CounterType = "counter"
	GaugeType   = "gauge"
)

func NewGaugeI(name string, value uint64) Encoder {
	v := gauge(value)
	return &metric{
		name:  name,
		value: &v,
	}
}

func NewGaugeF(name string, value float64) Encoder {
	v := gauge(value)
	return &metric{
		name:  name,
		value: &v,
	}
}

func NewCounter(name string, value int64) Encoder {
	v := counter(value)
	return &metric{
		name:  name,
		value: &v,
	}
}
