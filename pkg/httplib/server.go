package httplib

import (
	"fmt"
	"net/http"
)

func Error(writer http.ResponseWriter, code int, message interface{}) {
	var err string
	if message == nil {
		err = fmt.Sprintf("%d %s", code, http.StatusText(code))
	} else {
		err = fmt.Sprintf("%d %s: %v", code, http.StatusText(code), message)
	}
	http.Error(writer, err, code)
}
