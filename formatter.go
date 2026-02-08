package logo

import (
	"fmt"
	"io"
	"sync"
)

// Formatter defines how log entries are formatted for output.
type Formatter interface {
	Format(entry *Entry) ([]byte, error)
}

// TextFormatter formats log entries as human-readable text.
type TextFormatter struct{}

// NewTextFormatter creates a new text formatter.
func NewTextFormatter() *TextFormatter {
	return &TextFormatter{}
}

// Format formats a log entry as text.
func (f *TextFormatter) Format(entry *Entry) ([]byte, error) {
	// Format: timestamp [LEVEL] message key=value key=value
	timestamp := entry.Timestamp.Format("2006-01-02T15:04:05.000Z07:00")

	msg := fmt.Sprintf("%s [%s] %s", timestamp, entry.Level, entry.Message)

	// Add fields
	for key, value := range entry.Fields {
		msg += fmt.Sprintf(" %s=%v", key, value)
	}

	// Add caller if present
	if entry.Caller != nil {
		msg += fmt.Sprintf(" source=%s:%d", entry.Caller.File, entry.Caller.Line)
	}

	msg += "\n"
	return []byte(msg), nil
}

// WriterHandler is a handler that writes formatted entries to an io.Writer.
type WriterHandler struct {
	*BaseHandler
	writer    io.Writer
	formatter Formatter
	mu        sync.Mutex
}

// NewWriterHandler creates a new handler that writes to the given writer.
func NewWriterHandler(writer io.Writer, formatter Formatter, level Level) *WriterHandler {
	return &WriterHandler{
		BaseHandler: NewBaseHandler(level),
		writer:      writer,
		formatter:   formatter,
	}
}

// Handle processes and writes a log entry.
func (h *WriterHandler) Handle(entry *Entry) bool {
	if !h.Enabled(entry.Level) {
		return false
	}

	formatted, err := h.formatter.Format(entry)
	if err != nil {
		// If formatting fails, we can't do much - maybe write to stderr?
		return false
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	_, err = h.writer.Write(formatted)
	return err == nil
}

// Close closes the handler. If the underlying writer implements io.Closer,
// it will be closed as well.
func (h *WriterHandler) Close() error {
	if closer, ok := h.writer.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}
