package model

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/zhupanovdm/go-runtime-monitor/model/metric"
)

var _ json.Unmarshaler = (*Metrics)(nil)

type Metrics struct {
	ID    string   `json:"id"`              // metric name
	MType string   `json:"type"`            // metric type is enum value {"counter", "gauge"}
	Delta *int64   `json:"delta,omitempty"` // metric measure if MType is "counter"
	Value *float64 `json:"value,omitempty"` // metric measure if MType is "gauge"
	Hash  string   `json:"hash,omitempty"`  // packet hash sum
}

func (m *Metrics) Sign(key string) error {
	hash, err := m.calcHash([]byte(key))
	if err != nil {
		return err
	}
	m.Hash = hex.EncodeToString(hash)
	return nil
}

func (m Metrics) Verify(key string) error {
	mac1, err := hex.DecodeString(m.Hash)
	if err != nil {
		return fmt.Errorf("can't decode packet hash: %w", err)
	}
	mac2, err := m.calcHash([]byte(key))
	if err != nil {
		return fmt.Errorf("unable to recalculate packet hash: %w", err)
	}
	if !hmac.Equal(mac1, mac2) {
		return fmt.Errorf("sign verification failed")
	}
	return nil
}

func (m Metrics) String() string {
	switch metric.Type(m.MType) {
	case metric.GaugeType:
		return fmt.Sprintf("%s/%s/%v", m.ID, m.MType, m.Value)
	case metric.CounterType:
		return fmt.Sprintf("%s/%s/%v", m.ID, m.MType, m.Delta)
	}
	return fmt.Sprintf("unknown:%s/%s", m.ID, m.MType)
}

func (m *Metrics) UnmarshalJSON(bytes []byte) error {
	type MetricsAlias Metrics
	mtr := &struct {
		*MetricsAlias
		Delta json.RawMessage
		Value json.RawMessage
	}{
		MetricsAlias: (*MetricsAlias)(m),
	}
	if err := json.Unmarshal(bytes, mtr); err != nil {
		return err
	}

	switch metric.Type(mtr.MType) {
	case metric.GaugeType:
		if mtr.Value != nil {
			m.Value = new(float64)
			if err := json.Unmarshal(mtr.Value, m.Value); err != nil {
				return fmt.Errorf("value of type (%T): unmarshal: %w", m.Value, err)
			}
		}
	case metric.CounterType:
		if mtr.Delta != nil {
			m.Delta = new(int64)
			if err := json.Unmarshal(mtr.Delta, m.Delta); err != nil {
				return fmt.Errorf("value of type (%T): unmarshal: %w", m.Value, err)
			}
		}
	}
	return nil
}

func (m Metrics) ToCanonical() *metric.Metric {
	switch metric.Type(m.MType) {
	case metric.GaugeType:
		if m.Value == nil {
			return metric.NewGaugeMetric(m.ID, metric.Gauge(0))
		}
		return metric.NewGaugeMetric(m.ID, metric.Gauge(*m.Value))
	case metric.CounterType:
		if m.Delta == nil {
			return metric.NewCounterMetric(m.ID, metric.Counter(0))
		}
		return metric.NewCounterMetric(m.ID, metric.Counter(*m.Delta))
	}
	return nil
}

func NewFromCanonical(mtr *metric.Metric) *Metrics {
	m := &Metrics{
		ID:    mtr.ID,
		MType: string(mtr.Type()),
	}

	switch mtr.Type() {
	case metric.GaugeType:
		if gauge, ok := mtr.Value.(*metric.Gauge); ok {
			value := float64(*gauge)
			m.Value = &value
			return m
		}
	case metric.CounterType:
		if counter, ok := mtr.Value.(*metric.Counter); ok {
			delta := int64(*counter)
			m.Delta = &delta
			return m
		}
	}
	return nil
}

func (m Metrics) calcHash(key []byte) ([]byte, error) {
	var data string
	switch metric.Type(m.MType) {
	case metric.CounterType:
		data = fmt.Sprintf("%s:counter:%d", m.ID, *m.Delta)
	case metric.GaugeType:
		data = fmt.Sprintf("%s:gauge:%f", m.ID, *m.Value)
	default:
		return nil, fmt.Errorf("hash calc: can't calc for unknown type %s", m.MType)
	}

	h := hmac.New(sha256.New, key)
	h.Write([]byte(data))
	return h.Sum(nil), nil
}

func (m *Metrics) Validate(validators ...func(*Metrics) error) error {
	for _, validator := range validators {
		if err := validator(m); err != nil {
			return fmt.Errorf("validator: %w", err)
		}
	}
	return nil
}

func CheckValue(m *Metrics) error {
	switch metric.Type(m.MType) {
	case metric.GaugeType:
		if m.Value == nil {
			return fmt.Errorf("metrics validate: for type %v [Value] must not be empty", m.MType)
		}
	case metric.CounterType:
		if m.Delta == nil {
			return fmt.Errorf("metrics validate: for type %v [Delta] must not be empty", m.MType)
		}
	}
	return nil
}

func CheckID(m *Metrics) error {
	if len(m.ID) == 0 {
		return errors.New("metrics validate: metric ID is empty")
	}
	return nil
}

func CheckType(m *Metrics) error {
	if err := metric.Type(m.MType).Validate(); err != nil {
		return fmt.Errorf("metrics validate: type assertion: %w", err)
	}
	return nil
}
