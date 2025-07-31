package main

import (
	"fmt"
	"log/slog"

	logger "github.com/aN0mad/go-logo/logo"
)

func main() {

	logger.Init(
		logger.AddSource(),                                    // Log source caller information
		logger.SetLevel(slog.LevelDebug),                      // Set the minimum log level to Debug
		logger.AddFileOutput("logs/app.log", 10, 3, 30, true), // Log to file with rotation
	)

	defer logger.Close() // Ensure all file writers are closed properly

	log := logger.L()
	fmt.Println("Initialized logger")

	log.Trace("This is a trace message - It will not show up")
	log.Debug("This is a debug message")
	log.Info("This is an info message")
	log.Warn("This is a warning message")
	log.Error("This is an error message")
	log.Fatal("This is a fatal message - It will exit the program")
}
