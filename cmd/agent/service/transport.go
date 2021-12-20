package service

import (
	"fmt"
	"github.com/zhupanovdm/go-runtime-monitor/internal/encoder"
	"net/http"
	"net/url"
	"path"
)

func ServerTransmitter(baseUrl *url.URL) func(encoder.Encoder) error {
	return func(value encoder.Encoder) error {
		target, _ := url.Parse(baseUrl.String())
		target.Path = path.Join(target.Path, "update", value.Encode())
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
}
