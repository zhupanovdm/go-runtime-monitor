package logging

import "github.com/rs/zerolog"

const (
	// CorrelationIDKey is used to track unique request.
	CorrelationIDKey = "cid"

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

func UpdateLogCtxWith(ctx zerolog.Context, providers ...LogCtxProvider) zerolog.Context {
	for _, p := range providers {
		ctx = p.LoggerCtx(ctx)
	}
	return ctx
}

func LogCtxFrom(providers ...LogCtxProvider) func(ctx zerolog.Context) zerolog.Context {
	return func(ctx zerolog.Context) zerolog.Context {
		return UpdateLogCtxWith(ctx, providers...)
	}
}
