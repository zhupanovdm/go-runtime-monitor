package handlers

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zhupanovdm/go-runtime-monitor/config"
	"github.com/zhupanovdm/go-runtime-monitor/model/metric"
	"github.com/zhupanovdm/go-runtime-monitor/pkg/task"
	"github.com/zhupanovdm/go-runtime-monitor/service/monitor"
)

const notFoundSample = "not-found"

func NewServer(cfg *config.Config, svc monitor.Monitor) *httptest.Server {
	return httptest.NewServer(NewMetricsRouter(NewMetricsHandler(svc), NewMetricsAPIHandler(cfg, svc)))
}

var _ monitor.Monitor = (*monitorServiceStub)(nil)

type monitorServiceStub struct{}

func (s *monitorServiceStub) BackgroundTask() task.Task {
	return task.VoidTask
}

func (s *monitorServiceStub) GetAll(context.Context) (metric.List, error) {
	return make([]*metric.Metric, 0), nil
}

func (s *monitorServiceStub) Restore(context.Context) error {
	return nil
}

func (s *monitorServiceStub) Name() string {
	return "Stub monitor service"
}

func (s *monitorServiceStub) Get(_ context.Context, id string, typ metric.Type) (*metric.Metric, error) {
	if id == notFoundSample {
		return nil, nil
	}

	v, _ := typ.New()
	return &metric.Metric{
		ID:    id,
		Value: v,
	}, nil
}

func (s *monitorServiceStub) Update(context.Context, *metric.Metric) error {
	return nil
}

func (s *monitorServiceStub) Ping(context.Context) error {
	return nil
}

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body []byte) (int, []byte, http.Header) {
	if body == nil {
		body = []byte{}
	}
	req, err := http.NewRequest(method, ts.URL+path, bytes.NewBuffer(body))
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	respBody, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	require.NoError(t, err)

	return resp.StatusCode, respBody, resp.Header
}
