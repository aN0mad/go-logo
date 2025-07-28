package main

// Text logging example demonstrating how to use text format for structured logging.

import (
	"log/slog"
	logger "logo/logo"
	"time"
)

var (
	logLocation = "./logs/text_example.log"
)

func main() {
	// Initialize with default text format
	logger.Init(
		logger.SetLevel(slog.LevelDebug),
		logger.EnableTrace(),
		logger.AddSource(),
		logger.AddFileOutput(logLocation, 10, 3, 30, true),
	)

	log := logger.L()

	log.Info("TEXT FORMAT LOGGER EXAMPLES")
	log.Info("---------------------------")

	// Demonstrate all log levels
	log.Trace("This is a TRACE level message - for very detailed diagnostic information")
	log.Debug("This is a DEBUG level message - for debugging and development")
	log.Info("This is an INFO level message - general information about system operation")
	log.Warn("This is a WARN level message - warning conditions that should be addressed")
	log.Error("This is an ERROR level message - error conditions that might be recoverable")
	// Note: Fatal would terminate the program, so we're commenting it out
	// log.Fatal("This is a FATAL level message - critical errors causing program termination")

	// Demonstrate structured logging with attributes
	log.Info("User logged in",
		"user_id", 12345,
		"username", "example_user",
		"role", "admin",
		"login_time", time.Now(),
	)

	// Demonstrate error logging with additional context
	err := simulateError()
	if err != nil {
		log.Error("Failed to process request",
			"error", err,
			"request_id", "req-123456",
			"client_ip", "192.168.1.1",
		)
	}

	log.Info("Text logging example completed")
}

func simulateError() error {
	return &customError{message: "simulated error for demonstration"}
}

type customError struct {
	message string
}

func (e *customError) Error() string {
	return e.message
}
