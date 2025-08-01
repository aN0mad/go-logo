// Package logo provides functionality for structured logging.
//
// This file contains the file writer implementation which supports log rotation
// through the lumberjack package.
package logo

import "github.com/aN0mad/lumberjack/v2"

// NewLumberjackWriter creates a new log writer with rotation capabilities.
// It configures a lumberjack logger with the specified parameters for log rotation.
//
// Parameters:
//   - filename: Path to the log file
//   - maxSize: Maximum size of log files in megabytes before rotation
//   - maxBackups: Maximum number of old log files to retain
//   - maxAge: Maximum number of days to retain old log files
//   - compress: Whether to compress rotated log files
//
// Returns:
//   - *lumberjack.Logger: Configured log writer with rotation support
func NewLumberjackWriter(filename string, maxSize, maxBackups, maxAge int, compress bool) *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackups,
		MaxAge:     maxAge,
		Compress:   compress,
	}
}
