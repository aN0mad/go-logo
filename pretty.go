package logo

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"time"
)

// PrettyFormatter formats log entries in a human-readable, colorized format.
type PrettyFormatter struct {
	colorSupport *ColorSupport
	timeFormat   string
	showSource   bool
}

// PrettyFormatterOption configures a PrettyFormatter
type PrettyFormatterOption func(*PrettyFormatter)

// WithColorMode sets the color mode for the formatter
func WithColorMode(mode ColorMode) PrettyFormatterOption {
	return func(f *PrettyFormatter) {
		// Color support will be set when we know the writer
		// Store mode for later
	}
}

// WithTimeFormat sets a custom time format
func WithTimeFormat(format string) PrettyFormatterOption {
	return func(f *PrettyFormatter) {
		f.timeFormat = format
	}
}

// WithSourceInfo enables source location in output
func WithSourceInfo(enabled bool) PrettyFormatterOption {
	return func(f *PrettyFormatter) {
		f.showSource = enabled
	}
}

// NewPrettyFormatter creates a new pretty formatter with color support detection.
func NewPrettyFormatter(colorSupport *ColorSupport, opts ...PrettyFormatterOption) *PrettyFormatter {
	f := &PrettyFormatter{
		colorSupport: colorSupport,
		timeFormat:   "15:04:05.000",
		showSource:   false,
	}

	for _, opt := range opts {
		opt(f)
	}

	return f
}

// Format formats a log entry in a pretty, human-readable format.
// Example output: 15:04:05.123 INFO  service started version=1.0.0
func (f *PrettyFormatter) Format(entry *Entry) ([]byte, error) {
	var buf bytes.Buffer

	// Timestamp in local time with custom format
	timestamp := entry.Timestamp.Local().Format(f.timeFormat)
	buf.WriteString(f.colorSupport.Colorize(timestamp, colorGray))
	buf.WriteString(" ")

	// Level with color and padding
	levelStr := f.formatLevel(entry.Level)
	buf.WriteString(levelStr)
	buf.WriteString(" ")

	// Message
	buf.WriteString(entry.Message)

	// Fields (sorted for consistency)
	if len(entry.Fields) > 0 {
		buf.WriteString(" ")
		buf.WriteString(f.formatFields(entry.Fields))
	}

	// Source location if enabled and available
	if f.showSource && entry.Caller != nil {
		buf.WriteString(" ")
		source := fmt.Sprintf("<%s:%d>", entry.Caller.File, entry.Caller.Line)
		buf.WriteString(f.colorSupport.Colorize(source, colorGray))
	}

	buf.WriteString("\n")
	return buf.Bytes(), nil
}

// formatLevel formats the level with color and padding
func (f *PrettyFormatter) formatLevel(level Level) string {
	// Pad level to 5 characters for alignment
	levelStr := level.String()
	padding := 5 - len(levelStr)
	if padding > 0 {
		levelStr = levelStr + strings.Repeat(" ", padding)
	}

	color := LevelColorBold(level)
	return f.colorSupport.Colorize(levelStr, color)
}

// formatFields formats the fields as key=value pairs
func (f *PrettyFormatter) formatFields(fields map[string]interface{}) string {
	if len(fields) == 0 {
		return ""
	}

	// Sort keys for consistent output
	keys := make([]string, 0, len(fields))
	for k := range fields {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var parts []string
	for _, key := range keys {
		value := fields[key]

		// Format the value appropriately
		var valueStr string
		switch v := value.(type) {
		case string:
			// Quote strings if they contain spaces
			if strings.Contains(v, " ") {
				valueStr = fmt.Sprintf(`"%s"`, v)
			} else {
				valueStr = v
			}
		case time.Duration:
			valueStr = v.String()
		case error:
			valueStr = v.Error()
		default:
			valueStr = fmt.Sprintf("%v", v)
		}

		// Colorize the key
		keyColored := f.colorSupport.Colorize(key, colorCyan)
		parts = append(parts, fmt.Sprintf("%s=%s", keyColored, valueStr))
	}

	return strings.Join(parts, " ")
}
