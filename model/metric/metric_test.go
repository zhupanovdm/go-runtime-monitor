package metric

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetric_MarshalJSON(t *testing.T) {
	c := Counter(777)
	g := Gauge(77.7)

	tests := []struct {
		name    string
		metric  Metric
		want    string
		wantErr bool
	}{
		{
			name: "Encode counter",
			metric: Metric{
				ID:    "foo",
				Value: &c,
			},
			want: `{"ID":"foo","Type":"counter","Value":777}`,
		},
		{
			name: "Encode gauge",
			metric: Metric{
				ID:    "foo",
				Value: &g,
			},
			want: `{"ID":"foo","Type":"gauge","Value":77.7}`,
		},
		{
			name: "Encode failed",
			metric: Metric{
				ID: "foo",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := json.Marshal(tt.metric)
			if tt.wantErr {
				assert.Error(t, err)
			} else if assert.NoError(t, err) {
				assert.JSONEq(t, tt.want, string(result))
			}
		})
	}
}

func TestMetric_UnmarshalJSON(t *testing.T) {
	c := Counter(777)
	g := Gauge(77.7)

	tests := []struct {
		name    string
		sample  string
		want    *Metric
		wantErr bool
	}{
		{
			name:   "Decode counter",
			sample: `{"ID":"foo","Type":"counter","Value":777}`,
			want: &Metric{
				ID:    "foo",
				Value: &c,
			},
		},
		{
			name:   "Decode gauge",
			sample: `{"ID":"foo","Type":"gauge","Value":77.7}`,
			want: &Metric{
				ID:    "foo",
				Value: &g,
			},
		},
		{
			name:    "Decode failed",
			sample:  `{"ID":"foo","Type":"counter","Value":77.7}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Metric{}
			err := json.Unmarshal([]byte(tt.sample), m)
			if tt.wantErr {
				assert.Error(t, err)
			} else if assert.NoError(t, err) {
				assert.Equal(t, tt.want, m)
			}
		})
	}
}

func TestMetric_String(t *testing.T) {
	c := Counter(777)
	g := Gauge(77.7)

	tests := []struct {
		name   string
		metric Metric
		want   string
	}{
		{
			name: "Gauge representation",
			metric: Metric{
				ID:    "foo",
				Value: &g,
			},
			want: "gauge/foo/77.700",
		},
		{
			name: "Counter representation",
			metric: Metric{
				ID:    "foo",
				Value: &c,
			},
			want: "counter/foo/777",
		},
		{
			name: "Absent value representation",
			metric: Metric{
				ID: "foo",
			},
			want: "?/foo/?",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.metric.String())
		})
	}
}
