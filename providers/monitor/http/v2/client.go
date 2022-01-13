package v2

import (
	"bytes"
	"context"
	"encoding/json"

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
}

func (c httpClient) Update(ctx context.Context, mtr *metric.Metric) error {
	resp, err := c.R().SetContext(ctx).SetBody(model.NewFromCanonical(mtr)).Post("update")
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
	value = mtr.ToCanonical().Value
	return
}

func NewClient(cfg *monitor.Config) monitor.Provider {
	return &httpClient{
		Client: http.NewClient(cfg, clientName).SetHeader("Content-Type", "application/json"),
	}
}
