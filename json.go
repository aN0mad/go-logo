package logo

import (
	"encoding/json"
	"fmt"
	"time"
)

// JSONFormatter formats log entries as JSON according to the schema.
type JSONFormatter struct {
	pretty        bool
	includeSource bool
}

// JSONFormatterOption configures a JSONFormatter
type JSONFormatterOption func(*JSONFormatter)

// WithPrettyJSON enables pretty-printed JSON output
func WithPrettyJSON(pretty bool) JSONFormatterOption {
	return func(f *JSONFormatter) {
		f.pretty = pretty
	}
}

// WithJSONSource enables source location in JSON output
func WithJSONSource(enabled bool) JSONFormatterOption {
	return func(f *JSONFormatter) {
		f.includeSource = enabled
	}
}

// NewJSONFormatter creates a new JSON formatter.
func NewJSONFormatter(opts ...JSONFormatterOption) *JSONFormatter {
	f := &JSONFormatter{
		pretty:        false,
		includeSource: false,
	}

	for _, opt := range opts {
		opt(f)
	}

	return f
}

// Format formats a log entry as JSON according to the schema.
// Schema fields (required):
//   - ts: RFC3339Nano UTC timestamp
//   - level: log level string
//   - msg: log message
//   - logger: logger name
//
// Optional fields:
//   - source: source location (if enabled and available)
//
// User fields are added as top-level keys after filtering reserved keys.
func (f *JSONFormatter) Format(entry *Entry) ([]byte, error) {
	// Build the JSON object according to schema
	obj := make(map[string]interface{})

	// Required fields in canonical order
	obj[FieldTimestamp] = entry.Timestamp.UTC().Format(time.RFC3339Nano)
	obj[FieldLevel] = entry.Level.String()
	obj[FieldMessage] = entry.Message
	obj[FieldLogger] = entry.Logger

	// Optional source field
	if f.includeSource && entry.Caller != nil {
		obj[FieldSource] = fmt.Sprintf("%s:%d", entry.Caller.File, entry.Caller.Line)
	}

	// Add user fields (automatically filtered for reserved keys)
	// User fields are added after schema fields to ensure schema fields take precedence
	for key, value := range entry.Fields {
		if !IsReservedKey(key) {
			obj[key] = value
		}
	}

	// Marshal to JSON
	var data []byte
	var err error

	if f.pretty {
		data, err = json.MarshalIndent(obj, "", "  ")
	} else {
		data, err = json.Marshal(obj)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// Add newline for consistent output
	data = append(data, '\n')

	return data, nil
}

// JSONHandler is a specialized handler for JSON output.
type JSONHandler struct {
	*WriterHandler
}

// NewJSONHandler creates a new JSON handler with the specified writer and options.
func NewJSONHandler(writer interface{}, level Level, opts ...JSONFormatterOption) *JSONHandler {
	formatter := NewJSONFormatter(opts...)

	// Handle both io.Writer and *os.File
	var w interface {
		Write([]byte) (int, error)
	}

	switch v := writer.(type) {
	case interface{ Write([]byte) (int, error) }:
		w = v
	default:
		panic("writer must implement Write method")
	}

	writerHandler := NewWriterHandler(w, formatter, level)

	return &JSONHandler{
		WriterHandler: writerHandler,
	}
}
