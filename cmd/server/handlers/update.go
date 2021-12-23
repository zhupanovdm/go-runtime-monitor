package handlers

import (
	"fmt"
	"github.com/zhupanovdm/go-runtime-monitor/cmd/server/storage"
	"github.com/zhupanovdm/go-runtime-monitor/internal/measure"
	"net/http"
	"strings"
)

func UpdateGauge(gauges storage.GaugesRepository) Middleware {
	return func(w http.ResponseWriter, r *http.Request, next Handler) {
		metric, err := decodeMetric(r.URL.Path)
		if err != nil {
			status(w, http.StatusBadRequest)
			return
		}
		if err := gauges.Save(metric.Name, *(metric.Value.(*measure.Gauge))); err != nil {
			status(w, http.StatusInternalServerError)
			return
		}
		next.Do(w, r)
	}
}

func UpdateCounter(counters storage.CountersRepository) Middleware {
	return func(w http.ResponseWriter, r *http.Request, next Handler) {
		metric, err := decodeMetric(r.URL.Path)
		if err != nil {
			status(w, http.StatusBadRequest)
			return
		}
		if err := counters.Save(metric.Name, *(metric.Value.(*measure.Counter))); err != nil {
			status(w, http.StatusInternalServerError)
			return
		}
		next.Do(w, r)
	}
}

func CheckMetricKey(w http.ResponseWriter, r *http.Request, next Handler) {
	s := strings.Split(r.URL.Path, "/")[2:]
	if len(s) < 3 {
		http.NotFound(w, r)
		return
	}
	next.Do(w, r)
}

func decodeMetric(path string) (*measure.Metric, error) {
	metric := &measure.Metric{}
	s := strings.Split(path, "/")[2:]
	if err := metric.Decode(strings.Join(s, "/")); err != nil {
		return nil, fmt.Errorf("cant recognize %s: %v", s, err)
	}
	return metric, nil
}
