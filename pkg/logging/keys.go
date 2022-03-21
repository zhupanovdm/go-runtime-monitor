package logging

import "github.com/rs/zerolog"

const (
	// CorrelationIDKey is used to track unique request.
	CorrelationIDKey = "cid"

	// CorrelationIDHeader is used to transport Correlation ID context value via the HTTP header.
	CorrelationIDHeader = "X-CorrelationID"

	// ServiceKey is used to track concrete service side effects.
	ServiceKey = "svc"

	// MetricIDKey is used to track metrics by name.
	MetricIDKey = "metric_name"

	// MetricTypeKey is used to track metrics by its type.
	MetricTypeKey = "metric_type"

	// MetricValueKey is used for debug purposes, represents metrics value.
	MetricValueKey = "metric_value"
)

var _ LogCtxProvider = (LoggerCtxUpdate)(nil)

type (
	// LogCtxProvider should be implemented by types that is supposed to display itself on logging context.
	LogCtxProvider interface {
		LoggerCtx(ctx zerolog.Context) zerolog.Context
	}

	// LoggerCtxUpdate function should update specified context and return updated logger context.
	LoggerCtxUpdate func(ctx zerolog.Context) zerolog.Context
)

// LoggerCtx applies LoggerCtxUpdate to specified logger context if not nil.
func (upd LoggerCtxUpdate) LoggerCtx(ctx zerolog.Context) zerolog.Context {
	if upd != nil {
		return upd(ctx)
	}
	return ctx
}

// LogCtxUpdateWith updates specified logger context sequentially applying provided contexts.
func LogCtxUpdateWith(ctx zerolog.Context, providers ...LogCtxProvider) zerolog.Context {
	for _, p := range providers {
		ctx = p.LoggerCtx(ctx)
	}
	return ctx
}

// LogCtxFrom is helper function combines several context providers to a single LoggerCtxUpdate.
func LogCtxFrom(providers ...LogCtxProvider) LoggerCtxUpdate {
	return func(ctx zerolog.Context) zerolog.Context {
		return LogCtxUpdateWith(ctx, providers...)
	}
}

// LogCtxKeyStr constructs new LoggerCtxUpdate that will update specified context with provided string key\value.
func LogCtxKeyStr(key string, value string) LoggerCtxUpdate {
	return func(ctx zerolog.Context) zerolog.Context {
		return ctx.Str(key, value)
	}
}
