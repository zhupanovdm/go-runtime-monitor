package v1

import (
	"context"
	"path"

	"github.com/go-resty/resty/v2"

	"github.com/zhupanovdm/go-runtime-monitor/model/metric"
	"github.com/zhupanovdm/go-runtime-monitor/pkg/httplib"
	"github.com/zhupanovdm/go-runtime-monitor/providers/monitor"
	"github.com/zhupanovdm/go-runtime-monitor/providers/monitor/http"
)

const clientName = "Monitor HTTP Client v.1"

var _ monitor.Provider = (*httpClient)(nil)

type httpClient struct {
	*resty.Client
}

func (c httpClient) Update(ctx context.Context, mtr *metric.Metric) error {
	resp, err := c.R().
		SetContext(ctx).
		Post(path.Join("update", mtr.Type().String(), mtr.ID, mtr.Value.String()))
	if err != nil {
		return err
	}
	if err := httplib.MustBeOK(resp.StatusCode()); err != nil {
		return err
	}
	return nil
}

func (c httpClient) Value(ctx context.Context, id string, typ metric.Type) (metric.Value, error) {
	resp, err := c.R().
		SetContext(ctx).
		Get(path.Join("value", typ.String(), id))
	if err != nil {
		return nil, err
	}
	if err := httplib.MustBeOK(resp.StatusCode()); err != nil {
		return nil, err
	}
	return typ.Parse(string(resp.Body()))
}

func NewClient(cfg *monitor.Config) monitor.Provider {
	return &httpClient{
		Client: http.NewClient(cfg, clientName).SetHeader("Content-Type", "text/plain"),
	}
}
