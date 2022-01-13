package handlers

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetricsApiHandler(t *testing.T) {
	ts := NewServer(&monitorServiceStub{})
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
			wantStatus: http.StatusNotImplemented,
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
