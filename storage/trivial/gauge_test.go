package trivial

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zhupanovdm/go-runtime-monitor/model/metric"
	"github.com/zhupanovdm/go-runtime-monitor/storage"
)

func Test_trivialGaugeStorage_Get(t *testing.T) {
	var s storage.GaugeStorage = &trivialGaugeStorage{
		data: map[string]float64{
			"gauge0": 0,
			"gauge1": .1,
			"gauge2": .2,
		},
	}

	tests := []struct {
		name string
		id   string
		want *metric.Metric
	}{
		{
			name: "Basic test 1",
			id:   "gauge1",
			want: metric.NewGaugeMetric("gauge1", metric.Gauge(.1)),
		},
		{
			name: "Basic test 2",
			id:   "gauge2",
			want: metric.NewGaugeMetric("gauge2", metric.Gauge(.2)),
		},
		{
			name: "Not found",
			id:   "absent-gauge",
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.Get(context.TODO(), tt.id)
			if assert.NoError(t, err) {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_trivialGaugeStorage_GetAll(t *testing.T) {
	var s storage.GaugeStorage = &trivialGaugeStorage{
		data: map[string]float64{
			"gauge0": 0,
			"gauge1": .1,
			"gauge2": .2,
		},
	}

	tests := []struct {
		name     string
		wantList []*metric.Metric
	}{
		{
			name: "Basic test",
			wantList: []*metric.Metric{
				metric.NewGaugeMetric("gauge0", metric.Gauge(0)),
				metric.NewGaugeMetric("gauge1", metric.Gauge(.1)),
				metric.NewGaugeMetric("gauge2", metric.Gauge(.2)),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotList, err := s.GetAll(context.TODO())
			if assert.NoError(t, err) {
				assert.ElementsMatch(t, tt.wantList, gotList)
			}
		})
	}
}

func Test_trivialGaugeStorage_Update(t *testing.T) {
	data := map[string]float64{
		"gauge0": 0,
		"gauge1": .1,
		"gauge2": .2,
	}

	var s storage.GaugeStorage = &trivialGaugeStorage{data: data}

	tests := []struct {
		name  string
		id    string
		value metric.Gauge
		want  metric.Gauge
	}{
		{
			name:  "New value",
			id:    "foo-gauge",
			value: metric.Gauge(.9),
			want:  metric.Gauge(.9),
		},
		{
			name:  "Update existing value",
			id:    "gauge2",
			value: metric.Gauge(1),
			want:  metric.Gauge(1),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s.Update(context.TODO(), tt.id, tt.value)
			if assert.NoError(t, err) {
				m, _ := s.Get(context.TODO(), tt.id)
				assert.Equal(t, tt.want, *(m.Value.(*metric.Gauge)))
			}
		})
	}
}
