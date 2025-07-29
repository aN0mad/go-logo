package main

// Context logger example demonstrating how to use context with logging.

import (
	"context"
	"log/slog"
	logger "logo/logo"
)

var (
	logLocation = "./logs/context_example.log"
)

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

	// Create a context with request ID
	ctx := context.WithValue(context.Background(), "request_id", "req-abc-123")

	// Use WithContext to get a logger with context values
	requestLogger := logger.WithContext(ctx)

	// The request ID will automatically be included in these log entries
	requestLogger.Info("Processing request")
	requestLogger.Debug("Request details being parsed")

	// Simulate request handling
	processRequest(ctx, requestLogger)

	log.Info("Context logging example completed")
}

func processRequest(ctx context.Context, log *logger.Logger) {
	// Extract request ID from context for demonstration
	reqID := ctx.Value("request_id")

	// Log with the request logger
	log.Info("Request processing complete", "status", "success")

	// Even though we're adding request_id explicitly here, with proper WithContext implementation,
	// it would be automatically included by the logger
	log.Debug("Request details",
		"request_id", reqID,
		"processing_time_ms", 42,
		"resource_path", "/api/v1/users",
	)
}
