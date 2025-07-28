// logo/writer_file.go
package logger

import "gopkg.in/natefinch/lumberjack.v2"

// NewLumberjackWriter creates a new lumberjack.Logger instance for file logging.
func NewLumberjackWriter(filename string, maxSize, backups, maxAge int, compress bool) *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: backups,
		MaxAge:     maxAge,
		Compress:   compress,
	}
}
