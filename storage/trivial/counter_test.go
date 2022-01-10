package trivial

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zhupanovdm/go-runtime-monitor/model/metric"
	"github.com/zhupanovdm/go-runtime-monitor/storage"
)

func Test_trivialCounterStorage_Get(t *testing.T) {
	var s storage.CounterStorage = &trivialCounterStorage{
		data: map[string]int64{
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
			got, err := s.Get(context.TODO(), tt.id)
			if assert.NoError(t, err) {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_trivialCounterStorage_GetAll(t *testing.T) {
	var s storage.CounterStorage = &trivialCounterStorage{
		data: map[string]int64{
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
	data := map[string]int64{
		"counter0": 0,
		"counter1": 1,
		"counter2": 2,
	}

	var s storage.CounterStorage = &trivialCounterStorage{data: data}

	tests := []struct {
		name  string
		id    string
		value metric.Counter
		want  metric.Counter
	}{
		{
			name:  "New value",
			id:    "foo-counter",
			value: metric.Counter(9),
			want:  metric.Counter(9),
		},
		{
			name:  "Update existing value",
			id:    "counter2",
			value: metric.Counter(3),
			want:  metric.Counter(5),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s.Update(context.TODO(), tt.id, tt.value)
			if assert.NoError(t, err) {
				m, _ := s.Get(context.TODO(), tt.id)
				assert.Equal(t, tt.want, *(m.Value.(*metric.Counter)))
			}
		})
	}
}
