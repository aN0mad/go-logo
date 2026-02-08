package main

// This example shows how to migrate from v1's JSON logging to v2.

import (
	"time"

	logo "github.com/aN0mad/go-logo/v2"
)

func main() {
	// v2: JSON logging with file output
	cfg := logo.NewConfig()
	cfg.Level = "debug"
	cfg.FileFormat = "json"
	cfg.Source = true
	cfg.FilePath = "logs/app-json.log"
	cfg.RotationMaxSize = 10
	cfg.RotationMaxBackups = 3
	cfg.RotationMaxAge = 30
	cfg.RotationCompress = true

	logger, err := cfg.BuildLogger()
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	logger.Info("Application started in JSON mode")

	// Structured data in JSON
	logger.Info("API request received",
		logo.F("method", "POST"),
		logo.F("path", "/api/v1/users"),
		logo.F("status", 201),
		logo.F("duration_ms", 45.2),
		logo.F("request_id", "req-abc-123"),
	)

	// Nested objects
	logger.Info("Database query executed",
		logo.F("query", "SELECT * FROM users WHERE id = ?"),
		logo.F("params", []interface{}{12345}),
		logo.F("result", map[string]interface{}{
			"rows_affected":     1,
			"execution_time_ms": 12.5,
		}),
	)

	// Arrays and complex structures
	logger.Info("Batch processing completed",
		logo.F("total_items", 100),
		logo.F("successful", 95),
		logo.F("failed", 5),
		logo.F("failed_ids", []int{12, 34, 56, 78, 90}),
		logo.F("processing_time", 2*time.Second),
	)
}
