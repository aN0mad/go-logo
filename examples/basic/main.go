package main

import (
	"os"

	"github.com/aN0mad/go-logo/v2"
)

func main() {
	// Create a text formatter
	formatter := logo.NewTextFormatter()

	// Create a console handler that outputs to stdout
	consoleHandler := logo.NewWriterHandler(os.Stdout, formatter, logo.LevelDebug)

	// Create a router with the console handler
	router := logo.NewRouter(consoleHandler)

	// Create a logger with the router
	log := logo.New(
		logo.WithRouter(router),
		logo.WithName("example"),
	)

	// Log at different levels
	log.Debug("This is a debug message")
	log.Info("Application started", logo.F("version", "1.0.0"))
	log.Warn("This is a warning", logo.F("component", "main"))
	log.Error("This is an error", logo.F("error_code", 500))

	// Create a child logger with additional context
	requestLogger := log.With(
		logo.F("request_id", "abc-123"),
		logo.F("user", "john.doe"),
	)

	requestLogger.Info("Processing request")
	requestLogger.Info("Request completed", logo.F("duration_ms", 42))

	// Clean up
	log.Close()
}
