package handlers

import (
	"fmt"
	"log"
	"net/http"
	"sort"

	"github.com/go-chi/chi/v5"

	"github.com/zhupanovdm/go-runtime-monitor/model/metric"
	"github.com/zhupanovdm/go-runtime-monitor/pkg/httplib"
	"github.com/zhupanovdm/go-runtime-monitor/service/monitor"
	"github.com/zhupanovdm/go-runtime-monitor/view"
)

type MetricsHandler struct {
	*chi.Mux
	mon monitor.MetricsMonitorService
}

func (h *MetricsHandler) Update() http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		typ := metric.Type(chi.URLParam(req, "type"))
		if err := typ.Validate(); err != nil {
			httplib.Error(resp, http.StatusNotImplemented, fmt.Errorf("type %v is not supported yet", typ))
			return
		}

		value, err := typ.Parse(chi.URLParam(req, "value"))
		if err != nil {
			httplib.Error(resp, http.StatusBadRequest, fmt.Errorf("parsing error: %s", err))
			return
		}

		mtr := &metric.Metric{
			ID:    chi.URLParam(req, "id"),
			Value: value,
		}

		if err := h.mon.Save(ctx, mtr); err != nil {
			log.Printf("failed to persist metric %v: %v", req.URL, err)
			httplib.Error(resp, http.StatusInternalServerError, nil)
		}
	}
}

func (h *MetricsHandler) Value() http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		typ := metric.Type(chi.URLParam(req, "type"))
		if err := typ.Validate(); err != nil {
			httplib.Error(resp, http.StatusNotImplemented, fmt.Errorf("type %v is not supported yet", typ))
			return
		}

		id := chi.URLParam(req, "id")
		mtr, err := h.mon.Get(ctx, id, typ)
		if err != nil {
			log.Printf("failed to read metric %v: %v", req.URL, err)
			httplib.Error(resp, http.StatusInternalServerError, nil)
		}
		if mtr == nil {
			httplib.Error(resp, http.StatusNotFound, fmt.Errorf("%s (%v) metric not found", id, typ))
			return
		}

		if _, err := resp.Write([]byte(mtr.Value.String())); err != nil {
			log.Printf("error while writing body: %v", err)
			httplib.Error(resp, http.StatusInternalServerError, nil)
		}
	}
}

func (h *MetricsHandler) GetAll() http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		all, err := h.mon.GetAll(ctx)
		if err != nil {
			log.Printf("failed to fetch persisted mon: %v", err)
			httplib.Error(resp, http.StatusInternalServerError, nil)
		}

		sort.Sort(metric.ByString(all))

		if err := view.Index.Execute(resp, all); err != nil {
			log.Printf("error while writing body: %v", err)
			httplib.Error(resp, http.StatusInternalServerError, nil)
		}
	}
}

func NewMetricsHandler(mon monitor.MetricsMonitorService) *MetricsHandler {
	handler := &MetricsHandler{
		Mux: chi.NewRouter(),
		mon: mon,
	}

	handler.Get("/", handler.GetAll())
	handler.Route("/update", func(r chi.Router) {
		r.Post("/{type}/{id}/{value}", handler.Update())
	})
	handler.Route("/value", func(r chi.Router) {
		r.Get("/{type}/{id}", handler.Value())
	})

	return handler
}
