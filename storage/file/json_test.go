package file

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zhupanovdm/go-runtime-monitor/model/metric"
)

func Test_jsonReader_Read(t *testing.T) {
	tests := []struct {
		name     string
		sample   []string
		wantList metric.List
		wantErr  bool
	}{
		{
			name: "Basic test",
			sample: []string{
				`{"ID":"foo","Type":"gauge","Value":33.3}`,
				`{"ID":"bar","Type":"counter","Value":333}`,
			},
			wantList: metric.List{
				metric.NewGaugeMetric("foo", metric.Gauge(33.3)),
				metric.NewCounterMetric("bar", metric.Counter(333)),
			},
		},
		{
			name:     "Empty source",
			wantList: metric.List{},
		},
		{
			name:    "Failed decode",
			sample:  []string{`{"ID":"baz","Type":"counter","Value":3.33}`},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			src := io.NopCloser(bytes.NewBuffer([]byte(strings.Join(tt.sample, "\n"))))
			list, err := NewJsonReader(context.TODO(), src).Read()
			if tt.wantErr {
				assert.Error(t, err)
			} else if assert.NoError(t, err) {
				assert.Equal(t, tt.wantList, list)
			}
		})
	}
}

func Test_jsonWriter_Write(t *testing.T) {
	tests := []struct {
		name    string
		list    metric.List
		want    []string
		wantErr bool
	}{
		{
			name: "Basic test",
			list: metric.List{
				metric.NewGaugeMetric("foo", metric.Gauge(33.3)),
				metric.NewCounterMetric("bar", metric.Counter(333)),
			},
			want: []string{
				`{"ID":"foo","Type":"gauge","Value":33.3}`,
				`{"ID":"bar","Type":"counter","Value":333}`,
				``},
		},
		{
			name: "Empty list",
			want: []string{``},
		},
		{
			name: "Failed encode",
			list: metric.List{
				&metric.Metric{ID: "must-fail"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := bytes.NewBuffer([]byte{})
			dest := struct {
				io.ReadCloser
				io.Writer
			}{
				ReadCloser: io.NopCloser(buf),
				Writer:     buf,
			}
			err := NewJsonWriter(context.TODO(), dest).Write(tt.list)
			if tt.wantErr {
				assert.Error(t, err)
			} else if assert.NoError(t, err) {
				got := strings.Split(buf.String(), "\n")
				assert.Equal(t, len(tt.want), len(got))
				for i, expected := range tt.want {
					if len(expected) == 0 {
						assert.Equal(t, 0, len(got[i]))
					} else {
						assert.JSONEq(t, expected, got[i])
					}
				}
			}
		})
	}
}
