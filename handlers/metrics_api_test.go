package handlers

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zhupanovdm/go-runtime-monitor/config"
	"github.com/zhupanovdm/go-runtime-monitor/model/metric"
	"github.com/zhupanovdm/go-runtime-monitor/providers/monitor/model"
)

func TestMetricsApiHandler(t *testing.T) {
	ts := NewServer(&config.Config{}, &monitorServiceStub{})
	defer ts.Close()

	tests := []struct {
		name            string
		method          string
		url             string
		body            string
		wantStatus      int
		wantContentType string
		want            string
	}{
		{
			name:       "Update gauge",
			method:     "POST",
			url:        "/update",
			body:       `{"id":"foo","type":"gauge","value":1.23}`,
			wantStatus: http.StatusOK,
		},
		{
			name:            "Get gauge",
			method:          "POST",
			url:             "/value",
			body:            `{"id":"foo","type":"gauge"}`,
			wantStatus:      http.StatusOK,
			wantContentType: "application/json",
			want:            `{"id":"foo","type":"gauge","value":0}`,
		},
		{
			name:       "Update counter",
			method:     "POST",
			url:        "/update",
			body:       `{"id":"foo","type":"counter","delta":1}`,
			wantStatus: http.StatusOK,
		},
		{
			name:            "Get counter",
			method:          "POST",
			url:             "/value",
			body:            `{"id":"foo","type":"counter"}`,
			wantStatus:      http.StatusOK,
			wantContentType: "application/json",
			want:            `{"id":"foo","type":"counter","delta":0}`,
		},
		{
			name:       "Get absent metric",
			method:     "POST",
			url:        "/value",
			body:       `{"id":"not-found","type":"counter"}`,
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "Update with incorrect gauge value",
			method:     "POST",
			url:        "/update",
			body:       `{"id":"baz","type":"gauge","value":"foo"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Update with incorrect counter value",
			method:     "POST",
			url:        "/update",
			body:       `{"id":"baz","type":"counter","delta":1.23}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Update unknown metric type",
			method:     "POST",
			url:        "/update",
			body:       `{"id":"bar","type":"baz","value":1}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Get unknown metric type",
			method:     "POST",
			url:        "/value",
			body:       `{"id":"bar","type":"baz"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Update metric without id 1",
			method:     "POST",
			url:        "/update",
			body:       `{"type":"counter"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Update metric without id 2",
			method:     "POST",
			url:        "/update",
			body:       `{"id":"","type":"counter"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Not allowed get http method",
			method:     "GET",
			url:        "/update",
			body:       `{"id":"foo","type":"gauge","value":1.23}`,
			wantStatus: http.StatusMethodNotAllowed,
		},
		{
			name:       "Not allowed value http method",
			method:     "GET",
			url:        "/value",
			body:       `{"id":"foo","type":"gauge"}`,
			wantStatus: http.StatusMethodNotAllowed,
		},
		{
			name:       "Ignore hash",
			method:     "POST",
			url:        "/update",
			body:       `{"id":"foo","type":"gauge","value":1.23,"hash":"ff"}`,
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, result, hdr := testRequest(t, ts, tt.method, tt.url, []byte(tt.body))
			if assert.Equal(t, tt.wantStatus, status) {
				if len(tt.want) != 0 {
					assert.JSONEq(t, tt.want, string(result))
				}
				if len(tt.wantContentType) != 0 {
					assert.Contains(t, hdr.Get("Content-Type"), tt.wantContentType)
				}
			}
		})
	}
}

func TestMetricsApiHandlerWithVSignatureVerification(t *testing.T) {
	known := "secret"
	ts := NewServer(&config.Config{Key: known}, &monitorServiceStub{})
	defer ts.Close()

	var value = 0.99
	var zeroGauge float64 = 0
	var delta int64 = 99

	zeroGaugeBody := &model.Metrics{
		ID:    "foo",
		MType: string(metric.GaugeType),
		Value: &zeroGauge,
	}
	require.NoError(t, zeroGaugeBody.Sign(known), "failed to sign sample")

	tests := []struct {
		name        string
		method      string
		url         string
		body        model.Metrics
		signWithKey string
		wantStatus  int
		wantHash    string
	}{
		{
			name:   "Update gauge with signed value",
			method: "POST",
			url:    "/update",
			body: model.Metrics{
				ID:    "foo",
				MType: string(metric.GaugeType),
				Value: &value,
			},
			signWithKey: known,
			wantStatus:  http.StatusOK,
		},
		{
			name:   "Update counter with signed value",
			method: "POST",
			url:    "/update",
			body: model.Metrics{
				ID:    "bar",
				MType: string(metric.CounterType),
				Delta: &delta,
			},
			signWithKey: known,
			wantStatus:  http.StatusOK,
		},
		{
			name:   "Unsigned",
			method: "POST",
			url:    "/update",
			body: model.Metrics{
				ID:    "baz",
				MType: string(metric.GaugeType),
				Value: &value,
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:   "Update gauge with signed with unknown key",
			method: "POST",
			url:    "/update",
			body: model.Metrics{
				ID:    "foo",
				MType: string(metric.GaugeType),
				Value: &value,
			},
			signWithKey: "unknown",
			wantStatus:  http.StatusBadRequest,
		},
		{
			name:   "Get gauge signed",
			method: "POST",
			url:    "/value",
			body: model.Metrics{
				ID:    "foo",
				MType: string(metric.GaugeType),
			},
			wantStatus: http.StatusOK,
			wantHash:   zeroGaugeBody.Hash,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, result, _ := testRequest(t, ts, tt.method, tt.url, signedJsonBody(t, tt.body, tt.signWithKey))
			if assert.Equal(t, tt.wantStatus, status) {
				if len(tt.wantHash) != 0 {
					mtr := model.Metrics{}
					require.NoError(t, json.Unmarshal(result, &mtr), "can't unmarshal body")
					assert.Equal(t, tt.wantHash, mtr.Hash)
				}
			}
		})
	}
}

func signedJsonBody(t *testing.T, mtr model.Metrics, key string) []byte {
	body := &mtr
	if len(key) != 0 {
		require.NoError(t, body.Sign(key), "sample sign unsuccessful")
	}
	bytes, err := json.Marshal(body)
	require.NoError(t, err, "sample marshalling failed")
	return bytes
}
