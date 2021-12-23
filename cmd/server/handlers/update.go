package handlers

import (
	"fmt"
	"github.com/zhupanovdm/go-runtime-monitor/cmd/server/storage"
	"github.com/zhupanovdm/go-runtime-monitor/internal/measure"
	"net/http"
	"strings"
)

func UpdateGauge(gauges storage.GaugesRepository) Middleware {
	return func(resp http.ResponseWriter, req *http.Request, next Handler) {
		metric, err := decodeMetric(req.URL.Path)
		if err != nil {
			http.Error(resp, fmt.Sprint(err), http.StatusBadRequest)
			return
		}
		if err := gauges.Save(metric.Name, *(metric.Value.(*measure.Gauge))); err != nil {
			http.Error(resp, fmt.Sprintf("cant save metric: %v", err), http.StatusBadRequest)
			return
		}
		next.Do(resp, req)
	}
}

func UpdateCounter(counters storage.CountersRepository) Middleware {
	return func(resp http.ResponseWriter, req *http.Request, next Handler) {
		metric, err := decodeMetric(req.URL.Path)
		if err != nil {
			http.Error(resp, fmt.Sprint(err), http.StatusBadRequest)
			return
		}
		if err := counters.Save(metric.Name, *(metric.Value.(*measure.Counter))); err != nil {
			http.Error(resp, fmt.Sprintf("cant save metric: %v", err), http.StatusBadRequest)
			return
		}
		next.Do(resp, req)
	}
}

func decodeMetric(path string) (*measure.Metric, error) {
	metric := &measure.Metric{}
	s := strings.Split(path, "/")[2:]
	if err := metric.Decode(strings.Join(s, "/")); err != nil {
		return nil, fmt.Errorf("cant recognize %s: %v", s, err)
	}
	return metric, nil
}
