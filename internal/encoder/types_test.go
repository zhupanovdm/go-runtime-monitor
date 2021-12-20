package encoder

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCounterDecode(t *testing.T) {
	tests := []struct {
		name    string
		c       counter
		sample  string
		want    counter
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
			err := tt.c.Decode(tt.sample)
			if tt.wantErr {
				assert.Error(t, err)
			} else if assert.NoError(t, err) {
				assert.Equal(t, tt.want, tt.c)
			}
		})
	}
}

func TestCounterEncode(t *testing.T) {
	tests := []struct {
		name string
		c    counter
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
			assert.Equal(t, tt.want, tt.c.Encode())
		})
	}
}

func TestCounterType(t *testing.T) {
	var c counter
	assert.Equal(t, CounterType, c.Type())
}

func TestGaugeDecode(t *testing.T) {
	tests := []struct {
		name    string
		g       gauge
		sample  string
		want    gauge
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
			err := tt.g.Decode(tt.sample)
			if tt.wantErr {
				assert.Error(t, err)
			} else if assert.NoError(t, err) {
				assert.Equal(t, tt.want, tt.g)
			}
		})
	}
}

func TestGaugeEncode(t *testing.T) {
	tests := []struct {
		name string
		g    gauge
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
			assert.Equal(t, tt.want, tt.g.Encode())
		})
	}
}

func TestGaugeType(t *testing.T) {
	var g gauge
	assert.Equal(t, GaugeType, g.Type())
}

func TestMetricDecode(t *testing.T) {
	g := gauge(0.000001)
	c := counter(1)

	type fields struct {
		name  string
		value Encoder
	}
	tests := []struct {
		name    string
		want    fields
		sample  string
		wantErr bool
	}{
		{
			name:   "Gauge test",
			sample: "gauge/foo/0.000001",
			want: fields{
				name:  "foo",
				value: Encoder(&g),
			},
		},
		{
			name:   "Counter test",
			sample: "counter/bar/1",
			want: fields{
				name:  "bar",
				value: Encoder(&c),
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
			m := &metric{}
			err := m.Decode(tt.sample)
			if tt.wantErr {
				assert.Error(t, err)
			} else if assert.NoError(t, err) {
				assert.Equal(t, tt.want.name, m.Type())
				assert.Equal(t, tt.want.value, m.value)
			}
		})
	}
}

func TestMetricEncode(t *testing.T) {
	g := gauge(0.000001)
	c := counter(1)

	tests := []struct {
		name   string
		sample metric
		want   string
	}{
		{
			name: "Gauge test",
			sample: metric{
				name:  "foo",
				value: Encoder(&g),
			},
			want: "gauge/foo/" + g.Encode(),
		},
		{
			name: "Counter test",
			sample: metric{
				name:  "bar",
				value: Encoder(&c),
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
	var m = metric{name: "foo"}
	assert.Equal(t, "foo", m.Type())
}
