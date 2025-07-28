package main

// Pretty JSON logging example demonstrating how to use pretty JSON format for structured logging.

import (
	"log/slog"
	logger "logo/logo"
	"time"
)

var (
	logLocation = "./logs/pretty_json_example.log"
)

func main() {
	// Initialize with pretty JSON format
	logger.Init(
		logger.SetLevel(slog.LevelDebug),
		logger.EnableTrace(),
		logger.UseJSON(true), // Pretty JSON
		logger.AddSource(),
		logger.AddFileOutput(logLocation, 10, 3, 30, true),
	)

	log := logger.L()

	log.Info("PRETTY JSON FORMAT LOGGER EXAMPLES")
	log.Info("----------------------------------")

	// Demonstrate nested data structure that benefits from pretty printing
	log.Info("Application configuration",
		"app", map[string]interface{}{
			"name":        "Example Service",
			"version":     "1.2.3",
			"environment": "development",
			"features": map[string]bool{
				"authentication": true,
				"authorization":  true,
				"caching":        false,
				"metrics":        true,
			},
			"database": map[string]interface{}{
				"host": "localhost",
				"port": 5432,
				"credentials": map[string]string{
					"user":     "app_user",
					"password": "******",
				},
				"pool": map[string]int{
					"max_connections":      100,
					"min_connections":      5,
					"idle_timeout_seconds": 300,
				},
			},
			"api": map[string]interface{}{
				"rate_limit": 1000,
				"timeout_ms": 5000,
				"endpoints": []string{
					"/api/v1/users",
					"/api/v1/products",
					"/api/v1/orders",
				},
			},
		},
		"timestamp", time.Now(),
	)

	log.Info("Pretty JSON logging example completed")
}
