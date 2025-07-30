package main

import (
	"bytes"
	"context"
	"fmt"
	logger "logo/logo"
)

func main() {
	// Test the context logger in a controlled environment
	var buf bytes.Buffer

	// Initialize logger with the same settings as in your test
	logger.Init(
		logger.SetConsoleOutput(&buf),
		logger.DisableColors(),
		logger.AddSource(),
	)

	// Create context with request_id
	ctx := context.WithValue(context.Background(), "request_id", "test-id-123")

	// Create context logger similar to example
	loggerFromContext := func(ctx context.Context) *logger.Logger {
		baseLogger := logger.L()

		// Add application-specific context values
		if requestID, ok := ctx.Value("request_id").(string); ok {
			baseLogger = &logger.Logger{
				Logger: baseLogger.Logger.With("request_id", requestID),
			}
		}

		return baseLogger
	}

	// Get logger with context and log a message
	logWithCtx := loggerFromContext(ctx)
	logWithCtx.Info("test message")

	// Print the raw output to examine it
	fmt.Println("--- Test buffer output ---")
	fmt.Printf("Raw output: %q\n", buf.String())

	// Print for comparison what the direct console output looks like
	fmt.Println("\n--- Direct console output ---")
	loggerFromContext(ctx).Info("direct test message")
}
