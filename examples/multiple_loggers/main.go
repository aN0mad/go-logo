package main

// Multiple loggers example demonstrating how to use multiple independent logger
// instances with different configurations in the same application.

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	logger "github.com/aN0mad/go-logo/logo"
)

var (
	appLogLocation    = "./logs/multiple_loggers_app_example.log"
	accessLogLocation = "./logs/multiple_loggers_access_example.log"
	errorLogLocation  = "./logs/multiple_loggers_error_example.log"
)

func main() {
	// Ensure logs directory exists
	if err := os.MkdirAll("./logs", 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating logs directory: %v\n", err)
		os.Exit(1)
	}

	// Create an application logger with info level that logs to console with colors
	// and to a rotated file
	appLogger := logger.NewLogger(
		logger.SetLevel(slog.LevelInfo),
		logger.AddSource(),                                     // Include source information
		logger.AddFileOutput(appLogLocation, 10, 3, 30, false), // 10MB files, 3 backups, 30 days
	)
	defer appLogger.Close() // Properly close the logger at exit

	// Create an access logger with debug level that logs in JSON format to a separate file
	accessLogger := logger.NewLogger(
		logger.SetLevel(slog.LevelDebug),
		logger.UseJSON(true), // Enable JSON format
		logger.DisableColors(),
		logger.AddFileOutput(accessLogLocation, 20, 5, 7, true), // 20MB files, 5 backups, 7 days, compressed
	)
	defer accessLogger.Close()

	// Create an error logger with error level that logs to both console (with colors)
	// and a dedicated error log file
	errorLogger := logger.NewLogger(
		logger.SetLevel(slog.LevelError),
		logger.EnableStackTraces(),                               // Include stack traces for errors
		logger.AddFileOutput(errorLogLocation, 10, 10, 90, true), // 10MB files, 10 backups, 90 days
	)
	defer errorLogger.Close()

	appLogger.Info("MULTIPLE LOGGERS EXAMPLE")
	appLogger.Info("-----------------------")

	// Use the loggers for different purposes
	appLogger.Info("Application started", "version", "1.0.0")

	// Simulate some application activity
	for i := 0; i < 3; i++ {
		// Log an access record
		accessLogger.Debug("Request received",
			"method", "GET",
			"path", "/api/users",
			"client_ip", "192.168.1.100",
			"user_agent", "Mozilla/5.0",
		)

		// Log application events
		appLogger.Info("Processing request", "request_id", fmt.Sprintf("req-%d", i))

		// Simulate an occasional error
		if i == 1 {
			err := fmt.Errorf("database connection timeout")
			errorLogger.Error("Failed to process request",
				"request_id", fmt.Sprintf("req-%d", i),
				"error", err,
			)

			// The application logger still gets info about the error
			appLogger.Warn("Request processing failed", "request_id", fmt.Sprintf("req-%d", i))
		}

		// Log successful completion with the access logger
		accessLogger.Info("Request completed",
			"method", "GET",
			"path", "/api/users",
			"status", 200,
			"duration_ms", 30,
		)

		time.Sleep(100 * time.Millisecond)
	}

	// Change log levels dynamically at runtime
	logger.SetLoggerLevel(appLogger, slog.LevelDebug)
	appLogger.Info("Log level changed to DEBUG")

	// Now debug messages will appear
	appLogger.Debug("This debug message will now be visible")

	appLogger.Info("Multiple loggers example completed")
	fmt.Println("Logs have been written to:")
	fmt.Printf("- Application log: %s\n", appLogLocation)
	fmt.Printf("- Access log (JSON): %s\n", accessLogLocation)
	fmt.Printf("- Error log: %s\n", errorLogLocation)
}
