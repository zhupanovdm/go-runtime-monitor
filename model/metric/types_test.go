package metric

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestType_New(t *testing.T) {
	gauge := Gauge(0)
	counter := Counter(0)

	tests := []struct {
		name      string
		t         Type
		wantValue Value
		wantErr   bool
	}{
		{
			name:      "Gauge creation",
			t:         GaugeType,
			wantValue: &gauge,
		},
		{
			name:      "Counter creation",
			t:         CounterType,
			wantValue: &counter,
		},
		{
			name:    "Unknown type creation",
			t:       "foo",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, err := tt.t.New()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				if assert.NoError(t, err) {
					assert.Equal(t, v, tt.wantValue,
						fmt.Sprintf("Type(\"%v\").New() must produce %T(%v)", tt.t, tt.wantValue, tt.wantValue))
				}
			}
		})
	}
}

func TestType_Parse(t *testing.T) {
	gauge := Gauge(.5)
	counter := Counter(5)

	tests := []struct {
		name      string
		t         Type
		arg       string
		wantValue Value
		wantErr   bool
	}{
		{
			name:      "Gauge parse",
			t:         GaugeType,
			arg:       ".5",
			wantValue: &gauge,
		},
		{
			name:      "Counter parse",
			t:         CounterType,
			arg:       "5",
			wantValue: &counter,
		},
		{
			name:    "Gauge parse error",
			t:       GaugeType,
			arg:     "must-fail",
			wantErr: true,
		},
		{
			name:    "Counter parse error",
			t:       CounterType,
			arg:     "must-fail",
			wantErr: true,
		},
		{
			name:    "Unknown type parse",
			t:       "foo",
			arg:     "5",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, err := tt.t.Parse(tt.arg)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				if assert.NoError(t, err) {
					assert.Equal(t, v, tt.wantValue,
						fmt.Sprintf("Type(\"%v\").Parse(\"%s\") must produce %T(%v)", tt.t, tt.arg, tt.wantValue, tt.wantValue))
				}
			}
		})
	}
}

func TestType_String(t *testing.T) {
	tests := []struct {
		name string
		t    Type
		want string
	}{
		{
			name: "Custom type stringer",
			t:    "foo",
			want: "foo",
		},
		{
			name: "Gauge stringer",
			t:    GaugeType,
			want: "gauge",
		},
		{
			name: "Counter stringer",
			t:    CounterType,
			want: "counter",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.t.String(), tt.want)
		})
	}
}

func TestType_Validate(t *testing.T) {
	tests := []struct {
		name    string
		t       Type
		wantErr bool
	}{
		{
			name: "Gauge type validate",
			t:    GaugeType,
		},
		{
			name: "Counter type validate",
			t:    CounterType,
		},
		{
			name:    "Unknown type validate",
			t:       "foo",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.t.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			err := tt.t.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
