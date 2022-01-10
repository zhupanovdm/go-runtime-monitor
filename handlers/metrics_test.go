package handlers

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetricsHandler(t *testing.T) {
	ts := NewServer(&monitorServiceStub{})
	defer ts.Close()

	tests := []struct {
		name       string
		method     string
		url        string
		wantStatus int
		want       string
	}{
		{
			name:       "Update gauge",
			method:     "POST",
			url:        "/update/gauge/foo/1.23",
			wantStatus: http.StatusOK,
		},
		{
			name:       "Get gauge",
			method:     "GET",
			url:        "/value/gauge/foo",
			wantStatus: http.StatusOK,
			want:       "0.000",
		},
		{
			name:       "Get absent metric",
			method:     "GET",
			url:        "/value/gauge/not-found",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "Update counter",
			method:     "POST",
			url:        "/update/counter/foo/1",
			wantStatus: http.StatusOK,
		},
		{
			name:       "Get counter",
			method:     "GET",
			url:        "/value/counter/foo",
			wantStatus: http.StatusOK,
			want:       "0",
		},
		{
			name:       "Update with incorrect gauge value",
			method:     "POST",
			url:        "/update/gauge/baz/foo",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Update with incorrect counter value",
			method:     "POST",
			url:        "/update/counter/baz/1.23",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Update unknown metric type",
			method:     "POST",
			url:        "/update/bar/baz/1",
			wantStatus: http.StatusNotImplemented,
		},
		{
			name:       "Get unknown metric type",
			method:     "GET",
			url:        "/value/bar/baz",
			wantStatus: http.StatusNotImplemented,
		},
		{
			name:       "Update metric without id 1",
			method:     "POST",
			url:        "/update/counter/",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "Update metric without id 2",
			method:     "POST",
			url:        "/update/counter",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "Not allowed update http method",
			method:     "GET",
			url:        "/update/counter/foo/1",
			wantStatus: http.StatusMethodNotAllowed,
		},
		{
			name:       "Not allowed value http method",
			method:     "POST",
			url:        "/value/counter/foo",
			wantStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, result, _ := testRequest(t, ts, tt.method, tt.url, nil)
			if assert.Equal(t, tt.wantStatus, status) {
				if len(tt.want) != 0 {
					assert.Equal(t, []byte(tt.want), result)
				}
			}
		})
	}
}
