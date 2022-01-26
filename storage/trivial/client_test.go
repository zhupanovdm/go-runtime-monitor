package trivial

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zhupanovdm/go-runtime-monitor/model/metric"
	"github.com/zhupanovdm/go-runtime-monitor/storage"
)

func Test_trivialCounterStorage_Get(t *testing.T) {
	var s storage.Storage = &client{
		counters: map[string]metric.Counter{
			"counter0": 0,
			"counter1": 1,
			"counter2": 2,
		},
	}

	tests := []struct {
		name string
		id   string
		want *metric.Metric
	}{
		{
			name: "Basic test 1",
			id:   "counter1",
			want: metric.NewCounterMetric("counter1", metric.Counter(1)),
		},
		{
			name: "Basic test 2",
			id:   "counter2",
			want: metric.NewCounterMetric("counter2", metric.Counter(2)),
		},
		{
			name: "Not found",
			id:   "absent-counter",
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.Get(context.TODO(), tt.id, metric.CounterType)
			if assert.NoError(t, err) {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_trivialCounterStorage_GetAll(t *testing.T) {
	var s storage.Storage = &client{
		counters: map[string]metric.Counter{
			"counter0": 0,
			"counter1": 1,
			"counter2": 2,
		},
	}

	tests := []struct {
		name     string
		wantList []*metric.Metric
	}{
		{
			name: "Basic test",
			wantList: []*metric.Metric{
				metric.NewCounterMetric("counter0", metric.Counter(0)),
				metric.NewCounterMetric("counter1", metric.Counter(1)),
				metric.NewCounterMetric("counter2", metric.Counter(2)),
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

func Test_trivialCounterStorage_Update(t *testing.T) {
	data := map[string]metric.Counter{
		"counter0": 0,
		"counter1": 1,
		"counter2": 2,
	}

	var s storage.Storage = &client{counters: data}

	tests := []struct {
		name   string
		metric *metric.Metric
		want   *metric.Metric
	}{
		{
			name:   "New metric",
			metric: metric.NewCounterMetric("foo-counter", 9),
			want:   metric.NewCounterMetric("foo-counter", 9),
		},
		{
			name:   "Update existing metric",
			metric: metric.NewCounterMetric("counter2", 3),
			want:   metric.NewCounterMetric("counter2", 5),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s.Update(context.TODO(), tt.metric)
			if assert.NoError(t, err) {
				m, _ := s.Get(context.TODO(), tt.metric.ID, tt.metric.Type())
				assert.Equal(t, tt.want, m)
			}
		})
	}
}

func Test_trivialGaugeStorage_Get(t *testing.T) {
	var s storage.Storage = &client{
		gauges: map[string]metric.Gauge{
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
			got, err := s.Get(context.TODO(), tt.id, metric.GaugeType)
			if assert.NoError(t, err) {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_trivialGaugeStorage_GetAll(t *testing.T) {
	var s storage.Storage = &client{
		gauges: map[string]metric.Gauge{
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
	data := map[string]metric.Gauge{
		"gauge0": 0,
		"gauge1": .1,
		"gauge2": .2,
	}

	var s storage.Storage = &client{gauges: data}

	tests := []struct {
		name   string
		metric *metric.Metric
		want   *metric.Metric
	}{
		{
			name:   "New metric",
			metric: metric.NewGaugeMetric("foo-gauge", .9),
			want:   metric.NewGaugeMetric("foo-gauge", .9),
		},
		{
			name:   "Update existing metric",
			metric: metric.NewGaugeMetric("gauge2", 1),
			want:   metric.NewGaugeMetric("gauge2", 1),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s.Update(context.TODO(), tt.metric)
			if assert.NoError(t, err) {
				m, _ := s.Get(context.TODO(), tt.metric.ID, tt.metric.Type())
				assert.Equal(t, tt.want, m)
			}
		})
	}
}
