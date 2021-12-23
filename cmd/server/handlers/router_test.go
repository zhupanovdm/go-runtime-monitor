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
			name:       "Path without /update part",
			method:     "POST",
			url:        "/counter/foo/1",
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

			assert.Equal(t, tt.wantStatus, resp.Result().StatusCode)

			_ = resp.Result().Body.Close()
		})
	}
}
