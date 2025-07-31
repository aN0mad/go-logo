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
	consoleOn = true

	// outputs is the list of writers to send log output to
	outputs []io.Writer

	// includeSource determines whether to include source file and line information in logs
	includeSource = false

	// includeStackTraces determines whether to include stack traces in logs
	includeStackTraces = false

	// logLevel is the minimum log level that will be output
	logLevel = slog.LevelInfo

	// useJSONFormat determines whether to output logs in JSON format
	useJSONFormat = false

	// jsonPretty determines whether JSON output should be pretty-printed
	jsonPretty = false

	// attrOrder defines the default attribute order for structured log entries
	attrOrder = []string{"time", "level", "msg", "source"}

	// colorEnabled determines whether to use colored output in console logs
	colorEnabled = true

	// fileWriters tracks all lumberjack loggers for proper closing
	fileWriters []*lumberjack.Logger
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

// LoggerOption is a functional option type for configuring the logger.
// This allows for a flexible and extensible way to configure the logger
// with various options.
type LoggerOption func()

// Logger is the main logging structure that wraps slog.Logger.
// It provides structured logging capabilities with additional
// convenience methods for different log levels.
type Logger struct {
	*slog.Logger
}

// Init initializes the logger with the given options.
// It sets up outputs, formats, and handlers based on the provided options.
// If no options are provided, it defaults to console output in text format.
//
// Options can be combined to customize the logger behavior:
//
//	logger.Init(
//		logger.SetLevel(slog.LevelDebug),
//		logger.EnableTrace(),
//		logger.AddSource(),
//		logger.AddFileOutput("/var/log/app.log", 10, 3, 30, true),
//	)
//
// Parameters:
//   - opts: A variadic list of LoggerOption functions to configure the logger
//
// Returns:
//   - None
func Init(opts ...LoggerOption) {
	mu.Lock()
	defer mu.Unlock()

	// Reset the logger state
	outputs = nil

	// Reset to default values before applying options
	logLevel = slog.LevelInfo // Default to INFO
	useJSONFormat = false
	jsonPretty = false
	includeSource = false
	includeStackTraces = false
	consoleOn = true
	colorEnabled = true

	// Process options
	for _, opt := range opts {
		opt()
	}

	// If no outputs are specified, default to console output unless disabled manually
	if consoleOn {
		if !useJSONFormat {
			// Only use styled writer for text format
			outputs = append(outputs, NewStyledConsoleWriter(os.Stdout))
		} else {
			// For JSON, use stdout directly
			outputs = append(outputs, os.Stdout)
		}
	}

	// If no outputs are specified, default to console output
	handlerOptions := &slog.HandlerOptions{
		Level:     logLevel,
		AddSource: includeSource,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {

			return a
		},
	}

	// Create the handler directly without wrapping
	multiwriter := io.MultiWriter(outputs...)
	var handler slog.Handler

	// Choose the appropriate handler based on format
	if useJSONFormat {
		// Use JSON format
		handler = NewJSONHandler(multiwriter, handlerOptions, jsonPretty)
	} else {
		// Use text format
		handler = NewCustomTextHandler(multiwriter, handlerOptions)
	}

	// Create the logger
	logger = &Logger{slog.New(handler)}
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
	return func() {
		logLevel = level
	}
}

// DisableColors disables colored output in console logs.
// This is useful for environments where ANSI color codes might cause issues,
// such as when logging to files or in environments that don't support colors.
//
// Returns:
//   - LoggerOption: A function that can be passed to Init() to disable colored output
func DisableColors() LoggerOption {
	return func() {
		colorEnabled = false
	}
}

// EnableLogLevelTrace enables trace logging level (which is a level below DEBUG).
// This is useful for capturing detailed information during development or debugging.
//
// Returns:
//   - LoggerOption: A function that can be passed to Init() to set the trace log level
func EnableLogLevelTrace() LoggerOption {
	return func() {
		logLevel = LevelTrace
	}
}

// EnableStackTraces enables stack trace inclusion in log messages.
// When enabled, a stack trace will be included with log messages,
// which can be helpful for debugging and error tracking.
//
// Returns:
//   - LoggerOption: A function that can be passed to Init() to enable stack traces
func EnableStackTraces() LoggerOption {
	return func() {
		includeStackTraces = true
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
	return func() {
		useJSONFormat = true
		jsonPretty = pretty
	}
}

// AddSource enables adding source file and line information to log messages.
// This helps with debugging by showing where each log message originated from.
//
// Returns:
//   - LoggerOption: A function that can be passed to Init() to include source information
func AddSource() LoggerOption {
	return func() {
		includeSource = true
	}
}

// UseCustomHandler replaces the handler with a custom slog.Handler.
// This allows for complete customization of the logging behavior.
//
// Parameters:
//   - h: A custom slog.Handler implementation
//
// Returns:
//   - LoggerOption: A function that can be passed to Init() to use a custom handler
func UseCustomHandler(h slog.Handler) LoggerOption {
	return func() {
		logger = &Logger{slog.New(h)}
	}
}

// DisableConsole disables the console output.
// This is useful for applications that do not require console logging,
// such as background services or daemons.
//
// Returns:
//   - LoggerOption: A function that can be passed to Init() to disable console output
func DisableConsole() LoggerOption {
	return func() {
		consoleOn = false
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
func AddFileOutput(filepath string, maxSizeMB, maxBackups, maxAgeDays int, compress bool) LoggerOption {
	return func() {
		w := NewLumberjackWriter(filepath, maxSizeMB, maxBackups, maxAgeDays, compress)
		fileWriters = append(fileWriters, w)
		outputs = append(outputs, w)
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

	var lastErr error

	// Close each file writer
	for _, fw := range fileWriters {
		if err := fw.Close(); err != nil {
			lastErr = err
		}
	}

	// Clear the writers list
	fileWriters = nil

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
	return func() {
		outputs = append(outputs, NewChannelWriter(ch))
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

	if includeStackTraces {
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
	return func() {
		// Add the provided writer as a file output
		outputs = append(outputs, w)
	}
}
