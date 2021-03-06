package v1

import (
	"context"
	"errors"
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

func (c httpClient) UpdateBulk(context.Context, metric.List) error {
	return errors.New("unsupported operation")
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

// Deprecated: New version v.2 of provider should be used instead
func NewClient(cfg *monitor.Config) (monitor.Provider, error) {
	c, err := http.NewClient(cfg, clientName)
	if err != nil {
		return nil, err
	}
	return &httpClient{
		Client: c.SetHeader("Content-Type", "text/plain"),
	}, nil
}
