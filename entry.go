package logo

import (
	"time"
)

// Entry represents a single log event with its metadata and fields.
type Entry struct {
	Timestamp time.Time
	Level     Level
	Message   string
	Logger    string
	Fields    map[string]interface{}
	Caller    *Caller
}

// Caller contains source location information.
type Caller struct {
	File     string
	Line     int
	Function string
}

// NewEntry creates a new log entry with the given parameters.
func NewEntry(level Level, message string, logger string, fields map[string]interface{}) *Entry {
	return &Entry{
		Timestamp: time.Now().UTC(),
		Level:     level,
		Message:   message,
		Logger:    logger,
		Fields:    fields,
	}
}

// WithCaller adds caller information to the entry.
func (e *Entry) WithCaller(caller *Caller) *Entry {
	e.Caller = caller
	return e
}

// Clone creates a deep copy of the entry.
func (e *Entry) Clone() *Entry {
	fields := make(map[string]interface{}, len(e.Fields))
	for k, v := range e.Fields {
		fields[k] = v
	}

	entry := &Entry{
		Timestamp: e.Timestamp,
		Level:     e.Level,
		Message:   e.Message,
		Logger:    e.Logger,
		Fields:    fields,
		Caller:    e.Caller,
	}

	return entry
}
