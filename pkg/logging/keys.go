package logging

import "github.com/rs/zerolog"

const (
	// CorrelationIDKey is used to track unique request.
	CorrelationIDKey = "cid"

	// CorrelationIDChangedKey is used to track request multiplexing.
	CorrelationIDChangedKey = "cid_changed"

	// CorrelationIDHeader is used to transport Correlation ID context value via the HTTP header.
	CorrelationIDHeader = "X-CorrelationID"

	// ServiceKey is used to track concrete service side effects.
	ServiceKey = "svc"

	// PollCountKey is used to track metrics polling attempt of agent.CollectorService.
	PollCountKey = "poll"

	// MetricIDKey is used to track metrics by name.
	MetricIDKey = "metric_name"

	// MetricTypeKey is used to track metrics by its type.
	MetricTypeKey = "metric_type"

	// MetricValueKey is used for debug purposes, represents metrics value.
	MetricValueKey = "metric_value"
)

type LogCtxProvider interface {
	LoggerCtx(ctx zerolog.Context) zerolog.Context
}

var _ LogCtxProvider = (LoggerCtxUpdate)(nil)

type LoggerCtxUpdate func(ctx zerolog.Context) zerolog.Context

func (upd LoggerCtxUpdate) LoggerCtx(ctx zerolog.Context) zerolog.Context {
	if upd != nil {
		return upd(ctx)
	}
	return ctx
}

func LogCtxUpdateWith(ctx zerolog.Context, providers ...LogCtxProvider) zerolog.Context {
	for _, p := range providers {
		ctx = p.LoggerCtx(ctx)
	}
	return ctx
}

func LogCtxFrom(providers ...LogCtxProvider) LoggerCtxUpdate {
	return func(ctx zerolog.Context) zerolog.Context {
		return LogCtxUpdateWith(ctx, providers...)
	}
}

func LogCtxKeyStr(key string, value string) LoggerCtxUpdate {
	return func(ctx zerolog.Context) zerolog.Context {
		return ctx.Str(key, value)
	}
}
