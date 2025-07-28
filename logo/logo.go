// logo/logger.go
package logger

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
)

var (
	mu            sync.RWMutex
	consoleOn     = true
	outputs       []io.Writer
	includeSource = false          // Include source file and line number in logs
	includeTrace  = false          // Include stack trace in logs for log levels other than TRACE (FATAL)
	logLevel      = slog.LevelInfo // Default log level
	useJSONFormat = false
	jsonPretty    = false // Whether to use pretty JSON formatting
	attrOrder     = []string{"time", "level", "msg", "source"}
	colorEnabled  = true // Whether to enable colored output in console logs
)

var logger *Logger

const (
	LevelTrace slog.Level = slog.LevelDebug - 5
	LevelFatal slog.Level = slog.LevelError + 1
)

type LoggerOption func()

// Logger is a wrapper around slog.Logger that provides additional functionality.
// It allows for easy configuration of log outputs, levels, and custom handlers.
type Logger struct {
	*slog.Logger
}

// Init initializes the logger with the given options.
func Init(opts ...LoggerOption) {
	mu.Lock()
	defer mu.Unlock()

	// Reset the logger state
	outputs = nil

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

	// handler := slog.NewTextHandler(io.MultiWriter(outputs...), &slog.HandlerOptions{
	// 	Level:     logLevel,
	// 	AddSource: includeSource,
	// 	ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
	// 		// If this is a level attribute, replace it with our custom string
	// 		if a.Key == slog.LevelKey {
	// 			level, ok := a.Value.Any().(slog.Level)
	// 			if ok {
	// 				if levelStr := levelToString(level); levelStr != "" {
	// 					return slog.String("level", levelStr)
	// 				}
	// 			}
	// 		}
	// 		if includeSource && a.Key == slog.SourceKey {
	// 			source, ok := a.Value.Any().(*slog.Source)
	// 			if ok {
	// 				// Format as "file:line" without the full path
	// 				shortFile := source.File
	// 				if lastSlash := strings.LastIndex(shortFile, "/"); lastSlash >= 0 {
	// 					shortFile = shortFile[lastSlash+1:]
	// 				}
	// 				return slog.String("source", fmt.Sprintf("%s:%d", shortFile, source.Line))
	// 			}
	// 		}

	// 		return a
	// 	},
	// })
	// wrapped := NewLevelNameHandler(handler)
	// logger = &Logger{slog.New(wrapped)}
	handlerOptions := &slog.HandlerOptions{
		Level:     logLevel,
		AddSource: includeSource, // This controls whether source info is captured
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {

			// Handle source attribute - make sure we don't rename it
			// The standard key is "source" which matches what we want
			// Just make the display format nicer
			// Removed because while this is nice, if looking for the source file, the full verbose path is better
			// if includeSource && a.Key == slog.SourceKey {
			// 	source, ok := a.Value.Any().(*slog.Source)
			// 	if ok {
			// 		// Format as "file:line" without the full path
			// 		shortFile := source.File
			// 		if lastSlash := strings.LastIndex(shortFile, "/"); lastSlash >= 0 {
			// 			shortFile = shortFile[lastSlash+1:]
			// 		}
			// 		// Keep the original key
			// 		return slog.String(slog.SourceKey, fmt.Sprintf("%s:%d", shortFile, source.Line))
			// 	}
			// }

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

	// Use our level name handler
	// levelHandler := NewLevelNameHandler(handler)

	// Create the logger
	logger = &Logger{slog.New(handler)}
}

// SetLevel sets the log level for the logger.
func SetLevel(level slog.Level) LoggerOption {
	return func() {
		logLevel = level
	}
}

// DisableColors disables colored output in console logs.
// This is useful for environments where ANSI color codes might cause issues.
func DisableColors() LoggerOption {
	return func() {
		colorEnabled = false
	}
}

// EnableTrace enables trace logging, which is a level below DEBUG.
// This is useful for capturing detailed information during development or debugging.
func EnableTrace() LoggerOption {
	return func() {
		includeTrace = true
		// The issue might be here - we need to ensure this also lowers the log level
		// if trace is enabled, the log level needs to be set low enough to show trace
		if logLevel > LevelTrace {
			logLevel = LevelTrace
		}
	}
}

// UseConsole enables or disables the JSON console writer.
func UseJSON(pretty bool) LoggerOption {
	return func() {
		useJSONFormat = true
		jsonPretty = pretty
	}
}

// AddSource enables adding source file and line information to log messages
func AddSource() LoggerOption {
	return func() {
		includeSource = true
	}
}

// UseCustomHandler replaces the handler with a custom slog.Handler.
func UseCustomHandler(h slog.Handler) LoggerOption {
	return func() {
		logger = &Logger{slog.New(h)}
	}
}

// DisableConsole disables the console output.
// This is useful for applications that do not require console logging, such as background services or daemons.
func DisableConsole() LoggerOption {
	return func() {
		consoleOn = false
	}
}

// AddFileOutput adds a file output to the logger.
// It uses the Lumberjack package to manage log file rotation.
// This allows for log files to be rotated based on size, number of backups, and age
func AddFileOutput(path string, maxSize, backups, age int, compress bool) LoggerOption {
	return func() {
		lj := NewLumberjackWriter(path, maxSize, backups, age, compress)
		outputs = append(outputs, lj)
	}
}

// AddChannelOutput adds a channel output to the logger.
// This allows log messages to be sent to a channel for further processing or handling.
func AddChannelOutput(ch chan string) LoggerOption {
	return func() {
		outputs = append(outputs, NewChannelWriter(ch))
	}
}

// L returns the current logger instance.
// It is safe to call concurrently and returns the same logger instance.
// This is the main entry point for logging in the application.
func L() *Logger {
	mu.RLock()
	defer mu.RUnlock()
	return logger
}

// WithContext returns a logger with the request ID from the context.
// This should be customized based on your context handling for each application.
func WithContext(ctx context.Context) *Logger {
	mu.RLock()
	defer mu.RUnlock()
	if id, ok := ctx.Value("request_id").(string); ok {
		return &Logger{logger.With("request_id", id)}
	}
	return logger
}

// Trace logs with a level below DEBUG and includes a stack trace.
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

	if includeTrace {
		custom = append(custom, slog.String("trace", string(debug.Stack())))
	}

	rec := slog.NewRecord(timeNow(), LevelFatal, msg, pc)
	rec.AddAttrs(append(custom, filtered...)...)

	_ = l.Handler().Handle(context.Background(), rec)
	os.Exit(1)
}

// normalizeAttrs normalizes the attributes passed to the logger.
// It processes the attributes to ensure they are in the correct format for logging.
func normalizeAttrs(args ...any) []slog.Attr {
	var attrs []slog.Attr
	i := 0
	for i < len(args) {
		fmt.Println("Processing arg:", args[i])
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
			// Skip unknown
			i++
		}
	}
	return attrs
}

// timeNow returns the current time
func timeNow() time.Time {
	return time.Now()
}
