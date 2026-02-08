package main

// This example shows v2's new CLI configuration feature
// which has no equivalent in v1.

import (
	"flag"

	logo "github.com/aN0mad/go-logo/v2"
)

func main() {
	// v2 NEW: CLI-driven configuration
	cfg := logo.NewConfig()
	cfg.RegisterFlags(nil) // Uses flag.CommandLine
	flag.Parse()

	// Build logger from CLI/env config
	logger, err := cfg.BuildLogger()
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	logger.Info("Application started with CLI config")
	logger.Debug("This will only show if --log-level=debug")
	logger.Info("Running with profile", logo.F("profile", "custom"))

	// Application logic here...
	doWork(logger)
}

func doWork(logger *logo.Logger) {
	logger.Info("Performing work", logo.F("task", "process_data"))
}

// Run this example with different configurations:
//
// Basic:
//   go run cli_config.go
//
// Set log level:
//   go run cli_config.go --log-level=debug
//
// Use JSON output:
//   go run cli_config.go --log-format=json
//
// Write to file:
//   go run cli_config.go --log-file=app.log
//
// Use dev profile:
//   go run cli_config.go --log-profile=dev
//
// Use prod profile with overrides:
//   go run cli_config.go --log-profile=prod --log-level=debug
//
// Environment variables (alternative to flags):
//   LOG_LEVEL=debug LOG_FORMAT=json go run cli_config.go
