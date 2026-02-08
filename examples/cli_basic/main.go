package main

import (
	"flag"
	"fmt"

	logo "github.com/aN0mad/go-logo/v2"
)

func main() {
	// Create configuration
	cfg := logo.NewConfig()

	// Register CLI flags
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

	// Use the logger
	logger.Debug("This is a debug message")
	logger.Info("Application started", logo.F("version", "1.0.0"))
	logger.Warn("This is a warning", logo.F("threshold", 80))
	logger.Error("An error occurred", logo.F("code", 500))
}
