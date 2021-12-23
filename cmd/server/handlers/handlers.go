package handlers

import (
	"net/http"
)

type Middleware func(http.ResponseWriter, *http.Request, Handler)
type Handler func(http.ResponseWriter, *http.Request)

func (h Handler) Do(writer http.ResponseWriter, request *http.Request) {
	if h != nil {
		h(writer, request)
	}
}

func Handle(middlewares ...Middleware) http.HandlerFunc {
	if len(middlewares) == 0 {
		return func(writer http.ResponseWriter, request *http.Request) {}
	}

	var handler Handler
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = bind(middlewares[i], handler)
	}
	return http.HandlerFunc(handler)
}

func POST(resp http.ResponseWriter, req *http.Request, next Handler) {
	if req.Method != "POST" {
		http.Error(resp, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	next.Do(resp, req)
}

func bind(middleware Middleware, handler Handler) Handler {
	return func(writer http.ResponseWriter, request *http.Request) {
		middleware(writer, request, handler)
	}
}
