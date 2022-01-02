package httplib

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter(root http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	router := chi.NewRouter()
	for _, mw := range middlewares {
		router.Use(mw)
	}
	router.Mount("/", root)
	return router
}

func Status(writer http.ResponseWriter, code int, message interface{}) {
	var err string
	if message == nil {
		err = fmt.Sprintf("%d %s", code, http.StatusText(code))
	} else {
		err = fmt.Sprintf("%d %s: %v", code, http.StatusText(code), message)
	}
	http.Error(writer, err, code)
}
