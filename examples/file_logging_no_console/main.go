package main

// File logging example demonstrating how to use file output without console logging.

import (
	"fmt"
	"log/slog"
	logger "logo/logo"
	"time"
)

var (
	logLocation = "./logs/file_example.log"
)

func main() {
	// Initialize logger with file output only (no console)
	logger.Init(
		logger.SetLevel(slog.LevelDebug),
		logger.EnableStackTraces(),
		logger.DisableConsole(), // Disable console output
		logger.AddFileOutput(logLocation, 10, 3, 30, true),
	)

	defer logger.Close() // Ensure all file writers are closed properly

	log := logger.L()

	log.Info("FILE LOGGER EXAMPLES")
	log.Info("------------------")

	// Log various messages at different levels
	log.Debug("Debug message written to file only")
	log.Info("Info message written to file only")
	log.Warn("Warning message written to file only")

	// Log with timestamps and duration
	startTime := time.Now()
	// Simulate some work
	time.Sleep(100 * time.Millisecond)
	endTime := time.Now()
	duration := endTime.Sub(startTime)

	log.Info("Operation completed",
		"start_time", startTime,
		"end_time", endTime,
		"duration_ms", duration.Milliseconds(),
	)

	// Tell user where to find the logs
	fmt.Println("Logs have been written to /workspace/logs/file_example.log")
	log.Info("File logging example completed")
}
