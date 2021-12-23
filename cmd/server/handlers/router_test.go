package handlers

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewRouter(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		url        string
		wantStatus int
	}{
		{
			name:       "Update gauge",
			method:     "POST",
			url:        "/update/gauge/foo/1.23",
			wantStatus: http.StatusOK,
		},
		{
			name:       "Update counter",
			method:     "POST",
			url:        "/update/counter/foo/1",
			wantStatus: http.StatusOK,
		},
		{
			name:       "Unknown metric type",
			method:     "POST",
			url:        "/update/bar/baz/1",
			wantStatus: http.StatusNotImplemented,
		},
		{
			name:       "Metric without key",
			method:     "POST",
			url:        "/update/counter/",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "Metric without key",
			method:     "POST",
			url:        "/update/counter",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "Not allowed method in /update path",
			method:     "GET",
			url:        "/update/counter/foo/1",
			wantStatus: http.StatusMethodNotAllowed,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRouter()
			resp := httptest.NewRecorder()

			r.ServeHTTP(resp, httptest.NewRequest(tt.method, baseURL+tt.url, nil))

			result := resp.Result()
			defer func() { _ = result.Body.Close() }()

			assert.Equal(t, tt.wantStatus, result.StatusCode)
		})
	}
}
