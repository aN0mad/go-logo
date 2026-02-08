package main

import (
	"flag"
	"fmt"

	logo "github.com/aN0mad/go-logo/v2"
)

func main() {
	// Create configuration
	cfg := logo.NewConfig()

	// Load from environment first
	cfg.LoadFromEnv()

	// Register CLI flags (flags override environment)
	cfg.RegisterFlags(nil)

	// Parse flags
	flag.Parse()

	// Build logger from config
	logger, err := cfg.BuildLogger()
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		return
	}
	defer logger.Close()

	// Log some messages
	logger.Debug("Debug message - only visible in dev mode")
	logger.Info("Application starting", logo.F("profile", cfg.Profile))
	logger.Warn("Warning message")
	logger.Error("Error message", logo.F("details", "something went wrong"))
}
