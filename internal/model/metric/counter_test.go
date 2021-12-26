package metric

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCounterParse(t *testing.T) {
	tests := []struct {
		name    string
		c       Counter
		sample  string
		want    Counter
		wantErr bool
	}{
		{
			name:   "Basic test",
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
			name:    "Float test",
			sample:  "100.0",
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
			err := tt.c.Parse(tt.sample)
			if tt.wantErr {
				assert.Error(t, err)
			} else if assert.NoError(t, err) {
				assert.Equal(t, tt.want, tt.c)
			}
		})
	}
}

func TestCounterString(t *testing.T) {
	tests := []struct {
		name string
		c    Counter
		want string
	}{
		{
			name: "Basic test",
			c:    1,
			want: "1",
		},
		{
			name: "Zero value",
			c:    0,
			want: "0",
		},
		{
			name: "Negative value",
			c:    -1,
			want: "-1",
		},
		{
			name: "Long value",
			c:    1<<63 - 1,
			want: "9223372036854775807",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.c.String())
		})
	}
}

func TestCounterType(t *testing.T) {
	var c Counter
	assert.Equal(t, CounterType, c.Type())
}
