package main

// Context logger example demonstrating how to use context with logging.

import (
	"context"
	"log/slog"

	logger "github.com/aN0mad/go-logo/logo"
)

var (
	logLocation = "./logs/context_example.log"
)

// LoggerFromContext creates a logger with context values.
// This is an application-level extension of the core logger.
func LoggerFromContext(ctx context.Context) *logger.Logger {
	baseLogger := logger.L()

	// Add application-specific context values
	if requestID, ok := ctx.Value("request_id").(string); ok {
		// The fix is here - we need to create a new logger.Logger wrapper around
		// the slog.Logger that has the additional attribute
		baseLogger = &logger.Logger{
			Logger: baseLogger.Logger.With("request_id", requestID),
		}
	}

	if userID, ok := ctx.Value("user_id").(string); ok {
		baseLogger = &logger.Logger{
			Logger: baseLogger.Logger.With("user_id", userID),
		}
	}

	// Add other context values as needed for your application
	return baseLogger
}

func main() {
	// Initialize logger
	logger.Init(
		logger.SetLevel(slog.LevelDebug),
		logger.AddSource(),
		logger.AddFileOutput(logLocation, 10, 3, 30, true),
	)

	defer logger.Close() // Ensure all file writers are closed properly

	log := logger.L()

	log.Info("CONTEXT LOGGER EXAMPLES")
	log.Info("----------------------")

	// Create a context with request ID and user ID
	ctx := context.WithValue(context.Background(), "request_id", "req-abc-123")
	ctx = context.WithValue(ctx, "user_id", "user-456")

	// Use LoggerFromContext instead of WithContext
	requestLogger := LoggerFromContext(ctx)

	// The request ID and user ID will automatically be included in these log entries
	requestLogger.Info("Processing request")
	requestLogger.Debug("Request details being parsed")

	// Simulate request handling
	processRequest(ctx, requestLogger)

	log.Info("Context logging example completed")
}

func processRequest(ctx context.Context, log *logger.Logger) {
	// No need to extract request ID from context - it's already in the logger

	// Log with the context-aware logger
	log.Info("Request processing complete", "status", "success")

	// Additional attributes are merged with the context attributes
	log.Debug("Request details",
		"processing_time_ms", 42,
		"resource_path", "/api/v1/users",
	)

	// Create a sub-context for a specific operation
	operationCtx := context.WithValue(ctx, "operation_id", "op-789")

	// Create a new logger for this operation
	opLogger := LoggerFromContext(operationCtx)

	// This log will include request_id, user_id, and operation_id
	opLogger.Info("Operation completed", "result", "success")
}
