package handlers

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/go-chi/chi/v5"

	"github.com/zhupanovdm/go-runtime-monitor/model/metric"
	"github.com/zhupanovdm/go-runtime-monitor/pkg/httplib"
	"github.com/zhupanovdm/go-runtime-monitor/pkg/logging"
	"github.com/zhupanovdm/go-runtime-monitor/service/monitor"
	"github.com/zhupanovdm/go-runtime-monitor/view"
)

const handlerName = "Monitor HTTP handler"

type MetricsMonitorHandler struct {
	*chi.Mux
	mon monitor.MetricsMonitorService
}

func (h *MetricsMonitorHandler) Update() http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		ctx, _ := logging.SetIfAbsentCID(req.Context(), logging.NewCID())
		_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(handlerName), logging.WithCID(ctx))
		logger.Info().Msg("handling [Update]")

		typ := metric.Type(chi.URLParam(req, "type"))
		logger.UpdateContext(logging.LogCtxFrom(typ))

		if err := typ.Validate(); err != nil {
			logger.Err(err).Msg("unsupported type")
			httplib.Error(resp, http.StatusNotImplemented, fmt.Errorf("type %v is not supported yet", typ))
			return
		}

		value, err := typ.Parse(chi.URLParam(req, "value"))
		if err != nil {
			logger.Err(err).Msg("malformed metric value")
			httplib.Error(resp, http.StatusBadRequest, fmt.Errorf("parsing error: %s", err))
			return
		}

		mtr := &metric.Metric{
			ID:    chi.URLParam(req, "id"),
			Value: value,
		}
		logger.UpdateContext(logging.LogCtxFrom(mtr))

		ctx = logging.SetLogger(ctx, logger)
		if err := h.mon.Update(ctx, mtr); err != nil {
			logger.Err(err).Msg("failed to persist metric")
			httplib.Error(resp, http.StatusInternalServerError, nil)
		}
	}
}

func (h *MetricsMonitorHandler) Value() http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		ctx, _ := logging.SetIfAbsentCID(req.Context(), logging.NewCID())
		_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(handlerName), logging.WithCID(ctx))
		logger.Info().Msg("handling [Value]")

		typ := metric.Type(chi.URLParam(req, "type"))
		logger.UpdateContext(logging.LogCtxFrom(typ))

		if err := typ.Validate(); err != nil {
			logger.Err(err).Msg("unsupported type")
			httplib.Error(resp, http.StatusNotImplemented, fmt.Errorf("type %v is not supported yet", typ))
			return
		}

		id := chi.URLParam(req, "id")
		logger.UpdateContext(logging.LogCtxKeyStr(logging.MetricIDKey, id))

		ctx = logging.SetLogger(ctx, logger)
		mtr, err := h.mon.Get(ctx, id, typ)
		if err != nil {
			logger.Err(err).Msg("metric read failed")
			httplib.Error(resp, http.StatusInternalServerError, nil)
		}

		if mtr == nil {
			logger.Warn().Msg("requested metric not found")
			httplib.Error(resp, http.StatusNotFound, fmt.Errorf("%s (%v) metric not found", id, typ))
			return
		}

		if _, err := resp.Write([]byte(mtr.Value.String())); err != nil {
			logger.Err(err).Msg("failed to write response body")
			httplib.Error(resp, http.StatusInternalServerError, nil)
		}
	}
}

func (h *MetricsMonitorHandler) GetAll() http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		ctx, _ := logging.SetIfAbsentCID(req.Context(), logging.NewCID())
		_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(handlerName), logging.WithCID(ctx))
		logger.Info().Msg("handling [GetAll]")

		all, err := h.mon.GetAll(ctx)
		if err != nil {
			logger.Err(err).Msg("failed to query metrics")
			httplib.Error(resp, http.StatusInternalServerError, nil)
		}
		logger.Trace().Msgf("got %d records", len(all))

		sort.Sort(metric.ByString(all))

		if err := view.Index.Execute(resp, all); err != nil {
			logger.Err(err).Msg("failed to write response body")
			httplib.Error(resp, http.StatusInternalServerError, nil)
		}
	}
}

func NewMetricsHandler(mon monitor.MetricsMonitorService) *MetricsMonitorHandler {
	handler := &MetricsMonitorHandler{
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
