package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/zhupanovdm/go-runtime-monitor/config"
	"github.com/zhupanovdm/go-runtime-monitor/model/metric"
	"github.com/zhupanovdm/go-runtime-monitor/pkg/httplib"
	"github.com/zhupanovdm/go-runtime-monitor/pkg/logging"
	"github.com/zhupanovdm/go-runtime-monitor/providers/monitor/model"
	"github.com/zhupanovdm/go-runtime-monitor/service/monitor"
)

const metricsHandlerAPIName = "Metrics REST API handler"

type MetricsAPIHandler struct {
	monitor monitor.Monitor
	key     string
}

// Update godoc
// @Tags v2
// @Summary Updates single metric value
// @Description report monitor server of changed metric
// @ID v2metricsUpdate
// @Accept json
// @Param metric_data body model.Metrics true "Metric to update"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /update [post]
func (h *MetricsAPIHandler) Update(resp http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	ctx, _ := logging.SetIfAbsentCID(req.Context(), logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(metricsHandlerAPIName), logging.WithCID(ctx))
	logger.Info().Msg("handling [Update]")

	body, err := h.decodeRequestBody(req.Body)
	if err != nil {
		logger.Err(err).Msg("failed to process request body")
		httplib.Error(resp, http.StatusBadRequest, err)
		return
	}

	if err = body.Validate(model.CheckID, model.CheckValue, model.CheckType, model.CheckHash(h.key)); err != nil {
		logger.Err(err).Msg("validation failed")
		httplib.Error(resp, http.StatusBadRequest, err)
		return
	}

	mtr := body.ToCanonical()
	logger.UpdateContext(logging.LogCtxFrom(mtr))
	ctx = logging.SetLogger(ctx, logger)
	if err := h.monitor.Update(ctx, mtr); err != nil {
		logger.Err(err).Msg("failed to persist metric")
		httplib.Error(resp, http.StatusInternalServerError, nil)
		return
	}
}

// UpdateBulk godoc
// @Tags v2
// @Summary Updates multiple metrics values
// @Description Report monitor server of several changed metrics at once
// @ID v2metricsUpdateBulk
// @Accept json
// @Param metrics_list body []model.Metrics true "Metric to update"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /updates [post]
func (h *MetricsAPIHandler) UpdateBulk(resp http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	ctx, _ := logging.SetIfAbsentCID(req.Context(), logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(metricsHandlerAPIName), logging.WithCID(ctx))
	logger.Info().Msg("handling [Updates]")

	var metrics []model.Metrics
	if err := json.NewDecoder(req.Body).Decode(&metrics); err != nil {
		logger.Err(err).Msg("failed to process request body")
		httplib.Error(resp, http.StatusBadRequest, err)
		return
	}

	list := make(metric.List, 0, len(metrics))
	for _, m := range metrics {
		if err := m.Validate(model.CheckID, model.CheckValue, model.CheckType, model.CheckHash(h.key)); err != nil {
			logger.Err(err).Msg("validation failed")
			httplib.Error(resp, http.StatusBadRequest, err)
			return
		}
		list = append(list, m.ToCanonical())
	}
	if err := h.monitor.UpdateBulk(ctx, list); err != nil {
		logger.Err(err).Msg("failed to batch update metrics")
		httplib.Error(resp, http.StatusInternalServerError, nil)
	}
}

// Value godoc
// @Tags v2
// @Summary Queries requested metric value
// @Description Returns specified metric actual value
// @ID v2metricsValue
// @Accept json
// @Produce json
// @Success 200 {number} number "OK"
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "Not found"
// @Failure 500 {string} string "Internal server error"
// @Router /value [get]
func (h *MetricsAPIHandler) Value(resp http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	ctx, _ := logging.SetIfAbsentCID(req.Context(), logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(metricsHandlerAPIName), logging.WithCID(ctx))
	logger.Info().Msg("handling [Value]")

	body, err := h.decodeRequestBody(req.Body)
	if err != nil {
		logger.Err(err).Msg("failed to process request body")
		httplib.Error(resp, http.StatusBadRequest, err)
		return
	}

	if err = body.Validate(model.CheckID, model.CheckType); err != nil {
		logger.Err(err).Msg("validation failed")
		httplib.Error(resp, http.StatusBadRequest, err)
		return
	}

	mtr := body.ToCanonical()

	logger.UpdateContext(logging.LogCtxFrom(logging.LogCtxKeyStr(logging.MetricIDKey, mtr.ID), mtr.Type()))
	ctx = logging.SetLogger(ctx, logger)
	if mtr, err = h.monitor.Get(ctx, mtr.ID, mtr.Type()); err != nil {
		logger.Err(err).Msg("metric read failed")
		httplib.Error(resp, http.StatusInternalServerError, nil)
		return
	}
	if mtr == nil {
		logger.Warn().Msg("requested metric not found")
		httplib.Error(resp, http.StatusNotFound, errors.New("metric not found"))
		return
	}

	resp.Header().Set("Content-Type", "application/json")

	body = model.NewFromCanonical(mtr)
	if len(h.key) != 0 {
		if err := body.Sign(h.key); err != nil {
			logger.Err(err).Msg("signing failed")
			httplib.Error(resp, http.StatusBadRequest, err)
			return
		}
	}

	if err = json.NewEncoder(resp).Encode(body); err != nil {
		logger.Err(err).Msg("failed to encode response body")
		httplib.Error(resp, http.StatusInternalServerError, nil)
		return
	}
}

// Ping godoc
// @Tags Diag
// @Summary Service diagnostic method
// @Description Returns OK if service is up
// @ID diagPing
// @Success 200 {string} string "OK"
// @Failure 500 {string} string "Internal server error"
// @Router /ping [get]
func (h *MetricsAPIHandler) Ping(resp http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	ctx, _ := logging.SetIfAbsentCID(req.Context(), logging.NewCID())
	_, logger := logging.GetOrCreateLogger(ctx, logging.WithServiceName(metricsHandlerAPIName), logging.WithCID(ctx))
	logger.Info().Msg("handling [Ping]")

	if err := h.monitor.Ping(ctx); err != nil {
		logger.Err(err).Msg("failed to check monitor storage availability")
		httplib.Error(resp, http.StatusInternalServerError, nil) // IMHO should 503
	}
}

func (h *MetricsAPIHandler) decodeRequestBody(body io.Reader) (*model.Metrics, error) {
	metrics := &model.Metrics{}
	if err := json.NewDecoder(body).Decode(&metrics); err != nil {
		return nil, fmt.Errorf("decoder: error while decoding JSON: %w", err)
	}
	return metrics, nil
}

func NewMetricsAPIHandler(cfg *config.Config, service monitor.Monitor) *MetricsAPIHandler {
	return &MetricsAPIHandler{
		monitor: service,
		key:     cfg.Key,
	}
}
