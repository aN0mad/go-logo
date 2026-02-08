package logo

import (
	"io"
	"os"
)

// ConsoleHandler is a specialized handler for console output with pretty formatting.
type ConsoleHandler struct {
	*WriterHandler
}

// ConsoleHandlerOption configures a ConsoleHandler
type ConsoleHandlerOption func(*ConsoleHandler)

// WithConsoleLevel sets the minimum level for the console handler
func WithConsoleLevel(level Level) ConsoleHandlerOption {
	return func(h *ConsoleHandler) {
		h.BaseHandler.minLevel = level
	}
}

// NewConsoleHandler creates a new console handler with pretty formatting and color support.
func NewConsoleHandler(opts ...ConsoleHandlerOption) *ConsoleHandler {
	// Default to stdout
	writer := os.Stdout

	// Auto-detect color support
	colorSupport := NewColorSupport(writer, ColorAuto)

	// Create pretty formatter
	formatter := NewPrettyFormatter(colorSupport)

	// Create the writer handler
	writerHandler := NewWriterHandler(writer, formatter, LevelDebug)

	handler := &ConsoleHandler{
		WriterHandler: writerHandler,
	}

	// Apply options
	for _, opt := range opts {
		opt(handler)
	}

	return handler
}

// NewConsoleHandlerWithWriter creates a console handler with a custom writer.
// This is useful for testing or directing output to custom destinations.
func NewConsoleHandlerWithWriter(writer io.Writer, colorMode ColorMode, level Level) *ConsoleHandler {
	colorSupport := NewColorSupport(writer, colorMode)
	formatter := NewPrettyFormatter(colorSupport)
	writerHandler := NewWriterHandler(writer, formatter, level)

	return &ConsoleHandler{
		WriterHandler: writerHandler,
	}
}

// WithConsole is a logger option that adds a console handler
func WithConsole(opts ...ConsoleHandlerOption) Option {
	return func(l *Logger) {
		handler := NewConsoleHandler(opts...)

		if l.router == nil {
			l.router = NewRouter(handler)
		} else {
			l.router.AddHandler(handler)
		}
	}
}
