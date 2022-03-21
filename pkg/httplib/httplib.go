// Package httplib contains recently used helper functions and types for convenient HTTP operating.
package httplib

import (
	"fmt"
	"net/http"
	"strings"
)

// Header is used to represent HTTP headers as a string.
type Header http.Header

func (h Header) String() string {
	var builder strings.Builder
	for k, v := range h {
		if builder.Len() != 0 {
			builder.WriteByte(' ')
		}
		builder.WriteString(fmt.Sprintf("%s:%s", k, strings.Join(v, ";")))
	}
	return builder.String()
}

// MustBeOK returns true if a given HTTP status is OK, otherwise false.
func MustBeOK(code int) error {
	if code != http.StatusOK {
		return fmt.Errorf("server responded with %d: %s", code, http.StatusText(code))
	}
	return nil
}
