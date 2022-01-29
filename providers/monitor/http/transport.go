package http

import (
	"fmt"
	"net/url"

	"github.com/go-resty/resty/v2"

	"github.com/zhupanovdm/go-runtime-monitor/pkg/logging"
	"github.com/zhupanovdm/go-runtime-monitor/providers/monitor"
)

func NewClient(cfg *monitor.Config, name string) (*resty.Client, error) {
	baseURL, err := getURL(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to set client destination address: %w", err)
	}

	client := resty.New()
	client.SetBaseURL(baseURL.String())
	client.SetTimeout(cfg.Timeout)
	client.OnBeforeRequest(requestHandler(cfg, name))
	client.OnAfterResponse(responseHandler(cfg, name))
	return client, err
}

func getURL(cfg *monitor.Config) (*url.URL, error) {
	u, err := url.Parse(cfg.Address)
	if err != nil {
		return nil, fmt.Errorf("invalid client destination address: %s: %w", cfg.Address, err)
	}
	if u.Host == "" {
		u, err = url.Parse("http://" + cfg.Address)
		if err != nil {
			return nil, fmt.Errorf("invalid client destination address: %s: %w", cfg.Address, err)
		}
	}
	return u, nil
}

func requestHandler(*monitor.Config, string) func(*resty.Client, *resty.Request) error {
	return func(_ *resty.Client, req *resty.Request) error {
		ctx, cid := logging.SetIfAbsentCID(req.Context(), logging.NewCID())
		req.SetContext(ctx)
		req.SetHeader(logging.CorrelationIDHeader, cid)
		return nil
	}
}

func responseHandler(_ *monitor.Config, name string) func(*resty.Client, *resty.Response) error {
	return func(_ *resty.Client, resp *resty.Response) error {
		req := resp.Request
		ctx := req.Context()
		_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(name), logging.WithCID(ctx))
		logger.Trace().Msgf("%s %s [%s] %d", req.Method, req.URL, resp.Status(), resp.Size())
		return nil
	}
}
