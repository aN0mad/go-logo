package main

import (
	"time"

	"github.com/aN0mad/go-logo/v2"
)

func main() {
	// Create a logger with console handler (pretty formatting + auto color detection)
	log := logo.New(logo.WithConsole())

	log.Info("Application starting...")

	// Different log levels with colors
	log.Debug("Debug message - detailed diagnostic info")
	log.Info("Info message - general information")
	log.Warn("Warning message - something needs attention")
	log.Error("Error message - something went wrong")

	// Structured logging with fields
	log.Info("User logged in",
		logo.F("user_id", 12345),
		logo.F("username", "john.doe"),
		logo.F("ip", "192.168.1.100"),
	)

	// Nested context with child loggers
	requestLog := log.With(
		logo.F("request_id", "req-abc-123"),
		logo.F("method", "POST"),
		logo.F("path", "/api/users"),
	)

	requestLog.Info("Processing request")
	requestLog.Info("Database query executed", logo.F("duration", 45*time.Millisecond))
	requestLog.Info("Request completed", logo.F("status", 200))

	// Fields with special characters
	log.Info("Processing file",
		logo.F("filename", "my document.pdf"),
		logo.F("size_mb", 2.5),
	)

	// Multiple fields showing different types
	log.Info("System metrics",
		logo.F("cpu_percent", 45.2),
		logo.F("memory_mb", 1024),
		logo.F("uptime", 24*time.Hour),
		logo.F("healthy", true),
	)

	log.Info("Application shutdown complete")
}
