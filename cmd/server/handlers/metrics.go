package handlers

import (
	"fmt"
	"log"
	"net/http"
	"sort"

	"github.com/go-chi/chi/v5"

	"github.com/zhupanovdm/go-runtime-monitor/cmd/server/service"
	"github.com/zhupanovdm/go-runtime-monitor/cmd/server/view"

	"github.com/zhupanovdm/go-runtime-monitor/internal/model/metric"
)

type MetricsHandler struct {
	*chi.Mux
	metrics service.Metrics
}

func (h *MetricsHandler) Update() http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		typ := chi.URLParam(req, "type")
		value, ok := metric.Type(typ).NewValue()
		if !ok {
			status(resp, http.StatusNotImplemented, fmt.Errorf("type %s is not supported yet", typ))
			return
		}

		v := chi.URLParam(req, "value")
		if err := value.Parse(v); err != nil {
			status(resp, http.StatusBadRequest, fmt.Errorf("error while reading metric value: %s", err))
			return
		}

		m := metric.Metric{
			ID:    chi.URLParam(req, "id"),
			Value: value,
		}

		if err := h.metrics.Save(m); err != nil {
			log.Printf("failed to persist metric %v: %v", req.URL, err)
			status(resp, http.StatusInternalServerError, nil)
		}
	}
}

func (h *MetricsHandler) Value() http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		t := metric.Type(chi.URLParam(req, "type"))
		_, ok := t.NewValue()
		if !ok {
			status(resp, http.StatusNotImplemented, fmt.Errorf("type %s is not supported yet", t))
			return
		}

		id := chi.URLParam(req, "id")
		m, err := h.metrics.Get(id, t)
		if err != nil {
			log.Printf("failed to read metric %v: %v", req.URL, err)
			status(resp, http.StatusInternalServerError, nil)
		}
		if m == nil {
			status(resp, http.StatusNotFound, fmt.Errorf("%s (%v) metric not found", id, t))
			return
		}

		if _, err := resp.Write([]byte(m.String())); err != nil {
			log.Printf("error while writing body: %v", err)
			status(resp, http.StatusInternalServerError, nil)
		}
	}
}

func (h *MetricsHandler) GetAll() http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		all, err := h.metrics.GetAll()
		if err != nil {
			log.Printf("failed to fetch persisted metrics: %v", err)
			status(resp, http.StatusInternalServerError, nil)
		}

		sort.Sort(all)

		if err := view.Index.Execute(resp, all); err != nil {
			log.Printf("error while writing body: %v", err)
			status(resp, http.StatusInternalServerError, nil)
		}
	}
}

func NewMetricsHandler(metricsService service.Metrics) *MetricsHandler {
	h := &MetricsHandler{
		Mux:     chi.NewRouter(),
		metrics: metricsService,
	}

	h.Get("/", h.GetAll())
	h.Route("/update", func(r chi.Router) {
		r.Post("/{type}/{id}/{value}", h.Update())
	})
	h.Route("/value", func(r chi.Router) {
		r.Get("/{type}/{id}", h.Value())
	})

	return h
}

func status(writer http.ResponseWriter, code int, message interface{}) {
	var err string
	if message == nil {
		err = fmt.Sprintf("%d %s", code, http.StatusText(code))
	} else {
		err = fmt.Sprintf("%d %s: %v", code, http.StatusText(code), message)
	}
	http.Error(writer, err, code)
}
