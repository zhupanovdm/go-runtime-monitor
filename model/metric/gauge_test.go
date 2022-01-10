package metric

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGauge_Parse(t *testing.T) {
	tests := []struct {
		name    string
		g       Gauge
		arg     string
		want    Gauge
		wantErr bool
	}{
		{
			name: "Basic test",
			arg:  ".5",
			want: Gauge(.5),
		},
		{
			name: "Int",
			arg:  "5",
			want: Gauge(5),
		},
		{
			name: "Negative",
			arg:  "-.5",
			want: Gauge(-.5),
		},
		{
			name: "Zero",
			arg:  "0",
			want: Gauge(0),
		},
		{
			name: "Max value",
			arg:  "1.79769313486231570814527423731704356798070e+308",
			want: Gauge(math.MaxFloat64),
		},
		{
			name: "Min value",
			arg:  "4.9406564584124654417656879286822137236505980e-324",
			want: Gauge(math.SmallestNonzeroFloat64),
		},
		{
			name:    "Incorrect",
			arg:     "foo",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.g.Parse(tt.arg)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				if assert.NoError(t, err) {
					assert.Equal(t, tt.want, tt.g,
						fmt.Sprintf("%T.Parse(\"%s\") affects to %T(%v)", tt.g, tt.arg, tt.want, tt.want))
				}
			}
		})
	}
}

func TestGauge_String(t *testing.T) {
	tests := []struct {
		name string
		g    Gauge
		want string
	}{
		{
			name: "Basic test",
			g:    Gauge(.5),
			want: "0.500",
		},
		{
			name: "Zero",
			g:    Gauge(0),
			want: "0.000",
		},
		{
			name: "Negative",
			g:    Gauge(-.5),
			want: "-0.500",
		},
		{
			name: "Precision overflow (rounds)",
			g:    Gauge(.0005),
			want: "0.001",
		},
		{
			name: "Precision overflow (drops)",
			g:    Gauge(.0001),
			want: "0.000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.g.String())
		})
	}
}

func TestGauge_Type(t *testing.T) {
	t.Run("Basic test", func(t *testing.T) {
		assert.Equal(t, GaugeType, Gauge(0).Type())
	})
}
