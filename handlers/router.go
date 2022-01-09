package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewMetricsRouter(metricsHandler *MetricsHandler, metricsAPI *MetricsAPIHandler) http.Handler {
	router := chi.NewRouter()
	router.Get("/", metricsHandler.GetAll)
	router.Route("/update", func(r chi.Router) {
		r.Post("/", metricsAPI.Update)
		r.Post("/{type}/{id}/{value}", metricsHandler.Update)
	})
	router.Route("/value", func(r chi.Router) {
		r.Post("/", metricsAPI.Value)
		r.Get("/{type}/{id}", metricsHandler.Value)
	})
	return router
}
