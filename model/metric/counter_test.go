package metric

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestCounter_Parse(t *testing.T) {
	tests := []struct {
		name    string
		c       Counter
		arg     string
		want    Counter
		wantErr bool
	}{
		{
			name: "Basic test",
			arg:  "5",
			want: Counter(5),
		},
		{
			name: "Negative",
			arg:  "-5",
			want: Counter(-5),
		},
		{
			name: "Zero",
			arg:  "0",
			want: Counter(0),
		},
		{
			name: "Max value",
			arg:  "9223372036854775807",
			want: Counter(math.MaxInt64),
		},
		{
			name: "Min value",
			arg:  "-9223372036854775808",
			want: Counter(math.MinInt64),
		},
		{
			name:    "Float",
			arg:     "0.5",
			wantErr: true,
		},
		{
			name:    "Incorrect",
			arg:     "foo",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		err := tt.c.Parse(tt.arg)
		if tt.wantErr {
			assert.Error(t, err)
		} else {
			if assert.NoError(t, err) {
				assert.Equal(t, tt.c, tt.want,
					fmt.Sprintf("%T.Parse(\"%s\") affects to %T(%v)", tt.c, tt.arg, tt.want, tt.want))
			}
		}
	}
}

func TestCounter_String(t *testing.T) {
	tests := []struct {
		name string
		c    Counter
		want string
	}{
		{
			name: "Basic test",
			c:    Counter(5),
			want: "5",
		},
		{
			name: "Zero",
			c:    Counter(0),
			want: "0",
		},
		{
			name: "Negative",
			c:    Counter(-5),
			want: "-5",
		},
		{
			name: "Max value",
			c:    Counter(math.MaxInt64),
			want: "9223372036854775807",
		},
		{
			name: "Min value",
			c:    Counter(math.MinInt64),
			want: "-9223372036854775808",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.c.String())
		})
	}
}

func TestCounter_Type(t *testing.T) {
	t.Run("Basic test", func(t *testing.T) {
		assert.Equal(t, CounterType, Counter(0).Type())
	})
}
