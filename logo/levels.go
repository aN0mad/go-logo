// Package logo provides functionality for structured logging.
//
// This file contains functions to query and manipulate the current log level
// in the logger configuration.
package logo

import "log/slog"

// IsLevelEnabled checks if a log level is enabled based on the current logger configuration.
//
// Parameters:
//   - level: The log level to check
//
// Returns:
//   - bool: True if the level is enabled, false otherwise
func IsLevelEnabled(level slog.Level) bool {
	return level >= logLevel
}

// GetCurrentLevel returns the current minimum log level configured in the logger.
//
// Returns:
//   - slog.Level: The current log level
func GetCurrentLevel() slog.Level {
	return logLevel
}
