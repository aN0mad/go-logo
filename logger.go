package logo

import (
	"runtime"
	"strings"
	"sync"
)

// Logger is the main entry point for logging.
// It is immutable and safe for concurrent use.
type Logger struct {
	name       string
	router     *Router
	fields     map[string]interface{}
	addSource  bool
	sourceSkip int
	mu         sync.RWMutex
}

// New creates a new logger with the given options.
func New(opts ...Option) *Logger {
	logger := &Logger{
		name:       "default",
		fields:     make(map[string]interface{}),
		sourceSkip: 3, // Default skip for caller detection
	}

	// Apply options
	for _, opt := range opts {
		opt(logger)
	}

	// If no router was set, create a default one with no handlers
	if logger.router == nil {
		logger.router = NewRouter()
	}

	return logger
}

// Option is a functional option for configuring a logger.
type Option func(*Logger)

// WithRouter sets the router for the logger.
func WithRouter(router *Router) Option {
	return func(l *Logger) {
		l.router = router
	}
}

// WithName sets the logger name.
func WithName(name string) Option {
	return func(l *Logger) {
		l.name = name
	}
}

// WithFields sets initial fields for the logger.
func WithFields(fields map[string]interface{}) Option {
	return func(l *Logger) {
		l.fields = fields
	}
}

// WithSource enables source location tracking.
func WithSource(enabled bool) Option {
	return func(l *Logger) {
		l.addSource = enabled
	}
}

// Trace logs a trace-level message.
func (l *Logger) Trace(msg string, fields ...Field) {
	l.log(LevelTrace, msg, fields...)
}

// Debug logs a debug-level message.
func (l *Logger) Debug(msg string, fields ...Field) {
	l.log(LevelDebug, msg, fields...)
}

// Info logs an info-level message.
func (l *Logger) Info(msg string, fields ...Field) {
	l.log(LevelInfo, msg, fields...)
}

// Warn logs a warning-level message.
func (l *Logger) Warn(msg string, fields ...Field) {
	l.log(LevelWarn, msg, fields...)
}

// Error logs an error-level message.
func (l *Logger) Error(msg string, fields ...Field) {
	l.log(LevelError, msg, fields...)
}

// Fatal logs a fatal-level message.
func (l *Logger) Fatal(msg string, fields ...Field) {
	l.log(LevelFatal, msg, fields...)
}

// With creates a new logger with additional fields.
// The original logger is not modified.
func (l *Logger) With(fields ...Field) *Logger {
	l.mu.RLock()
	defer l.mu.RUnlock()

	// Create a new field map with existing fields
	newFields := make(map[string]interface{}, len(l.fields)+len(fields))
	for k, v := range l.fields {
		newFields[k] = v
	}

	// Add new fields
	for _, f := range fields {
		newFields[f.Key] = f.Value
	}

	return &Logger{
		name:       l.name,
		router:     l.router,
		fields:     newFields,
		addSource:  l.addSource,
		sourceSkip: l.sourceSkip,
	}
}

// WithName creates a new logger with a different name.
func (l *Logger) WithName(name string) *Logger {
	l.mu.RLock()
	defer l.mu.RUnlock()

	// Copy fields
	newFields := make(map[string]interface{}, len(l.fields))
	for k, v := range l.fields {
		newFields[k] = v
	}

	return &Logger{
		name:       name,
		router:     l.router,
		fields:     newFields,
		addSource:  l.addSource,
		sourceSkip: l.sourceSkip,
	}
}

// log is the internal logging method.
func (l *Logger) log(level Level, msg string, fields ...Field) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	// Merge logger fields with provided fields
	entryFields := make(map[string]interface{}, len(l.fields)+len(fields))
	for k, v := range l.fields {
		entryFields[k] = v
	}
	for _, f := range fields {
		entryFields[f.Key] = f.Value
	}

	// Create entry
	entry := NewEntry(level, msg, l.name, entryFields)

	// Add caller information if enabled
	if l.addSource {
		if caller := l.getCaller(); caller != nil {
			entry.WithCaller(caller)
		}
	}

	// Route the entry
	l.router.Route(entry)
}

// getCaller retrieves caller information.
func (l *Logger) getCaller() *Caller {
	pc, file, line, ok := runtime.Caller(l.sourceSkip)
	if !ok {
		return nil
	}

	// Get function name
	fn := runtime.FuncForPC(pc)
	funcName := ""
	if fn != nil {
		funcName = fn.Name()
		// Simplify function name (remove package path)
		if lastSlash := strings.LastIndex(funcName, "/"); lastSlash >= 0 {
			funcName = funcName[lastSlash+1:]
		}
	}

	// Simplify file path (just filename)
	if lastSlash := strings.LastIndex(file, "/"); lastSlash >= 0 {
		file = file[lastSlash+1:]
	}

	return &Caller{
		File:     file,
		Line:     line,
		Function: funcName,
	}
}

// Close closes the logger and its router.
func (l *Logger) Close() error {
	return l.router.Close()
}

// Field represents a key-value pair for structured logging.
type Field struct {
	Key   string
	Value interface{}
}

// F is a shorthand for creating a Field.
func F(key string, value interface{}) Field {
	return Field{Key: key, Value: value}
}
