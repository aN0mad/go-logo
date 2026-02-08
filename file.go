package logo

import (
	"os"

	"github.com/aN0mad/lumberjack/v2"
)

// RotationConfig configures log file rotation behavior.
type RotationConfig struct {
	// MaxSize is the maximum size in megabytes of the log file before it gets rotated.
	// Default is 100 MB.
	MaxSize int

	// MaxAge is the maximum number of days to retain old log files based on the
	// timestamp encoded in their filename. The default is not to remove old log files.
	MaxAge int

	// MaxBackups is the maximum number of old log files to retain. The default
	// is to retain all old log files (though MaxAge may still cause them to get deleted.)
	MaxBackups int

	// Compress determines if the rotated log files should be compressed using gzip.
	// The default is not to perform compression.
	Compress bool

	// LocalTime determines if the time used for formatting the timestamps in
	// backup files is the computer's local time. The default is to use UTC time.
	LocalTime bool
}

// DefaultRotationConfig returns sensible defaults for log rotation.
func DefaultRotationConfig() *RotationConfig {
	return &RotationConfig{
		MaxSize:    100, // 100 MB
		MaxAge:     30,  // 30 days
		MaxBackups: 10,  // keep 10 old files
		Compress:   true,
		LocalTime:  false,
	}
}

// FileHandler is a handler that writes logs to a file with rotation support.
type FileHandler struct {
	*WriterHandler
	logger *lumberjack.Logger
}

// FileHandlerOption configures a FileHandler
type FileHandlerOption func(*FileHandler)

// WithRotation sets the rotation configuration for the file handler.
func WithRotation(config *RotationConfig) FileHandlerOption {
	return func(h *FileHandler) {
		if config != nil {
			h.logger.MaxSize = config.MaxSize
			h.logger.MaxAge = config.MaxAge
			h.logger.MaxBackups = config.MaxBackups
			h.logger.Compress = config.Compress
			h.logger.LocalTime = config.LocalTime
		}
	}
}

// WithFileLevel sets the minimum level for the file handler
func WithFileLevel(level Level) FileHandlerOption {
	return func(h *FileHandler) {
		h.BaseHandler.minLevel = level
	}
}

// NewFileHandler creates a new file handler with optional rotation.
// The file handler writes JSON-formatted logs to the specified file path.
func NewFileHandler(path string, formatter Formatter, opts ...FileHandlerOption) (*FileHandler, error) {
	// Create lumberjack logger
	ljLogger := &lumberjack.Logger{
		Filename: path,
	}

	// Apply default rotation config
	defaultConfig := DefaultRotationConfig()
	ljLogger.MaxSize = defaultConfig.MaxSize
	ljLogger.MaxAge = defaultConfig.MaxAge
	ljLogger.MaxBackups = defaultConfig.MaxBackups
	ljLogger.Compress = defaultConfig.Compress
	ljLogger.LocalTime = defaultConfig.LocalTime

	// Create writer handler with the lumberjack logger
	writerHandler := NewWriterHandler(ljLogger, formatter, LevelDebug)

	handler := &FileHandler{
		WriterHandler: writerHandler,
		logger:        ljLogger,
	}

	// Apply options
	for _, opt := range opts {
		opt(handler)
	}

	return handler, nil
}

// Close closes the file handler and flushes any remaining data.
func (h *FileHandler) Close() error {
	if h.logger != nil {
		return h.logger.Close()
	}
	return nil
}

// Rotate causes the file handler to close the existing log file and immediately
// create a new one. This is useful for manual rotation triggers.
func (h *FileHandler) Rotate() error {
	if h.logger != nil {
		return h.logger.Rotate()
	}
	return nil
}

// WithJSONFile is a logger option that adds a JSON file handler with rotation.
func WithJSONFile(path string, opts ...FileHandlerOption) Option {
	return func(l *Logger) {
		// Combine rotation options with JSON formatter
		formatter := NewJSONFormatter()

		handler, err := NewFileHandler(path, formatter, opts...)
		if err != nil {
			// In a real scenario, you might want to handle this differently
			// For now, we'll just skip adding the handler
			return
		}

		if l.router == nil {
			l.router = NewRouter(handler)
		} else {
			l.router.AddHandler(handler)
		}
	}
}

// OpenFile opens or creates a file for writing logs with appropriate permissions.
// This is a helper function for cases where you need more control over file creation.
func OpenFile(path string, perm os.FileMode) (*os.File, error) {
	return os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, perm)
}
