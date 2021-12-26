package service

import (
	"fmt"
	"net/url"
	"path"

	"github.com/go-resty/resty/v2"
)

func sendToMonitorServer(baseURL *url.URL, value string) error {
	target, _ := url.Parse(baseURL.String())
	target.Path = path.Join(target.Path, "update", value)

	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "text/plain").
		Post(target.String())

	if err != nil {
		return fmt.Errorf("error quering server: %w", err)
	}
	if resp.StatusCode() != 200 {
		return fmt.Errorf("server responded: %d", resp.StatusCode())
	}

	return nil
}
