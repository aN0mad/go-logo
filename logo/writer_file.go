// Package logger provides functionality for structured logging.
//
// This file contains the file writer implementation which supports log rotation
// through the lumberjack package.
package logger

import "github.com/aN0mad/lumberjack/v2"

// NewLumberjackWriter creates a new lumberjack.Logger instance for file logging
// with rotation capabilities.
//
// Parameters:
//   - filename: Path to the log file
//   - maxSize: Maximum size of the log file in megabytes before it's rotated
//   - backups: Maximum number of old log files to retain
//   - maxAge: Maximum number of days to retain old log files
//   - compress: Whether to compress old log files
//
// Returns a configured lumberjack.Logger that implements io.Writer
func NewLumberjackWriter(filename string, maxSize, backups, maxAge int, compress bool) *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: backups,
		MaxAge:     maxAge,
		Compress:   compress,
	}
}
