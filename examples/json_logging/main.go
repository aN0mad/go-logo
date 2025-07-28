package main

// JSON logging example demonstrating how to use JSON format for structured logging.

import (
	"log/slog"
	logger "logo/logo"
	"time"
)

var (
	logLocation = "./logs/json_example.log"
)

func main() {
	// Initialize with JSON format
	logger.Init(
		logger.SetLevel(slog.LevelDebug),
		logger.EnableTrace(),
		logger.UseJSON(false), // Regular JSON (not pretty)
		logger.AddSource(),
		logger.AddFileOutput(logLocation, 10, 3, 30, true),
	)

	log := logger.L()

	log.Info("JSON FORMAT LOGGER EXAMPLES")
	log.Info("---------------------------")

	// Demonstrate all log levels
	log.Trace("This is a TRACE level message in JSON format")
	log.Debug("This is a DEBUG level message in JSON format")
	log.Info("This is an INFO level message in JSON format")
	log.Warn("This is a WARN level message in JSON format")
	log.Error("This is an ERROR level message in JSON format")

	// Demonstrate complex structured data
	log.Info("System metrics",
		"cpu", map[string]interface{}{
			"usage_percent": 45.2,
			"temperature_c": 62,
			"cores": []map[string]interface{}{
				{"id": 0, "usage": 30.5},
				{"id": 1, "usage": 60.2},
				{"id": 2, "usage": 45.0},
				{"id": 3, "usage": 50.1},
			},
		},
		"memory", map[string]interface{}{
			"total_mb": 16384,
			"used_mb":  8192,
			"free_mb":  8192,
		},
		"timestamp", time.Now(),
	)

	log.Info("JSON logging example completed")
}
