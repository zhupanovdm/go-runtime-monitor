package logging

import (
	"context"
	"os"

	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"github.com/zhupanovdm/go-runtime-monitor/pkg"
	"github.com/zhupanovdm/go-runtime-monitor/pkg/app"
)

const (
	// ctxKeyLogger identifies logger instance bound within request context.
	ctxKeyLogger = app.CtxKey("Logger")

	// ctxKeyCorrelationID identifies request's correlation ID.
	ctxKeyCorrelationID = app.CtxKey("CorrelationID")
)

// Option specifies logger functional option.
type Option func(zerolog.Logger) zerolog.Logger

// GetOrCreateLogger returns context bound logger.
// Creates a new one with correlation ID field than binds it to context.
func GetOrCreateLogger(ctx context.Context, options ...Option) (context.Context, zerolog.Logger) {
	if ctx == nil {
		ctx = context.Background()
	}
	if value := ctx.Value(ctxKeyLogger); value != nil {
		if logger, ok := value.(zerolog.Logger); ok {
			return ctx, ApplyOptions(logger, options...)
		}
	}

	logger := NewLogger(options...)
	return SetLogger(ctx, logger), logger
}

// WithService option adds service log key corresponding to specified service.
func WithService(service pkg.Service) Option {
	return WithServiceName(service.Name())
}

// WithServiceName option adds service log key corresponding to specified service name.
func WithServiceName(service string) Option {
	return func(logger zerolog.Logger) zerolog.Logger {
		return logger.With().Str(ServiceKey, service).Logger()
	}
}

// WithCID option extracts CID from specified context and adds it to log context.
func WithCID(ctx context.Context) Option {
	return func(logger zerolog.Logger) zerolog.Logger {
		if value := ctx.Value(ctxKeyCorrelationID); value != nil {
			if correlationID, ok := value.(string); ok {
				return logger.With().Str(CorrelationIDKey, correlationID).Logger()
			}
		}
		return logger
	}
}

// SetIfAbsentCID stores specified CID within context if CID was not previously set.
func SetIfAbsentCID(ctx context.Context, cid string) (context.Context, string) {
	if value := ctx.Value(ctxKeyCorrelationID); value != nil {
		if cid, ok := value.(string); ok {
			return ctx, cid
		}
	}
	return SetCID(ctx, cid)
}

// SetCID stores specified CID within context. Previously stored CID value will be overridden.
func SetCID(ctx context.Context, cid string) (context.Context, string) {
	return context.WithValue(ctx, ctxKeyCorrelationID, cid), cid
}

// NewCID generates new unique CID.
func NewCID() string {
	cid, _ := uuid.NewUUID()
	return cid.String()
}

// SetLogger binds specified logger to the context.
func SetLogger(ctx context.Context, logger zerolog.Logger) context.Context {
	return context.WithValue(ctx, ctxKeyLogger, logger)
}

// NewLogger creates new logger instance with specified functional options.
func NewLogger(options ...Option) zerolog.Logger {
	logger := zerolog.New(os.Stdout).
		Output(zerolog.ConsoleWriter{Out: os.Stdout}).
		Level(zerolog.TraceLevel).
		With().
		Timestamp().
		Logger()
	return ApplyOptions(logger, options...)
}

// ApplyOptions applies given functional options to specified logger.
func ApplyOptions(logger zerolog.Logger, options ...Option) zerolog.Logger {
	for _, opt := range options {
		logger = opt(logger)
	}
	return logger
}
