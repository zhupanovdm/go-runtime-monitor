package metric

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/rs/zerolog"

	"github.com/zhupanovdm/go-runtime-monitor/pkg/logging"
)

var _ logging.LogCtxProvider = (*Metric)(nil)
var _ json.Marshaler = (*Metric)(nil)
var _ json.Unmarshaler = (*Metric)(nil)

type Metric struct {
	ID string
	Value
}

func (m *Metric) String() string {
	if m.Value == nil {
		return fmt.Sprintf("?/%s/?", m.ID)
	}
	return fmt.Sprintf("%s/%s/%v", m.Value.Type(), m.ID, m.Value)
}

func (m *Metric) LoggerCtx(ctx zerolog.Context) zerolog.Context {
	return logging.LogCtxUpdateWith(ctx.Str(logging.MetricIDKey, m.ID), m.Value)
}

func (m Metric) MarshalJSON() ([]byte, error) {
	if m.Value == nil {
		return nil, errors.New("metric value is not specified")
	}

	type MetricAlias Metric
	mtr := &struct {
		MetricAlias
		Type Type
	}{
		MetricAlias: MetricAlias(m),
		Type:        m.Value.Type(),
	}
	return json.Marshal(mtr)
}

func (m *Metric) UnmarshalJSON(bytes []byte) error {
	type MetricAlias Metric
	mtr := &struct {
		*MetricAlias
		Type  Type
		Value json.RawMessage
	}{
		MetricAlias: (*MetricAlias)(m),
	}

	if err := json.Unmarshal(bytes, mtr); err != nil {
		return err
	}

	switch mtr.Type {
	case GaugeType:
		v := Gauge(0)
		if err := json.Unmarshal(mtr.Value, &v); err != nil {
			return err
		}
		m.Value = &v
	case CounterType:
		v := Counter(0)
		if err := json.Unmarshal(mtr.Value, &v); err != nil {
			return err
		}
		m.Value = &v
	default:
		return fmt.Errorf("json decoder: unknown type %v", mtr.Type)
	}

	return nil
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
