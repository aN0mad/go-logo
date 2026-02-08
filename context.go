package logo

import (
	"context"
)

// contextKey is a private type for context keys to avoid collisions.
type contextKey string

const (
	loggerKey    contextKey = "logo.logger"
	requestIDKey contextKey = "logo.request_id"
)

// WithLogger adds a logger to the context.
// This allows passing a logger through the request lifecycle.
func WithLogger(ctx context.Context, logger *Logger) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, loggerKey, logger)
}

// FromContext retrieves the logger from the context.
// If no logger is found, returns a default logger with no handlers.
func FromContext(ctx context.Context) *Logger {
	if ctx == nil {
		return New()
	}

	if logger, ok := ctx.Value(loggerKey).(*Logger); ok {
		return logger
	}

	// Return a default logger if none is found
	return New()
}

// WithRequestID adds a request ID to the context.
// The request ID can be used for request tracing and correlation.
func WithRequestID(ctx context.Context, requestID string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, requestIDKey, requestID)
}

// RequestIDFromContext retrieves the request ID from the context.
// Returns an empty string if no request ID is found.
func RequestIDFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	if requestID, ok := ctx.Value(requestIDKey).(string); ok {
		return requestID
	}

	return ""
}

// ContextWithFields returns a logger from context with additional fields.
// This is a convenience method that combines FromContext and With.
func ContextWithFields(ctx context.Context, fields ...Field) *Logger {
	logger := FromContext(ctx)
	return logger.With(fields...)
}

// L is a shorthand for FromContext.
// Useful for quick logging from context.
func L(ctx context.Context) *Logger {
	return FromContext(ctx)
}
