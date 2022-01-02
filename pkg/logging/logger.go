package logging

import (
	"context"
	"os"

	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"github.com/zhupanovdm/go-runtime-monitor/pkg/app"
)

const (
	// contextKeyLogger defines the key for the logger to be stored within request context.
	contextKeyLogger = app.ContextKey("Logger")

	// contextKeyCorrelationID defines the key for the correlation ID to be stored within request context.
	contextKeyCorrelationID = app.ContextKey("CorrelationID")
)

// Get returns a logger stored within the context.
// Creates a new one with correlation ID field and stores it in the context.
func Get(ctx context.Context) (context.Context, zerolog.Logger) {
	if value := ctx.Value(contextKeyLogger); value != nil {
		if logger, ok := value.(zerolog.Logger); ok {
			return ctx, logger
		}
	}

	correlationID, _ := uuid.NewUUID()
	logger := NewLogger().With().Str(CorrelationIDKey, correlationID.String()).Logger()
	ctx = context.WithValue(ctx, contextKeyCorrelationID, correlationID.String())
	return Set(ctx, logger), logger
}

// Set adds the logger to the context overwriting and existing one.
func Set(ctx context.Context, logger zerolog.Logger) context.Context {
	return context.WithValue(ctx, contextKeyLogger, logger)
}

func NewLogger() zerolog.Logger {
	return zerolog.New(os.Stdout).
		Output(zerolog.ConsoleWriter{Out: os.Stdout}).
		Level(zerolog.TraceLevel).
		With().
		Timestamp().
		Logger()
}
