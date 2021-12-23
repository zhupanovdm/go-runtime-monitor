package handlers

import (
	"github.com/stretchr/testify/assert"
	"github.com/zhupanovdm/go-runtime-monitor/cmd/server/storage"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdateCounter(t *testing.T) {
	var counters = storage.NewCounters()

	tests := []struct {
		name       string
		req        string
		wantStatus int
		dataCheck  func() bool
	}{
		{
			name:       "Basic test",
			req:        "/update/counter/bar/777",
			wantStatus: http.StatusOK,
			dataCheck: func() bool {
				v, ok := counters.Get("bar")
				return ok && v == 777
			},
		},
		{
			name:       "Incorrect metric value (float)",
			req:        "/update/counter/foo/1.5",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Incorrect metric value (string)",
			req:        "/update/counter/foo/abc",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Long request path",
			req:        "/update/counter/bar/1/",
			wantStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := httptest.NewRecorder()

			UpdateCounter(counters)(resp, httptest.NewRequest("POST", baseURL+tt.req, nil), nil)

			result := resp.Result()
			defer result.Body.Close()

			assert.Equal(t, tt.wantStatus, result.StatusCode)
			if tt.dataCheck != nil {
				assert.True(t, tt.dataCheck(), "persisted data check")
			}
		})
	}
}

func TestUpdateGauge(t *testing.T) {
	gauges := storage.NewGauges()

	tests := []struct {
		name       string
		req        string
		wantStatus int
		dataCheck  func() bool
	}{
		{
			name:       "Basic test",
			req:        "/update/gauge/foo/7.77",
			wantStatus: http.StatusOK,
			dataCheck: func() bool {
				v, ok := gauges.Get("foo")
				return ok && v == 7.77
			},
		},
		{
			name:       "Incorrect metric value (string)",
			req:        "/update/gauge/foo/abc",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Long request path",
			req:        "/update/gauge/bar/1.1/",
			wantStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := httptest.NewRecorder()

			UpdateGauge(gauges)(resp, httptest.NewRequest("POST", baseURL+tt.req, nil), nil)

			result := resp.Result()
			defer result.Body.Close()

			assert.Equal(t, tt.wantStatus, result.StatusCode)
			if tt.dataCheck != nil {
				assert.True(t, tt.dataCheck(), "persisted data check")
			}
		})
	}
}
