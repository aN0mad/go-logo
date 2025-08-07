// Package logo provides a flexible and extensible structured logging framework
// built on top of Go's standard library slog package. It enhances slog with features
// like customizable outputs, log levels including TRACE and FATAL, colorized console
// output, file rotation, log channels, and support for JSON formatting.
//
// This package supports multiple logging backends simultaneously, including console,
// files (with rotation via lumberjack), and channels for custom processing.
//
// Basic usage:
//
//	logo.Init(
//		logo.SetLevel(slog.LevelInfo),
//		logo.EnableTrace(),
//		logo.AddSource(),
//	)
//
//	log := logo.L()
//	log.Info("Hello, world!", "user", "gopher")
//
// For file logging with rotation:
//
//	logo.Init(
//		logo.AddFileOutput("/path/to/log.file", 10, 3, 30, true),
//	)
//
// For JSON output:
//
//	logo.Init(
//		logo.UseJSON(true), // true for pretty-printed JSON
//	)
package logo

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"sync"
	"time"

	"github.com/aN0mad/lumberjack/v2"
)

var (
	// mu protects shared state in the logger package
	mu sync.RWMutex

	// consoleOn determines whether console output is enabled
	CONSOLEON = true

	// outputs is the list of writers to send log output to
	OUTPUTS []io.Writer

	// includeSource determines whether to include source file and line information in logs
	INCLUDESOURCE = false

	// includeStackTraces determines whether to include stack traces in logs
	INCLUDESTACKTRACES = false

	// logLevel is the minimum log level that will be output
	LOGLEVEL = slog.LevelInfo

	// useJSONFormat determines whether to output logs in JSON format
	USEJSONFORMAT = false

	// jsonPretty determines whether JSON output should be pretty-printed
	JSONPRETTY = false

	// attrOrder defines the default attribute order for structured log entries
	attrOrder = []string{"time", "level", "msg", "source"}

	// colorEnabled determines whether to use colored output in console logs
	COLORENABLED = true

	// fileWriters tracks all lumberjack loggers for proper closing
	FILEWRITERS []*lumberjack.Logger

	// Predefined log levels for convenience
	LevelInfo  = slog.LevelInfo
	LevelDebug = slog.LevelDebug
	LevelError = slog.LevelError
)

// osExit is a variable that points to os.Exit to allow for testing
// of the Fatal function without actually terminating the program.
var osExit = os.Exit

// Global logger instance
var logger *Logger

// Constants for additional log levels not provided by the standard slog package.
const (
	// LevelTrace is a level below DEBUG for extremely verbose diagnostic information
	LevelTrace slog.Level = slog.LevelDebug - 5

	// LevelFatal is a level above ERROR that indicates a fatal error condition
	// which will cause the application to terminate after logging
	LevelFatal slog.Level = slog.LevelError + 1
)

// loggerContext holds all the configuration for a specific logger instance
type loggerContext struct {
	outputs            []io.Writer
	consoleOn          bool
	useJSONFormat      bool
	jsonPretty         bool
	includeSource      bool
	includeStackTraces bool
	logLevel           slog.Level
	colorEnabled       bool
	fileWriters        []*lumberjack.Logger
	customHandler      slog.Handler
}

// LoggerOption is a functional option type for configuring the logger.
// This allows for a flexible and extensible way to configure the logger
// with various options.
type LoggerOption func(*loggerContext)

// Logger is the main logging structure that wraps slog.Logger.
// It provides structured logging capabilities with additional
// convenience methods for different log levels.
type Logger struct {
	*slog.Logger
	ctx *loggerContext // Contains all configuration including file writers
}

// Init initializes the global default logger with the given options.
// This configures a single global logger instance used by L().
// To create independent loggers with their own configurations, use NewLogger() instead.
func Init(opts ...LoggerOption) {
	mu.Lock()
	defer mu.Unlock()

	// Create a new logger with the provided options
	logger = NewLogger(opts...)
}

// NewLogger creates a new independent logger instance with its own configuration.
// Unlike Init() which configures a global singleton logger, NewLogger returns a
// completely separate logger that can be configured differently from other loggers.
func NewLogger(opts ...LoggerOption) *Logger {
	// Create a configuration context for this specific logger
	ctx := &loggerContext{
		outputs:            OUTPUTS,
		consoleOn:          CONSOLEON,
		useJSONFormat:      USEJSONFORMAT,
		jsonPretty:         JSONPRETTY,
		includeSource:      INCLUDESOURCE,
		includeStackTraces: INCLUDESTACKTRACES,
		logLevel:           LOGLEVEL,
		colorEnabled:       COLORENABLED,
		fileWriters:        nil,
	}

	// Apply all options to this context
	for _, opt := range opts {
		opt(ctx)
	}

	// If a custom handler was specified, use it directly
	if ctx.customHandler != nil {
		return &Logger{
			Logger: slog.New(ctx.customHandler),
			ctx:    ctx,
		}
	}

	// If no outputs are specified, default to console output unless disabled manually
	if ctx.consoleOn && len(ctx.outputs) == 0 {
		if !ctx.useJSONFormat {
			// Only use styled writer for text format
			ctx.outputs = append(ctx.outputs, NewStyledConsoleWriter(os.Stdout, ctx))
		} else {
			// For JSON, use stdout directly
			ctx.outputs = append(ctx.outputs, os.Stdout)
		}
	}

	// Configure handler options
	handlerOptions := &slog.HandlerOptions{
		Level:     ctx.logLevel,
		AddSource: ctx.includeSource,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Handle attribute replacements if needed
			return a
		},
	}

	// Create the handler
	var handler slog.Handler
	if len(ctx.outputs) > 0 {
		multiwriter := io.MultiWriter(ctx.outputs...)

		// Choose handler based on format
		if ctx.useJSONFormat {
			handler = NewJSONHandler(multiwriter, handlerOptions, ctx.jsonPretty)
		} else {
			handler = NewCustomTextHandler(multiwriter, handlerOptions)
		}
	} else {
		// Fallback to a no-op handler if no outputs
		handler = slog.NewTextHandler(io.Discard, handlerOptions)
	}

	// Create and return the logger
	return &Logger{
		Logger: slog.New(handler),
		ctx:    ctx, // Store the context with file writers
	}
}

// SetLevel sets the minimum log level that will be logged.
// Any log messages with a level lower than this will be ignored.
//
// Parameters:
//   - level: The minimum log level to log (e.g., slog.LevelDebug, slog.LevelInfo)
//
// Returns:
//   - LoggerOption: A function that can be passed to Init() to configure the logger
func SetLevel(level slog.Level) LoggerOption {
	return func(ctx *loggerContext) {
		ctx.logLevel = level
	}
}

// DisableColors disables colored output in console logs.
// This is useful for environments where ANSI color codes might cause issues,
// such as when logging to files or in environments that don't support colors.
//
// Returns:
//   - LoggerOption: A function that can be passed to Init() to disable colored output
func DisableColors() LoggerOption {
	return func(ctx *loggerContext) {
		ctx.colorEnabled = false
	}
}

// EnableLogLevelTrace enables trace logging level (which is a level below DEBUG).
// This is useful for capturing detailed information during development or debugging.
//
// Returns:
//   - LoggerOption: A function that can be passed to Init() to set the trace log level
func EnableLogLevelTrace() LoggerOption {
	return func(ctx *loggerContext) {
		ctx.logLevel = LevelTrace
	}
}

// EnableStackTraces enables stack trace inclusion in log messages.
// When enabled, a stack trace will be included with log messages,
// which can be helpful for debugging and error tracking.
//
// Returns:
//   - LoggerOption: A function that can be passed to Init() to enable stack traces
func EnableStackTraces() LoggerOption {
	return func(ctx *loggerContext) {
		ctx.includeStackTraces = true
	}
}

// UseJSON configures the logger to output logs in JSON format.
//
// Parameters:
//   - pretty: If true, JSON will be formatted with indentation for better readability.
//     If false, JSON will be compact without extra whitespace.
//
// Returns:
//   - LoggerOption: A function that can be passed to Init() to use JSON formatting
func UseJSON(pretty bool) LoggerOption {
	return func(ctx *loggerContext) {
		ctx.useJSONFormat = true
		ctx.jsonPretty = pretty
	}
}

// Add this function to the logo package

// AddConsoleOutput adds console output to the logger.
// This is useful when you want to explicitly enable console output
// even when other outputs are configured.
//
// Returns:
//   - LoggerOption: A function that can be passed to Init() to enable console output
func AddConsoleOutput() LoggerOption {
	return func(ctx *loggerContext) {
		hasConsoleWriter := false
		for _, output := range ctx.outputs {
			if _, ok := output.(*StyledConsoleWriter); ok {
				hasConsoleWriter = true
				break
			}
			if output == os.Stdout {
				hasConsoleWriter = true
				break
			}
		}

		if !hasConsoleWriter {
			if ctx.useJSONFormat {
				ctx.outputs = append(ctx.outputs, os.Stdout)
			} else {
				ctx.outputs = append(ctx.outputs, NewStyledConsoleWriter(os.Stdout, ctx))
			}
		}
		ctx.consoleOn = true
	}
}

// AddSource enables adding source file and line information to log messages.
// This helps with debugging by showing where each log message originated from.
//
// Returns:
//   - LoggerOption: A function that can be passed to Init() to include source information
func AddSource() LoggerOption {
	return func(ctx *loggerContext) {
		ctx.includeSource = true
	}
}

// UseCustomHandler replaces the handler with a custom slog.Handler.
// This allows for complete customization of the logging behavior.
//
// Parameters:
//   - h: A custom slog.Handler implementation
//
// Returns:
//   - LoggerOption: A function that can be passed to Init() or NewLogger() to use a custom handler
func UseCustomHandler(h slog.Handler) LoggerOption {
	return func(ctx *loggerContext) {
		// When used with NewLogger, just set the custom handler in the context
		if ctx != nil {
			ctx.customHandler = h
		} else {
			// For backward compatibility with direct calls
			mu.Lock()
			defer mu.Unlock()
			logger = &Logger{
				Logger: slog.New(h),
				ctx:    &loggerContext{},
			}
		}
	}
}

// DisableConsole disables the console output.
// This is useful for applications that do not require console logging,
// such as background services or daemons.
//
// Returns:
//   - LoggerOption: A function that can be passed to Init() to disable console output
func DisableConsole() LoggerOption {
	return func(ctx *loggerContext) {
		ctx.consoleOn = false
	}
}

// AddFileOutput adds file output to the logger with rotation support.
// This allows log messages to be written to a file, with automatic rotation
// when the file reaches the specified maximum size.
//
// Parameters:
//   - filepath: The path to the log file
//   - maxSizeMB: Maximum size of the log file in megabytes before rotation
//   - maxBackups: Maximum number of old log files to retain
//   - maxAgeDays: Maximum number of days to retain old log files
//   - compress: If true, rotated log files will be compressed using gzip
//
// Returns:
//   - LoggerOption: A function that can be passed to Init() to add file output
func AddFileOutput(filename string, maxSize, maxBackups, maxAge int, compress bool) LoggerOption {
	return func(ctx *loggerContext) {
		// Ensure directory exists
		dir := filepath.Dir(filename)
		if dir != "" && dir != "." {
			if err := os.MkdirAll(dir, 0755); err != nil {
				fmt.Fprintf(os.Stderr, "Error creating log directory: %v\n", err)
				return
			}
		}

		fileWriter := &lumberjack.Logger{
			Filename:   filename,
			MaxSize:    maxSize,
			MaxBackups: maxBackups,
			MaxAge:     maxAge,
			Compress:   compress,
		}

		// Test that the file can be created
		if f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644); err == nil {
			f.Close()
		} else {
			fmt.Fprintf(os.Stderr, "Error testing log file creation: %v\n", err)
			return
		}

		ctx.outputs = append(ctx.outputs, fileWriter)
		ctx.fileWriters = append(ctx.fileWriters, fileWriter)
	}
}

// Close properly closes all resources used by the logger.
// This ensures that all log messages are flushed and file handles are closed.
// It should be called before the application exits.
//
// Returns:
//   - error: Any error encountered while closing resources
func Close() error {
	mu.Lock()
	defer mu.Unlock()

	if logger == nil {
		return nil
	}

	// Use the logger.ctx to access file writers
	return logger.Close()
}

// Close properly closes all resources used by this logger instance
func (l *Logger) Close() error {
	if l == nil || l.ctx == nil {
		return nil
	}

	var lastErr error

	// Close all file writers
	for _, fw := range l.ctx.fileWriters {
		if fw != nil {
			if err := fw.Close(); err != nil {
				lastErr = err
			}
		}
	}

	// Also sync any other writers that might implement Sync()
	for _, out := range l.ctx.outputs {
		if syncer, ok := out.(interface{ Sync() error }); ok {
			if err := syncer.Sync(); err != nil && lastErr == nil {
				lastErr = err
			}
		}
	}

	return lastErr
}

// AddChannelOutput adds a channel output to the logger.
// This allows log messages to be sent to a channel for further processing or handling.
//
// Parameters:
//   - ch: A channel of strings that will receive log messages
//
// Returns:
//   - LoggerOption: A function that can be passed to Init() to add channel output
func AddChannelOutput(ch chan string) LoggerOption {
	return func(ctx *loggerContext) {
		ctx.outputs = append(ctx.outputs, NewChannelWriter(ch))
	}
}

// L returns the current logger instance.
// It is safe to call concurrently and returns the same logger instance.
// This is the main entry point for logging in the application.
//
// Returns:
//   - *Logger: The configured Logger instance
func L() *Logger {
	mu.RLock()
	defer mu.RUnlock()
	return logger
}

// Trace logs with a level below DEBUG and includes a stack trace.
// This is useful for very detailed diagnostic information typically
// needed during development or troubleshooting.
//
// Parameters:
//   - msg: The message to log
//   - attrs: Additional attributes to include with the log entry,
//     provided as alternating keys and values
//
// Returns:
//   - None
func (l *Logger) Trace(msg string, attrs ...any) {
	if !l.Enabled(context.Background(), LevelTrace) {
		return
	}
	pc, file, line, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc).Name()

	userAttrs := normalizeAttrs(attrs...)

	filtered := userAttrs[:0]
	for _, a := range userAttrs {
		if a.Key != slog.SourceKey {
			filtered = append(filtered, a)
		}
	}

	custom := []slog.Attr{
		slog.String("trace", string(debug.Stack())),
		slog.String("source", fmt.Sprintf("%s:%d (%s)", file, line, fn)),
	}

	rec := slog.NewRecord(timeNow(), LevelTrace, msg, pc)
	rec.AddAttrs(append(custom, filtered...)...)

	_ = l.Handler().Handle(context.Background(), rec)
}

// Fatal logs the message and exits the program with status 1.
// This should be used for critical errors that require immediate termination.
//
// Parameters:
//   - msg: The message to log
//   - attrs: Additional attributes to include with the log entry,
//     provided as alternating keys and values
//
// Returns:
//   - None: This function does not return as it calls os.Exit(1)
func (l *Logger) Fatal(msg string, attrs ...any) {
	if !l.Enabled(context.Background(), LevelFatal) {
		osExit(1) // Still exit even if logging is disabled
		return
	}

	pc, file, line, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc).Name()

	userAttrs := normalizeAttrs(attrs...)

	filtered := userAttrs[:0]
	for _, a := range userAttrs {
		if a.Key != slog.SourceKey {
			filtered = append(filtered, a)
		}
	}

	custom := []slog.Attr{
		slog.String("source", fmt.Sprintf("%s:%d (%s)", file, line, fn)),
	}

	// Check if this specific logger has stack traces enabled
	includeStackTracesForThisLogger := false
	if l.ctx != nil {
		includeStackTracesForThisLogger = l.ctx.includeStackTraces
	} else {
		// Fall back to global setting for backward compatibility
		mu.RLock()
		includeStackTracesForThisLogger = INCLUDESTACKTRACES
		mu.RUnlock()
	}

	if includeStackTracesForThisLogger {
		custom = append(custom, slog.String("trace", string(debug.Stack())))
	}

	rec := slog.NewRecord(timeNow(), LevelFatal, msg, pc)
	rec.AddAttrs(append(custom, filtered...)...)

	_ = l.Handler().Handle(context.Background(), rec)
	osExit(1)
}

// normalizeAttrs normalizes the attributes passed to the logger.
// It processes the attributes to ensure they are in the correct format for logging.
//
// Parameters:
//   - args: Variable arguments that should be pairs of string keys and arbitrary values
//
// Returns:
//   - []slog.Attr: A slice of slog.Attr representing the normalized attributes
func normalizeAttrs(args ...any) []slog.Attr {
	var attrs []slog.Attr
	i := 0
	for i < len(args) {
		switch v := args[i].(type) {
		case slog.Attr:
			attrs = append(attrs, v)
			i++
		case string:
			if i+1 < len(args) {
				attrs = append(attrs, slog.Any(v, args[i+1]))
				i += 2
			} else {
				// Malformed key without value
				attrs = append(attrs, slog.Any(v, "(MISSING)"))
				i++
			}
		default:
			// Skip non-string keys and their values if possible
			if i+1 < len(args) {
				// If there's a potential value after this non-string key,
				// skip both the key and the value
				i += 2
			} else {
				// Otherwise just skip this single item
				i++
			}
		}
	}
	return attrs
}

// timeNow returns the current time.
// This function exists to make testing easier by allowing time to be mocked.
//
// Returns:
//   - time.Time: The current time
func timeNow() time.Time {
	return time.Now()
}

// SetFileHandlerForTesting is a special helper for test files
// to ensure proper handling of file outputs during testing
//
// Parameters:
//   - w: The io.Writer to use for log output during testing
//
// Returns:
//   - LoggerOption: A function that can be passed to Init() to set a test file handler
func SetFileHandlerForTesting(w io.Writer) LoggerOption {
	return func(ctx *loggerContext) {
		// Add the provided writer as a file output
		ctx.outputs = append(ctx.outputs, w)
	}
}
