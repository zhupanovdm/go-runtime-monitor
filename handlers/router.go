package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewMetricsRouter(mon *MetricsHandler) http.Handler {
	router := chi.NewRouter()
	router.Get("/", mon.GetAll)
	router.Route("/update", func(r chi.Router) {
		r.Post("/{type}/{id}/{value}", mon.Update)
	})
	router.Route("/value", func(r chi.Router) {
		r.Get("/{type}/{id}", mon.Value)
	})
	return router
}
