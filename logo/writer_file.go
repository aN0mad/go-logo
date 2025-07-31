// Package logo provides functionality for structured logging.
//
// This file contains the file writer implementation which supports log rotation
// through the lumberjack package.
package logo

import "github.com/aN0mad/lumberjack/v2"

// NewLumberjackWriter creates a new writer that writes to a file with rotation support.
// It uses the lumberjack library to provide log rotation capabilities.
//
// Parameters:
//   - filepath: The path to the log file
//   - maxSizeMB: Maximum size of the log file in megabytes before rotation
//   - maxBackups: Maximum number of old log files to retain
//   - maxAgeDays: Maximum number of days to retain old log files
//   - compress: If true, rotated log files will be compressed using gzip
//
// Returns:
//   - *lumberjack.Logger: A configured lumberjack logger that implements io.Writer
func NewLumberjackWriter(filepath string, maxSizeMB, maxBackups, maxAgeDays int, compress bool) *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   filepath,
		MaxSize:    maxSizeMB,
		MaxBackups: maxBackups,
		MaxAge:     maxAgeDays,
		Compress:   compress,
	}
}
