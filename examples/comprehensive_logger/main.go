package main

// Comprehensive logger example demonstrating various features
// including structured logging, error handling, and channel output.

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	logger "github.com/aN0mad/go-logo/logo"
)

var (
	logLocation = "./logs/comprehensive.log"
)

func main() {
	// Create a channel for log messages
	logChan := make(chan string, 100)

	// Start goroutine to process channel messages
	go func() {
		for msg := range logChan {
			fmt.Println("Channel received:", msg)
		}
	}()

	// Create a comprehensive logger with all features
	logger.Init(
		logger.SetLevel(slog.LevelDebug),
		logger.EnableStackTraces(),
		logger.AddSource(),
		logger.UseJSON(true), // Pretty JSON
		logger.AddFileOutput(logLocation, 10, 3, 30, true),
		logger.AddChannelOutput(logChan),
	)

	defer logger.Close() // Ensure all file writers are closed properly

	log := logger.L()

	log.Info("COMPREHENSIVE LOGGER EXAMPLE")
	log.Info("---------------------------")

	// Basic logging at different levels
	log.Trace("Trace message - most detailed level")
	log.Debug("Debug message - development information")
	log.Info("Info message - general operational messages")
	log.Warn("Warning message - potential issues")
	log.Error("Error message - actual problems")

	// Structured logging with rich context
	log.Info("Application startup",
		"version", "2.1.0",
		"environment", os.Getenv("APP_ENV"),
		"pid", os.Getpid(),
		"go_version", "1.18.3",
	)

	// Error logging with stacktrace
	err := generateError()
	if err != nil {
		log.Error("Failed operation",
			"error", err,
			"operation", "data_sync",
			"retry_count", 3,
		)
	}

	// Simulate periodic metrics logging
	log.Info("System metrics",
		"cpu_usage", 32.5,
		"memory_used_mb", 1024,
		"connections", 42,
		"timestamp", time.Now().Format(time.RFC3339),
	)

	log.Info("Comprehensive logging example completed")

	// Allow time for channel processing
	time.Sleep(200 * time.Millisecond)
}

func generateError() error {
	return fmt.Errorf("simulated error for demonstration purposes")
}
