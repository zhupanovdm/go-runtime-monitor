package httplib

import (
	"fmt"
	"io"
	"net/http"
)

type ResponseCustomWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w ResponseCustomWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func Error(writer http.ResponseWriter, code int, message interface{}) {
	var err string
	if message == nil {
		err = fmt.Sprintf("%d %s", code, http.StatusText(code))
	} else {
		err = fmt.Sprintf("%d %s: %v", code, http.StatusText(code), message)
	}
	http.Error(writer, err, code)
}
