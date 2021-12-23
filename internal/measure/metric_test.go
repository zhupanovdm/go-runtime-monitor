package measure

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMetricDecode(t *testing.T) {
	g := Gauge(0.000001)
	c := Counter(1)

	tests := []struct {
		name    string
		want    *Metric
		sample  string
		wantErr bool
	}{
		{
			name:   "Gauge test",
			sample: "gauge/foo/0.000001",
			want: &Metric{
				Name:  "foo",
				Value: Value(&g),
			},
		},
		{
			name:   "Counter test",
			sample: "counter/bar/1",
			want: &Metric{
				Name:  "bar",
				Value: Value(&c),
			},
		},
		{
			name:    "Unknown type",
			sample:  "foo/bar/1",
			wantErr: true,
		},
		{
			name:    "Type mismatch",
			sample:  "counter/bar/1.1",
			wantErr: true,
		},
		{
			name:    "Omit name element",
			sample:  "counter//1",
			wantErr: true,
		},
		{
			name:    "Unexpected element",
			sample:  "counter/foo/bar/baz",
			wantErr: true,
		},
		{
			name:    "Unexpected element 2",
			sample:  "/counter/bar/1",
			wantErr: true,
		},
		{
			name:    "Unexpected element 3",
			sample:  "counter/bar/1/",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Metric{}
			err := m.Decode(tt.sample)
			if tt.wantErr {
				assert.Error(t, err)
			} else if assert.NoError(t, err) {
				assert.Equal(t, tt.want, m)
			}
		})
	}
}

func TestMetricEncode(t *testing.T) {
	g := Gauge(0.000001)
	c := Counter(1)

	tests := []struct {
		name   string
		sample *Metric
		want   string
	}{
		{
			name: "Gauge test",
			sample: &Metric{
				Name:  "foo",
				Value: Value(&g),
			},
			want: "gauge/foo/" + g.Encode(),
		},
		{
			name: "Counter test",
			sample: &Metric{
				Name:  "bar",
				Value: Value(&c),
			},
			want: "counter/bar/" + c.Encode(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.sample.Encode())
		})
	}
}

func TestMetricType(t *testing.T) {
	var m = Metric{}
	assert.Equal(t, MetricType, m.Type())
}
