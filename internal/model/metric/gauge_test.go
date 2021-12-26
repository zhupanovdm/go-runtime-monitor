package metric

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGaugeParse(t *testing.T) {
	tests := []struct {
		name    string
		g       Gauge
		sample  string
		want    Gauge
		wantErr bool
	}{
		{
			name:   "Basic test",
			sample: "1.1",
			want:   1.1,
		},
		{
			name:   "Int test",
			sample: "1",
			want:   1,
		},
		{
			name:   "Zero value",
			sample: "0",
			want:   0,
		},
		{
			name:   "Negative value",
			sample: "-1",
			want:   -1,
		},
		{
			name:   "Long number",
			sample: "9223372036854775807",
			want:   1<<63 - 1,
		},
		{
			name:    "Non number",
			sample:  "foo",
			wantErr: true,
		},
		{
			name:    "Empty value",
			sample:  "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.g.Parse(tt.sample)
			if tt.wantErr {
				assert.Error(t, err)
			} else if assert.NoError(t, err) {
				assert.Equal(t, tt.want, tt.g)
			}
		})
	}
}

func TestGaugeString(t *testing.T) {
	tests := []struct {
		name string
		g    Gauge
		want string
	}{
		{
			name: "Basic test",
			g:    0.000001,
			want: "0.000001",
		},
		{
			name: "Zero value",
			g:    0,
			want: "0.000000",
		},
		{
			name: "Negative value",
			g:    -1,
			want: "-1.000000",
		},
		{
			name: "Long value",
			g:    1<<63 - 1,
			want: "9223372036854775808.000000",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.g.String())
		})
	}
}

func TestGaugeType(t *testing.T) {
	var g Gauge
	assert.Equal(t, GaugeType, g.Type())
}
