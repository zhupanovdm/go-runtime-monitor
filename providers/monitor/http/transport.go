package http

import (
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"

	"github.com/zhupanovdm/go-runtime-monitor/pkg/logging"
	"github.com/zhupanovdm/go-runtime-monitor/providers/monitor"
)

func NewClient(cfg *monitor.Config, name string) *resty.Client {
	client := resty.New()

	if strings.HasPrefix(cfg.Address, "http") {
		client.SetBaseURL(cfg.Address)
	} else {
		client.SetBaseURL(fmt.Sprintf("http://%s", cfg.Address))
	}

	client.SetTimeout(cfg.Timeout)
	client.OnBeforeRequest(func(client *resty.Client, req *resty.Request) error {
		ctx, cid := logging.SetIfAbsentCID(req.Context(), logging.NewCID())
		req.SetContext(ctx)
		req.SetHeader(logging.CorrelationIDHeader, cid)
		return nil
	})
	client.OnAfterResponse(func(client *resty.Client, resp *resty.Response) error {
		req := resp.Request
		ctx := req.Context()
		_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(name), logging.WithCID(ctx))
		logger.Trace().Msgf("%s %s [%s] %d", req.Method, req.URL, resp.Status(), resp.Size())
		return nil
	})
	return client
}
