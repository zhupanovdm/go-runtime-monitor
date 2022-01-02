package http

import (
	"context"
	"flag"
	"fmt"
	metric2 "github.com/zhupanovdm/go-runtime-monitor/model/metric"
	"path"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/zhupanovdm/go-runtime-monitor/pkg/httplib"
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

func (c httpClient) Update(_ context.Context, mtr *metric2.Metric) error {
	resp, err := c.R().Post(path.Join("update", mtr.Type().String(), mtr.ID, mtr.Value.String()))
	if err != nil {
		return fmt.Errorf("error quering server: %w", err)
	}
	if err := httplib.MustBeOK(resp); err != nil {
		return err
	}
	return nil
}

func (c httpClient) Value(_ context.Context, id string, typ metric2.Type) (metric2.Value, error) {
	resp, err := c.R().Get(path.Join("value", typ.String(), id))
	if err != nil {
		return nil, fmt.Errorf("error quering server: %w", err)
	}
	if err := httplib.MustBeOK(resp); err != nil {
		return nil, err
	}
	return typ.Parse(string(resp.Body()))
}

func NewClient(cfg *Config) monitor.Provider {
	client := resty.New()
	client.SetBaseURL(cfg.Server)
	client.SetTimeout(cfg.Timeout)
	client.SetHeader("Content-Type", "text/plain")

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
