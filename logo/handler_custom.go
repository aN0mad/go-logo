// Package logger provides functionality for structured logging.
//
// This file contains the custom text handler implementation which formats
// log messages in a human-readable text format with ordered attributes.
package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"runtime"
	"slices"
	"strings"
)

// CustomTextHandler is a custom handler that produces text output with ordered attributes.
// It implements the slog.Handler interface and formats log messages in a consistent,
// readable format with configurable attribute ordering.
type CustomTextHandler struct {
	out       io.Writer
	opts      *slog.HandlerOptions
	attrOrder []string
}

// NewCustomTextHandler creates a new text handler with ordered attributes.
//
// Parameters:
//   - out: The io.Writer where log entries will be written
//   - opts: Handler options including log level and attribute replacements
//
// Returns a slog.Handler implementation
func NewCustomTextHandler(out io.Writer, opts *slog.HandlerOptions) slog.Handler {

	return &CustomTextHandler{
		out:       out,
		opts:      opts,
		attrOrder: attrOrder, // Use the global attrOrder defined in this package
	}
}

// Enabled implements Handler.Enabled.
// It checks if the given log level should be processed based on the configured minimum level.
//
// Parameters:
//   - ctx: The context for the logging operation
//   - level: The log level to check
//
// Returns true if the log level should be processed, false otherwise
func (h *CustomTextHandler) Enabled(ctx context.Context, level slog.Level) bool {
	minLevel := h.opts.Level.Level()
	return level >= minLevel
}

// Handle implements Handler.Handle.
// It processes a log record and outputs it in text format.
//
// Parameters:
//   - ctx: The context for the logging operation
//   - r: The log record to process
//
// Returns any error encountered during formatting or writing
func (h *CustomTextHandler) Handle(ctx context.Context, r slog.Record) error {
	// Collect all attributes in a map for reordering
	attrs := make(map[string]string)

	// Add standard attributes
	attrs["time"] = r.Time.Format("2006-01-02T15:04:05.000Z07:00")
	attrs["level"] = levelToString(r.Level)
	attrs["msg"] = r.Message

	// Add source if enabled
	if h.opts.AddSource {
		if source := r.PC; source != 0 {
			fs := runtime.CallersFrames([]uintptr{source})
			frame, _ := fs.Next()
			if frame.File != "" {
				shortFile := frame.File
				if lastSlash := strings.LastIndex(shortFile, "/"); lastSlash >= 0 {
					shortFile = shortFile[lastSlash+1:]
				}
				attrs["source"] = fmt.Sprintf("%s:%d", shortFile, frame.Line)
			}
		}
	}

	// Collect other attributes
	r.Attrs(func(a slog.Attr) bool {
		if a.Key != "level" && a.Key != "msg" && a.Key != "time" && a.Key != "source" {
			// Apply ReplaceAttr if provided
			if h.opts.ReplaceAttr != nil {
				a = h.opts.ReplaceAttr(nil, a)
			}

			// Only include non-empty attributes
			if !a.Equal(slog.Attr{}) {
				attrs[a.Key] = a.Value.String()
			}
		}
		return true
	})

	// Build the output string with ordered attributes
	var sb strings.Builder

	// First add the ordered attributes
	for _, key := range h.attrOrder {
		if val, ok := attrs[key]; ok {
			// Skip empty values
			if val == "" {
				continue
			}

			if sb.Len() > 0 {
				sb.WriteString(" ")
			}
			sb.WriteString(key)
			sb.WriteString("=")
			sb.WriteString(val)

			// Remove from map to avoid duplicates
			delete(attrs, key)
		}
	}

	// Then add remaining attributes in alphabetical order
	var remainingKeys []string
	for k := range attrs {
		remainingKeys = append(remainingKeys, k)
	}
	slices.Sort(remainingKeys)

	for _, key := range remainingKeys {
		if sb.Len() > 0 {
			sb.WriteString(" ")
		}
		sb.WriteString(key)
		sb.WriteString("=")
		sb.WriteString(attrs[key])
	}

	sb.WriteString("\n")

	// Write to output
	_, err := h.out.Write([]byte(sb.String()))
	return err
}

// WithAttrs implements Handler.WithAttrs.
// It returns a new handler with the given attributes.
//
// Parameters:
//   - attrs: The attributes to add to the handler
//
// Returns a new handler instance with the attributes
func (h *CustomTextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return NewCustomTextHandler(h.out, h.opts)
}

// WithGroup implements Handler.WithGroup.
// It returns a handler that adds the given group name to the attribute key path.
//
// Parameters:
//   - name: The group name
//
// Returns a handler that adds the group name to the attribute key path
func (h *CustomTextHandler) WithGroup(name string) slog.Handler {
	// Groups not implemented in this simple handler
	return h
}

// levelToString converts a slog.Level to its string representation.
//
// Parameters:
//   - l: The slog.Level to convert
//
// Returns the string representation of the log level, or empty string if not recognized
func levelToString(l slog.Level) string {
	switch l {
	case LevelTrace:
		return "TRACE"
	case slog.LevelDebug:
		return "DEBUG"
	case slog.LevelInfo:
		return "INFO"
	case slog.LevelWarn:
		return "WARN"
	case slog.LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return ""
	}
}
