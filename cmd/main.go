package main

import (
	"encoding/json"
	"fmt"
	logger "logo/logo"
	"time"
)

func main() {

	// Console-only logger (uses charmbracelet/log for color output)
	log := logger.New(nil, true) // Initialize logger with console output enabled

	log.Info("📦 Starting the CLI tool...")

	log.Debug("Initializing subsystems...")
	time.Sleep(500 * time.Millisecond)

	log.Info("Processing data stream...")
	for i := 1; i <= 5; i++ {
		log.Debug("Processing item %d", i)
		time.Sleep(250 * time.Millisecond)

		if i == 3 {
			log.Warn("Potential issue with item %d", i)
		}
	}

	log.Info("✔️ Completed all tasks")
	fmt.Printf("\n")

	// End console logging

	// start channel-based logger

	log.Info("📦 Starting channel-based logger...")

	logChan := make(chan logger.LogMessage, 100) // Channel for log messages
	defer close(logChan)                         // Ensure channel is closed when done

	// Start a goroutine to listen for log messages
	// and print them to the console
	// This simulates a logging system that could be used in a larger application
	go func() {
		for {
			select {
			case msg := <-logChan:
				data, _ := json.Marshal(msg)
				fmt.Printf("data: %s\n", data)

			}
		}
	}()

	appLog := logger.New(logChan, false) // Initialize logger with channel
	simulateAppLogging(appLog)
	log.Info("✔️ Stopping channel-based logger...")

}

// simulateAppLogging periodically logs events
func simulateAppLogging(appLog *logger.Logger) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	count := 0
	for range ticker.C {
		count++
		appLog.Info("Heartbeat %d", count)
		if count == 5 {
			appLog.Warn("Slow operation detected")
		}
		if count == 8 {
			appLog.Error("Something went wrong")
		}
		if count == 10 {
			appLog.Info("Shutting down logging loop")
			return
		}
	}
}
