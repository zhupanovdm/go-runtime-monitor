package httplib

import (
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
)

func MustBeOK(resp *resty.Response) error {
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("server responded with %d: %s", resp.StatusCode(), http.StatusText(resp.StatusCode()))
	}
	return nil
}
