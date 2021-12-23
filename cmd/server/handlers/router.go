package handlers

import (
	"github.com/zhupanovdm/go-runtime-monitor/cmd/server/storage"
	"net/http"
)

func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/update/", Handle(POST, Status(http.StatusNotImplemented)))

	mux.Handle("/update/gauge", http.NotFoundHandler())
	mux.Handle("/update/gauge/", Handle(POST, CheckMetricKey, UpdateGauge(storage.NewGauges())))

	mux.Handle("/update/counter", http.NotFoundHandler())
	mux.Handle("/update/counter/", Handle(POST, CheckMetricKey, UpdateCounter(storage.NewCounters())))

	return mux
}
