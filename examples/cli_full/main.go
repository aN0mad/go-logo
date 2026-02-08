package main

import (
	"flag"
	"fmt"
	"time"

	logo "github.com/aN0mad/go-logo/v2"
)

func main() {
	// Create configuration
	cfg := logo.NewConfig()

	// Load environment variables first
	cfg.LoadFromEnv()

	// Register CLI flags (will override env vars)
	cfg.RegisterFlags(nil)

	// Parse command line flags
	flag.Parse()

	// Build logger from the configuration
	logger, err := cfg.BuildLogger()
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		return
	}
	defer logger.Close()

	// Demo logging at all levels
	logger.Trace("Trace level logging - very detailed", logo.F("step", "initialization"))
	logger.Debug("Debug level logging", logo.F("config", "loaded"))
	logger.Info("Application started successfully",
		logo.F("version", "2.0.0"),
		logo.F("profile", cfg.Profile))

	// Simulate some work
	for i := 1; i <= 3; i++ {
		logger.Info("Processing item",
			logo.F("item", i),
			logo.F("total", 3))
		time.Sleep(100 * time.Millisecond)
	}

	logger.Warn("Resource usage high",
		logo.F("cpu", 85.2),
		logo.F("memory", 78.9))

	logger.Error("Failed to connect to service",
		logo.F("service", "database"),
		logo.F("retry", 3),
		logo.F("timeout", "5s"))

	logger.Fatal("Critical system failure - application terminating",
		logo.F("reason", "out of memory"))

	logger.Info("Application shutdown complete")
}
