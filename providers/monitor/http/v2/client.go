package v2

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"

	"github.com/zhupanovdm/go-runtime-monitor/model/metric"
	"github.com/zhupanovdm/go-runtime-monitor/pkg/httplib"
	"github.com/zhupanovdm/go-runtime-monitor/providers/monitor"
	"github.com/zhupanovdm/go-runtime-monitor/providers/monitor/http"
	"github.com/zhupanovdm/go-runtime-monitor/providers/monitor/model"
)

const clientName = "Monitor HTTP Client v.2"

var _ monitor.Provider = (*httpClient)(nil)

type httpClient struct {
	*resty.Client
	key string
}

func (c httpClient) Update(ctx context.Context, mtr *metric.Metric) error {
	body := model.NewFromCanonical(mtr)
	if len(c.key) != 0 {
		if err := body.Sign(c.key); err != nil {
			return err
		}
	}
	resp, err := c.R().SetContext(ctx).SetBody(body).Post("update")
	if err != nil {
		return err
	}
	if err = httplib.MustBeOK(resp.StatusCode()); err != nil {
		return err
	}
	return nil
}

func (c httpClient) Value(ctx context.Context, id string, typ metric.Type) (value metric.Value, err error) {
	mtr := &model.Metrics{
		ID:    id,
		MType: string(typ),
	}

	var resp *resty.Response
	if resp, err = c.R().SetContext(ctx).SetBody(mtr).Post("value"); err != nil {
		return
	}
	if err = httplib.MustBeOK(resp.StatusCode()); err != nil {
		return
	}

	mtr = &model.Metrics{}
	if err = json.NewDecoder(bytes.NewBuffer(resp.Body())).Decode(mtr); err != nil {
		return
	}
	if len(c.key) != 0 {
		if err := mtr.Verify(c.key); err != nil {
			return nil, fmt.Errorf("response verification failed: %w", err)
		}
	}
	value = mtr.ToCanonical().Value
	return
}

func NewClient(cfg *monitor.Config) monitor.Provider {
	return &httpClient{
		Client: http.NewClient(cfg, clientName).SetHeader("Content-Type", "application/json"),
		key:    cfg.Key,
	}
}
