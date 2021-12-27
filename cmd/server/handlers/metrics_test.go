package handlers

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zhupanovdm/go-runtime-monitor/cmd/server/service"

	"github.com/zhupanovdm/go-runtime-monitor/internal/model/metric"
)

func TestMetricsHandler(t *testing.T) {
	ts := httptest.NewServer(NewMetricsHandler(serviceStub{}))
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
			code, result := testRequest(t, ts, tt.method, tt.url)
			if assert.Equal(t, tt.wantStatus, code) {
				if len(tt.want) != 0 {
					assert.Equal(t, tt.want, result)
				}
			}
		})
	}
}

func testRequest(t *testing.T, ts *httptest.Server, method, path string) (int, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	respBody, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	require.NoError(t, err)

	return resp.StatusCode, string(respBody)
}

var _ service.Metrics = (*serviceStub)(nil)

type serviceStub struct {
}

func (s serviceStub) Save(_ metric.Metric) error {
	return nil
}

func (s serviceStub) Get(_ string, t metric.Type) (metric.Value, error) {
	v, _ := t.NewValue()
	return v, nil
}

func (s serviceStub) GetAll() (metric.List, error) {
	return make(metric.List, 0), nil
}
