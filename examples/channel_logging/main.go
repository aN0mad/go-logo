package main

// Channel-based logger example demonstrating how to use a channel for log messages.

import (
	"fmt"
	"log/slog"
	logger "logo/logo"
	"time"
)

var (
	logLocation = "./logs/channel_example.log"
)

func main() {
	// Create a channel for log messages
	logChan := make(chan string, 100)
	defer close(logChan) // Ensure channel is closed when done

	// Start goroutine to process log messages from channel
	go processLogMessages(logChan)

	// Initialize logger with channel output
	logger.Init(
		logger.SetLevel(slog.LevelDebug),
		logger.EnableTrace(),
		logger.AddChannelOutput(logChan),
		logger.AddFileOutput(logLocation, 10, 3, 30, true),
	)

	log := logger.L()

	log.Info("CHANNEL LOGGER EXAMPLES")
	log.Info("----------------------")

	// Log some messages that will be sent to the channel
	log.Info("This message is sent to a channel for custom processing")
	log.Debug("Debug information sent to channel")
	log.Warn("Warning sent to channel")

	// Simulate work to allow time for channel processing
	time.Sleep(100 * time.Millisecond)

	// Log with structured data
	log.Info("Structured data via channel",
		"process_id", 1234,
		"status", "running",
		"uptime_seconds", 3600,
	)

	// Give time for channel processing to complete
	time.Sleep(500 * time.Millisecond)
	log.Info("Channel logging example completed")
}

func processLogMessages(ch chan string) {
	fmt.Println("Log message processor started")
	for msg := range ch {
		// In a real application, you could:
		// - Send logs to an external monitoring system
		// - Parse and transform log data
		// - Filter logs based on content
		// - Trigger alerts based on log content

		fmt.Printf("CHANNEL RECEIVED: %s\n", msg)
	}

	// Allow some time for processing before exiting
	time.Sleep(time.Second * 2)
	fmt.Println("Log message processor stopped")
}
