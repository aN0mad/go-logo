package main

import (
"fmt"
"time"

"github.com/aN0mad/go-logo/v2"
)

func main() {
	// Create logger with file output and rotation
	log := logo.New(logo.WithJSONFile("app.log"))
	defer log.Close()

	fmt.Println("Logging to app.log with default rotation settings...")

	// Log various messages
	log.Info("application started", logo.F("version", "2.0.0"))
	log.Debug("initialization complete")

	// Simulate some activity
	for i := 0; i < 10; i++ {
		log.Info("processing item",
logo.F("item_id", i),
logo.F("timestamp", time.Now().Unix()),
		)
	}

	log.Warn("approaching memory limit", logo.F("usage_percent", 85))
	log.Error("failed to connect", logo.F("host", "db.example.com"))

	log.Info("application shutdown")

	fmt.Println("Done! Check app.log for output")
}
