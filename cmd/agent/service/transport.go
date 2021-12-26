package service

import (
	"fmt"
	"path"
	"time"

	"github.com/go-resty/resty/v2"
)

var MonitorClientTimeout = 30 * time.Second

func sendToMonitorServer(client *resty.Client, value string) error {
	resp, err := client.R().Post(path.Join("update", value))
	if err != nil {
		return fmt.Errorf("error quering server: %w", err)
	}
	if resp.StatusCode() != 200 {
		return fmt.Errorf("server responded: %d", resp.StatusCode())
	}
	return nil
}

func monitorClient(baseURL string) *resty.Client {
	client := resty.New()
	client.SetTimeout(MonitorClientTimeout)
	client.SetHeader("Content-Type", "text/plain")
	client.SetBaseURL(baseURL)
	return client
}
