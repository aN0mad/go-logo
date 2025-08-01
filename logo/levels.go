// Package logo provides functionality for structured logging.
//
// This file contains functions to query and manipulate the current log level
// in the logger configuration.
package logo

import (
	"io"
	"log/slog"
)

// IsLevelEnabled checks if a log level is enabled based on the current logger configuration.
// When used with the global logger, it checks against the global log level.
// When used with a specific logger instance, it checks against that logger's level.
//
// Parameters:
//   - level: The log level to check
//   - logger: Optional specific logger to check against (can be nil for global logger)
//
// Returns:
//   - bool: True if the level is enabled, false otherwise
func IsLevelEnabled(level slog.Level, logger *Logger) bool {
	if logger != nil && logger.ctx != nil {
		// Check against the specific logger's level
		return level >= logger.ctx.logLevel
	}

	// Fall back to global level
	mu.RLock()
	defer mu.RUnlock()
	return level >= LOGLEVEL
}

// GetCurrentLevel returns the log level configured for a logger.
// When used with the global logger, it returns the global log level.
// When used with a specific logger instance, it returns that logger's level.
//
// Parameters:
//   - logger: Optional specific logger to get level from (can be nil for global logger)
//
// Returns:
//   - slog.Level: The current log level
func GetCurrentLevel(logger *Logger) slog.Level {
	if logger != nil && logger.ctx != nil {
		// Return the specific logger's level
		return logger.ctx.logLevel
	}

	// Fall back to global level
	mu.RLock()
	defer mu.RUnlock()
	return LOGLEVEL
}

// SetGlobalLevel sets the log level for the global logger.
// This affects all subsequent log messages through the global logger (via L()).
//
// Parameters:
//   - level: The new log level to set
func SetGlobalLevel(level slog.Level) {
	mu.Lock()
	defer mu.Unlock()
	LOGLEVEL = level
}

// SetLoggerLevel sets the log level for a specific logger instance.
//
// Parameters:
//   - logger: The logger instance to configure
//   - level: The new log level to set
func SetLoggerLevel(logger *Logger, level slog.Level) {
	if logger != nil && logger.ctx != nil {
		// Update the context log level
		logger.ctx.logLevel = level

		// Update the handler's log level if it supports it
		if handler := logger.Handler(); handler != nil {
			// Try to update the level via LevelVar if available
			if leveler, ok := handler.(interface{ SetLevel(slog.Level) }); ok {
				leveler.SetLevel(level)
				return
			}

			// If the handler doesn't directly support SetLevel, we might need to replace it
			// Create a new handler with the updated level
			var newHandler slog.Handler

			// Configure handler options with the new level
			handlerOptions := &slog.HandlerOptions{
				Level:     level,
				AddSource: logger.ctx.includeSource,
				ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
					// Preserve any existing attribute replacement logic
					return a
				},
			}

			// Create a multi-writer for all outputs
			if len(logger.ctx.outputs) > 0 {
				multiwriter := io.MultiWriter(logger.ctx.outputs...)

				// Recreate the appropriate handler type
				if logger.ctx.useJSONFormat {
					newHandler = NewJSONHandler(multiwriter, handlerOptions, logger.ctx.jsonPretty)
				} else {
					newHandler = NewCustomTextHandler(multiwriter, handlerOptions)
				}

				// Replace the logger's handler
				logger.Logger = slog.New(newHandler)
			}
		}
	}
}
