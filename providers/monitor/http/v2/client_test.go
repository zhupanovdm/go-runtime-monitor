package v2

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zhupanovdm/go-runtime-monitor/config"
	"github.com/zhupanovdm/go-runtime-monitor/model/metric"
	"github.com/zhupanovdm/go-runtime-monitor/providers/monitor"
	"github.com/zhupanovdm/go-runtime-monitor/providers/monitor/model"
)

func TestHttpClientUpdateWithSignVerification(t *testing.T) {
	key := "secret"
	server := httptest.NewServer(mockUpdateHandler(t, key))
	defer server.Close()

	var gauge = metric.Gauge(0)

	tests := []struct {
		name    string
		key     string
		metric  *metric.Metric
		wantErr bool
	}{
		{
			name: "Basic test",
			metric: &metric.Metric{
				ID:    "foo",
				Value: &gauge,
			},
			key: key,
		},
		{
			name: "Verification failed",
			metric: &metric.Metric{
				ID:    "foo",
				Value: &gauge,
			},
			key:     "unknown",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := newTestClient(server.URL, tt.key)
			require.NoError(t, err, "failed to create client")

			err = client.Update(context.TODO(), tt.metric)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestHttpClientValueWithSignVerification(t *testing.T) {
	key := "secret"
	server := httptest.NewServer(mockValueHandler(t, key))
	defer server.Close()

	tests := []struct {
		name    string
		key     string
		id      string
		typ     metric.Type
		wantErr bool
	}{
		{
			name: "Basic test",
			id:   "foo",
			typ:  metric.CounterType,
			key:  key,
		},
		{
			name:    "Verification failed",
			id:      "foo",
			typ:     metric.CounterType,
			key:     "unknown",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := newTestClient(server.URL, tt.key)
			require.NoError(t, err, "failed to create client")

			_, err = client.Value(context.TODO(), tt.id, tt.typ)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func newTestClient(addr string, key string) (monitor.Provider, error) {
	return NewClient(&monitor.Config{
		Config:  &config.Config{Address: addr, Key: key},
		Timeout: 1 * time.Second,
	})
}

func mockUpdateHandler(t *testing.T, key string) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		m := &model.Metrics{}
		require.NoError(t, json.NewDecoder(request.Body).Decode(m), "failed to decode body")
		defer request.Body.Close()
		if len(key) != 0 {
			if err := m.Verify(key); err != nil {
				writer.WriteHeader(http.StatusBadRequest)
			}
		}
	}
}

func mockValueHandler(t *testing.T, key string) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		m := &model.Metrics{}
		require.NoError(t, json.NewDecoder(request.Body).Decode(m), "failed to decode body")
		defer request.Body.Close()

		canonical := m.ToCanonical()
		require.NotNil(t, canonical)

		m = model.NewFromCanonical(canonical)

		require.NoError(t, m.Sign(key), "failed to sign response")
		require.NoError(t, json.NewEncoder(writer).Encode(m), "failed to marshall response")
	}
}
