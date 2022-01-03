package http

import (
	"context"
	"flag"
	"path"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/zhupanovdm/go-runtime-monitor/model/metric"
	"github.com/zhupanovdm/go-runtime-monitor/pkg/httplib"
	"github.com/zhupanovdm/go-runtime-monitor/pkg/logging"
	"github.com/zhupanovdm/go-runtime-monitor/providers/monitor"
)

var _ monitor.Provider = (*httpClient)(nil)

type httpClient struct {
	*resty.Client
}

type Config struct {
	Server  string
	Timeout time.Duration
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

func NewClient(cfg *Config) monitor.Provider {
	client := resty.New()
	client.SetBaseURL(cfg.Server)
	client.SetTimeout(cfg.Timeout)
	client.SetHeader("Content-Type", "text/plain")
	client.OnBeforeRequest(func(client *resty.Client, req *resty.Request) error {
		ctx, cid := logging.SetIfAbsentCID(req.Context(), logging.NewCID())
		req.SetContext(ctx)
		req.SetHeader(logging.CorrelationIDHeader, cid)
		return nil
	})
	client.OnAfterResponse(func(client *resty.Client, resp *resty.Response) error {
		req := resp.Request
		ctx := req.Context()
		_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName("Monitor HTTP Client"), logging.WithCID(ctx))
		logger.Trace().Msgf("%s %s [%s] %d", req.Method, req.URL, resp.Status(), resp.Size())
		return nil
	})
	return &httpClient{
		Client: client,
	}
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) FromCLI(flag *flag.FlagSet) *Config {
	flag.StringVar(&c.Server, "monitor-srv", "http://localhost:8080", "monitor server URL")
	flag.DurationVar(&c.Timeout, "client-timeout", 30*time.Second, "client timeout")
	return c
}
