package main

import (
	"fmt"
	"log/slog"
	logger "logo/logo"
)

func main() {
	logChan := make(chan string, 100)
	defer close(logChan) // Ensure channel is closed when done

	logger.Init(
		logger.AddSource(),
		// logger.EnableTrace(),
		logger.SetLevel(slog.LevelInfo),
		// logger.DisableColors(),
		logger.AddFileOutput("logs/app.log", 10, 3, 30, true),
	)

	log := logger.L()
	fmt.Println("Initialized logger")

	log.Trace("This is a trace message")
	log.Debug("This is a debug message")
	log.Info("This is an info message")
	log.Warn("This is a warning message")
	log.Error("This is an error message")
	log.Fatal("This is a fatal message")
}
