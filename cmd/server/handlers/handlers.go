package handlers

import (
	"fmt"
	"net/http"
)

type Middleware func(http.ResponseWriter, *http.Request, Handler)
type Handler func(http.ResponseWriter, *http.Request)

func (h Handler) Do(w http.ResponseWriter, r *http.Request) {
	if h != nil {
		h(w, r)
	}
}

func Handle(m ...Middleware) http.HandlerFunc {
	if len(m) == 0 {
		return func(http.ResponseWriter, *http.Request) {}
	}

	var h Handler
	for i := len(m) - 1; i >= 0; i-- {
		h = bind(m[i], h)
	}
	return http.HandlerFunc(h)
}

func POST(w http.ResponseWriter, r *http.Request, next Handler) {
	if r.Method != "POST" {
		status(w, http.StatusMethodNotAllowed)
		return
	}
	next.Do(w, r)
}

func Status(code int) Middleware {
	return func(w http.ResponseWriter, _ *http.Request, _ Handler) {
		status(w, code)
	}
}

func status(w http.ResponseWriter, code int) {
	http.Error(w, fmt.Sprintf("%d %s", code, http.StatusText(code)), code)
}

func bind(m Middleware, next Handler) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		m(w, r, next)
	}
}
