package main

// This file shows the same functionality migrated to v2.

import (
	"time"

	logo "github.com/aN0mad/go-logo/v2"
)

func main() {
	// v2: Use CLI configuration (simplest approach)
	cfg := logo.NewConfig()
	cfg.Level = "debug"
	cfg.Console = true
	cfg.Source = true
	cfg.FilePath = "logs/app.log"
	cfg.RotationMaxSize = 10
	cfg.RotationMaxBackups = 3
	cfg.RotationMaxAge = 30
	cfg.RotationCompress = true

	logger, err := cfg.BuildLogger()
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	// v2: Use logger instance directly (no global L())
	logger.Info("Application started")

	// v2: Structured logging with Field constructors
	logger.Info("User login",
		logo.F("user_id", 12345),
		logo.F("username", "johndoe"),
		logo.F("ip_address", "192.168.1.100"),
		logo.F("timestamp", time.Now()),
	)

	// v2: Different log levels (same methods)
	logger.Trace("Trace level - very detailed")
	logger.Debug("Debug level - development info")
	logger.Info("Info level - general information")
	logger.Warn("Warn level - warning condition")
	logger.Error("Error level - error occurred")
	// logger.Fatal("Fatal level - critical error") // Would exit

	// v2: Complex structured data (same as v1)
	logger.Info("System metrics",
		logo.F("cpu", map[string]interface{}{
			"usage": 45.2,
			"temp":  62,
		}),
		logo.F("memory", map[string]interface{}{
			"used":  4096,
			"total": 8192,
		}),
	)

	// v2: Error with context
	err = processOrder(12345)
	if err != nil {
		logger.Error("Order processing failed",
			logo.F("order_id", 12345),
			logo.F("error", err.Error()),
			logo.F("retry_count", 3),
		)
	}
}

func processOrder(orderID int) error {
	// Simulate processing
	return nil
}
