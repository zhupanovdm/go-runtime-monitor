package httplib

import (
	"fmt"
	"net/http"
)

func MustBeOK(code int) error {
	if code != http.StatusOK {
		return fmt.Errorf("server responded with %d: %s", code, http.StatusText(code))
	}
	return nil
}
