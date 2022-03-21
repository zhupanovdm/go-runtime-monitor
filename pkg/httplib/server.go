package httplib

import (
	"fmt"
	"io"
	"net/http"
)

// ResponseCustomWriter ables to override standard http.ResponseWriter
type ResponseCustomWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

// Write wries bytes to specified custom writer
func (w ResponseCustomWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

// Error writes response with specified HTTP status and corresponding explaining text
func Error(writer http.ResponseWriter, code int, message interface{}) {
	var err string
	if message == nil {
		err = fmt.Sprintf("%d %s", code, http.StatusText(code))
	} else {
		err = fmt.Sprintf("%d %s: %v", code, http.StatusText(code), message)
	}
	http.Error(writer, err, code)
}
