package handlers

import (
	"github.com/zhupanovdm/go-runtime-monitor/cmd/server/storage"
	"net/http"
)

func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/update/gauge/", Handle(POST, UpdateGauge(storage.NewGauges())))
	mux.Handle("/update/counter/", Handle(POST, UpdateCounter(storage.NewCounters())))
	return mux
}
