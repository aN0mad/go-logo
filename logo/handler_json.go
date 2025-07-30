// Package logger provides functionality for structured logging.
//
// This file contains the JSON handler implementation which formats log messages
// as JSON objects, with optional pretty-printing.
package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"runtime"
	"slices"
)

// JSONHandler is a custom handler that produces JSON output with configurable formatting.
// It implements the slog.Handler interface and supports ordering attributes
// and pretty printing options.
type JSONHandler struct {
	out         io.Writer
	opts        *slog.HandlerOptions
	prettyPrint bool
	attrOrder   []string
}

// NewJSONHandler creates a new JSON handler with optional pretty printing.
//
// Parameters:
//   - out: The io.Writer where JSON log entries will be written
//   - opts: Handler options including log level and attribute replacements
//   - pretty: Whether to format the JSON with indentation for better readability
//
// Returns a slog.Handler implementation
func NewJSONHandler(out io.Writer, opts *slog.HandlerOptions, pretty bool) slog.Handler {

	return &JSONHandler{
		out:         out,
		opts:        opts,
		prettyPrint: pretty,
		attrOrder:   attrOrder, // Use the global attrOrder defined in this package
	}
}

// Enabled implements slog.Handler interface.
// It checks if the given log level should be processed based on the configured minimum level.
//
// Parameters:
//   - ctx: The context for the logging operation
//   - level: The log level to check
//
// Returns true if the log level should be processed, false otherwise
func (h *JSONHandler) Enabled(ctx context.Context, level slog.Level) bool {
	minLevel := h.opts.Level.Level()
	return level >= minLevel
}

// Handle implements slog.Handler interface.
// It processes a log record and outputs it as JSON.
//
// Parameters:
//   - ctx: The context for the logging operation
//   - r: The log record to process
//
// Returns any error encountered during formatting or writing
func (h *JSONHandler) Handle(ctx context.Context, r slog.Record) error {
	// Create an ordered map to maintain attribute order
	orderedMap := make(map[string]interface{})

	// Add standard attributes in desired order
	orderedMap["time"] = r.Time.Format("2006-01-02T15:04:05.000Z07:00")
	orderedMap["level"] = levelToString(r.Level)
	orderedMap["msg"] = r.Message

	// Add source if requested
	if h.opts.AddSource {
		if source := r.PC; source != 0 {
			fs := runtime.CallersFrames([]uintptr{source})
			frame, _ := fs.Next()
			if frame.File != "" {
				shortFile := frame.File
				// if lastSlash := strings.LastIndex(shortFile, "/"); lastSlash >= 0 { // Removed short source for full path
				// 	shortFile = shortFile[lastSlash+1:]
				// }
				orderedMap["source"] = fmt.Sprintf("%s:%d", shortFile, frame.Line)
			}
		}
	}

	// Add attributes
	otherAttrs := make(map[string]interface{})

	r.Attrs(func(a slog.Attr) bool {
		// Apply attribute transformations if specified
		if h.opts.ReplaceAttr != nil {
			a = h.opts.ReplaceAttr(nil, a)
		}

		// Skip empty attributes
		if a.Equal(slog.Attr{}) {
			return true
		}

		// Skip attributes we've already handled
		if slices.Contains(h.attrOrder, a.Key) {
			return true
		}

		// Add the attribute to our collection
		otherAttrs[a.Key] = a.Value.Any()
		return true
	})

	// Add remaining attributes in alphabetical order
	var keys []string
	for k := range otherAttrs {
		keys = append(keys, k)
	}
	slices.Sort(keys)

	for _, k := range keys {
		orderedMap[k] = otherAttrs[k]
	}

	// Convert to JSON
	var jsonData []byte
	var err error

	if h.prettyPrint {
		jsonData, err = json.MarshalIndent(orderedMap, "", "  ")
	} else {
		jsonData, err = json.Marshal(orderedMap)
	}

	if err != nil {
		return err
	}

	// Write to output
	_, err = h.out.Write(append(jsonData, '\n'))
	return err
}

// WithAttrs implements slog.Handler interface.
// It returns a new handler with the given attributes.
//
// Parameters:
//   - attrs: The attributes to add to the handler
//
// Returns a new handler instance with the attributes
func (h *JSONHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	// Create a new handler with the same settings
	return NewJSONHandler(h.out, h.opts, h.prettyPrint)
}

// WithGroup implements slog.Handler interface.
// It returns a handler that adds the given group name to the attribute key path.
//
// Parameters:
//   - name: The group name
//
// Returns a handler that adds the group name to the attribute key path
func (h *JSONHandler) WithGroup(name string) slog.Handler {
	// Groups are not fully implemented in this simple handler
	return h
}
