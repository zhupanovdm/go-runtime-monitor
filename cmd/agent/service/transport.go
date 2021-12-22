package service

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
)

func sendToMonitorServer(baseURL *url.URL, value string) error {
	target, _ := url.Parse(baseURL.String())
	target.Path = path.Join(target.Path, "update", value)
	resp, err := http.Post(target.String(), "text/plain", nil)
	if err != nil {
		return fmt.Errorf("unable to connect to server: %w", err)
	}
	err = resp.Body.Close()
	if err != nil {
		return fmt.Errorf("unable to close response body: %w", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("server responded: %d", resp.StatusCode)
	}
	return nil
}
