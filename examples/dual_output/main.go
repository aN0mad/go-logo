package main

import (
"fmt"
"os"

"github.com/aN0mad/go-logo/v2"
)

func main() {
	// Create console handler for human-readable output
	consoleHandler := logo.NewConsoleHandler()

	// Create JSON handler for structured logging
	jsonHandler := logo.NewJSONHandler(os.Stderr, logo.LevelInfo, logo.WithPrettyJSON(true))

	// Router sends to both handlers
	router := logo.NewRouter(consoleHandler, jsonHandler)
	log := logo.New(logo.WithRouter(router))

	fmt.Println("=== Dual Output Demo ===")
	fmt.Println("Console (stdout) gets pretty colored logs")
	fmt.Println("JSON (stderr) gets structured JSON")
	fmt.Println()

	log.Info("application started", logo.F("version", "2.0.0"))
	log.Warn("memory usage high", logo.F("percent", 85))
	log.Error("connection lost", logo.F("retry", 3))

	log.Info("shutdown complete")
}
