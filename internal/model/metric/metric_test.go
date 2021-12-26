package metric

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var gauge = Gauge(0.000001)
var counter = Counter(1)

func TestMetricString(t *testing.T) {
	tests := []struct {
		name   string
		sample Metric
		want   string
	}{
		{
			name: "Gauge test",
			sample: Metric{
				ID:    "foo",
				Value: Value(&gauge),
			},
			want: "gauge/foo/" + gauge.String(),
		},
		{
			name: "Counter test",
			sample: Metric{
				ID:    "bar",
				Value: Value(&counter),
			},
			want: "counter/bar/" + counter.String(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.sample.String())
		})
	}
}
