package main

import (
"fmt"
"os"

"github.com/aN0mad/go-logo/v2"
)

func main() {
	// Custom rotation configuration
	rotationConfig := &logo.RotationConfig{
		MaxSize:    10,    // 10 MB per file
		MaxAge:     7,     // Keep logs for 7 days
		MaxBackups: 5,     // Keep 5 old files
		Compress:   true,  // Compress rotated files
		LocalTime:  false, // Use UTC for timestamps
	}

	// Create file handler with custom rotation
	formatter := logo.NewJSONFormatter()
	fileHandler, err := logo.NewFileHandler("logs/app.log", formatter,
logo.WithRotation(rotationConfig),
logo.WithFileLevel(logo.LevelInfo), // Only log INFO and above to file
)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create file handler: %v\n", err)
		os.Exit(1)
	}
	defer fileHandler.Close()

	// Also add console handler for immediate feedback
	consoleHandler := logo.NewConsoleHandler()

	// Create logger with both handlers
	router := logo.NewRouter(consoleHandler, fileHandler)
	log := logo.New(logo.WithRouter(router))

	fmt.Println("=== File Rotation Demo ===")
	fmt.Println("Rotation Settings:")
	fmt.Printf("  MaxSize: %d MB\n", rotationConfig.MaxSize)
	fmt.Printf("  MaxAge: %d days\n", rotationConfig.MaxAge)
	fmt.Printf("  MaxBackups: %d\n", rotationConfig.MaxBackups)
	fmt.Printf("  Compress: %v\n", rotationConfig.Compress)
	fmt.Println()

	// Log messages (both to console and file)
	log.Info("application started with rotation enabled")
	log.Debug("this debug message only goes to console")
	log.Info("processing data batch", logo.F("batch_id", 1))
	log.Warn("disk usage high", logo.F("percent", 90))
	log.Error("backup failed", logo.F("error", "connection timeout"))

	log.Info("application shutdown")

	fmt.Println("\nLogs written to logs/app.log")
	fmt.Println("Check logs/ directory for rotated files")
}
