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
	attrs     []slog.Attr
	groups    []string
}

// NewCustomTextHandler creates a new text handler with ordered attributes.
//
// Parameters:
//   - out: The io.Writer where log entries will be written
//   - opts: Handler options including log level and attribute replacements
//
// Returns a slog.Handler implementation
func NewCustomTextHandler(out io.Writer, opts *slog.HandlerOptions) slog.Handler {

	// if opts == nil {
	// 	opts = &slog.HandlerOptions{}
	// }

	return &CustomTextHandler{
		out:       out,
		opts:      opts,
		attrOrder: attrOrder, // Use the global attrOrder defined in this package
		attrs:     []slog.Attr{},
		groups:    nil,
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
				if lastSlash := strings.LastIndex(shortFile, "/"); lastSlash >= 0 { // TODO: Fix to remove short source
					shortFile = shortFile[lastSlash+1:]
				}
				attrs["source"] = fmt.Sprintf("%s:%d", shortFile, frame.Line)
			}
		}
	}

	// Process handler attributes (added via With())
	for _, attr := range h.attrs {
		if attr.Key != "level" && attr.Key != "msg" && attr.Key != "time" && attr.Key != "source" { // TODO: Replace with not in attrOrder `if !slices.Contains(h.attrOrder, attr.Key) {`
			attrs[attr.Key] = attr.Value.String()
		}
	}

	// Process record attributes
	r.Attrs(func(a slog.Attr) bool {
		if a.Key != "level" && a.Key != "msg" && a.Key != "time" && a.Key != "source" { // TODO: Replace with not in attrOrder `if !slices.Contains(h.attrOrder, attr.Key) {`
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
// WithAttrs implements slog.Handler.WithAttrs.
func (h *CustomTextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	// Create a new handler with the same settings
	newHandler := &CustomTextHandler{
		out:       h.out,
		opts:      h.opts,
		attrOrder: h.attrOrder,
		attrs:     append([]slog.Attr{}, h.attrs...), // Copy existing attributes
		groups:    append([]string{}, h.groups...),   // Copy existing groups
	}

	// Process and store the new attributes
	for _, attr := range attrs {
		// Apply ReplaceAttr if set
		if h.opts != nil && h.opts.ReplaceAttr != nil {
			attr = h.opts.ReplaceAttr(h.groups, attr)
		}

		// Skip empty attributes
		if attr.Equal(slog.Attr{}) {
			continue
		}

		// Add the attribute
		newHandler.attrs = append(newHandler.attrs, attr)
	}

	return newHandler
}

// WithGroup implements Handler.WithGroup.
// It returns a handler that adds the given group name to the attribute key path.
//
// Parameters:
//   - name: The group name
//
// Returns a handler that adds the group name to the attribute key path
// WithGroup implements slog.Handler.WithGroup.
func (h *CustomTextHandler) WithGroup(name string) slog.Handler {
	// Skip empty group names
	if name == "" {
		return h
	}

	// Create a new handler with the same settings
	newHandler := &CustomTextHandler{
		out:       h.out,
		opts:      h.opts,
		attrOrder: h.attrOrder,
		attrs:     append([]slog.Attr{}, h.attrs...),             // Copy existing attributes
		groups:    append(append([]string{}, h.groups...), name), // Add the new group
	}

	return newHandler
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
