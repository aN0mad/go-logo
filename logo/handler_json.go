package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"runtime"
	"slices"
	"strings"
)

// JSONHandler is a custom handler that produces JSON output with configurable formatting
type JSONHandler struct {
	out         io.Writer
	opts        *slog.HandlerOptions
	prettyPrint bool
	attrOrder   []string
}

// NewJSONHandler creates a new JSON handler with optional pretty printing
func NewJSONHandler(out io.Writer, opts *slog.HandlerOptions, pretty bool) slog.Handler {
	// Define the attribute order: time, level, msg, source, followed by other attrs
	// attrOrder := []string{"time", "level", "msg", "source"}

	return &JSONHandler{
		out:         out,
		opts:        opts,
		prettyPrint: pretty,
		attrOrder:   attrOrder,
	}
}

// Enabled implements slog.Handler interface
func (h *JSONHandler) Enabled(ctx context.Context, level slog.Level) bool {
	minLevel := h.opts.Level.Level()
	return level >= minLevel
}

// Handle implements slog.Handler interface
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
				if lastSlash := strings.LastIndex(shortFile, "/"); lastSlash >= 0 {
					shortFile = shortFile[lastSlash+1:]
				}
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
		if a.Key == "time" || a.Key == "level" || a.Key == "msg" || a.Key == "source" {
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

// WithAttrs implements slog.Handler interface
func (h *JSONHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	// Create a new handler with the same settings
	return NewJSONHandler(h.out, h.opts, h.prettyPrint)
}

// WithGroup implements slog.Handler interface
func (h *JSONHandler) WithGroup(name string) slog.Handler {
	// Groups are not fully implemented in this simple handler
	return h
}
