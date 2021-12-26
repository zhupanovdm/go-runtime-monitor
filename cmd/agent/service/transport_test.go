package service

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendToMonitorServer(t *testing.T) {
	var req, content string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req = r.Method + " " + r.URL.Path
		content = r.Header.Get("Content-Type")
	}))

	baseURL, _ := url.Parse(server.URL)
	assert.NoError(t, sendToMonitorServer(baseURL, "foo"))
	assert.Equal(t, "POST /update/foo", req)
	assert.Contains(t, content, "text/plain")
}

func TestSendToMonitorServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))

	baseURL, _ := url.Parse(server.URL)
	assert.Error(t, sendToMonitorServer(baseURL, "foo"))
}
