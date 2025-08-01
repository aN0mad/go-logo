package logo

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"runtime"
	"strings"
	"testing"
	"time"
)

// TestNewJSONHandler tests the creation of a JSONHandler.
// It verifies that the handler is properly initialized with the provided parameters.
//
// Parameters:
//   - t: The testing instance used for assertions and test control
func TestNewJSONHandler(t *testing.T) {
	// Suppress log output for this test
	defer SuppressLogOutput(t)()

	var buf bytes.Buffer
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	// Test with pretty print disabled
	handler := NewJSONHandler(&buf, opts, false)
	if handler == nil {
		t.Fatal("NewJSONHandler returned nil")
	}

	jsonHandler, ok := handler.(*JSONHandler)
	if !ok {
		t.Fatal("Handler is not a JSONHandler")
	}

	if jsonHandler.out != &buf {
		t.Error("Handler output writer not set correctly")
	}

	if jsonHandler.opts != opts {
		t.Error("Handler options not set correctly")
	}

	if jsonHandler.prettyPrint {
		t.Error("prettyPrint should be false when not requested")
	}

	// Test with pretty print enabled
	handler = NewJSONHandler(&buf, opts, true)
	jsonHandler, ok = handler.(*JSONHandler)
	if !ok {
		t.Fatal("Handler is not a JSONHandler")
	}

	if !jsonHandler.prettyPrint {
		t.Error("prettyPrint should be true when requested")
	}

	// Test that attrOrder was initialized correctly
	expectedOrder := attrOrder
	if len(jsonHandler.attrOrder) != len(expectedOrder) {
		t.Errorf("attrOrder length = %d, want %d", len(jsonHandler.attrOrder), len(expectedOrder))
	}

	for i, attr := range expectedOrder {
		if i < len(jsonHandler.attrOrder) && jsonHandler.attrOrder[i] != attr {
			t.Errorf("attrOrder[%d] = %q, want %q", i, jsonHandler.attrOrder[i], attr)
		}
	}
}

// TestJSONHandler_Enabled tests the Enabled method of JSONHandler.
// It verifies that the method correctly determines whether a log level should be processed.
//
// Parameters:
//   - t: The testing instance used for assertions and test control
func TestJSONHandler_Enabled(t *testing.T) {
	// Suppress log output for this test
	defer SuppressLogOutput(t)()

	tests := []struct {
		name       string
		minLevel   slog.Level
		checkLevel slog.Level
		want       bool
	}{
		{
			name:       "info enabled at info level",
			minLevel:   slog.LevelInfo,
			checkLevel: slog.LevelInfo,
			want:       true,
		},
		{
			name:       "warn enabled at info level",
			minLevel:   slog.LevelInfo,
			checkLevel: slog.LevelWarn,
			want:       true,
		},
		{
			name:       "debug disabled at info level",
			minLevel:   slog.LevelInfo,
			checkLevel: slog.LevelDebug,
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			var levelVar slog.LevelVar
			levelVar.Set(tt.minLevel)

			opts := &slog.HandlerOptions{
				Level: &levelVar,
			}

			handler := NewJSONHandler(&buf, opts, false)

			if got := handler.Enabled(context.Background(), tt.checkLevel); got != tt.want {
				t.Errorf("JSONHandler.Enabled() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestJSONHandler_Handle tests the Handle method of JSONHandler.
// It verifies that the method correctly formats log records as JSON.
//
// Parameters:
//   - t: The testing instance used for assertions and test control
func TestJSONHandler_Handle(t *testing.T) {
	// Suppress log output for this test
	defer SuppressLogOutput(t)()

	fixedTime := time.Date(2023, 1, 2, 3, 4, 5, 0, time.UTC)

	tests := []struct {
		name      string
		opts      *slog.HandlerOptions
		record    slog.Record
		pretty    bool
		checkJSON map[string]interface{}
	}{
		{
			name: "basic message",
			opts: &slog.HandlerOptions{
				Level: slog.LevelInfo,
			},
			record: func() slog.Record {
				r := slog.NewRecord(fixedTime, slog.LevelInfo, "test message", 0)
				return r
			}(),
			pretty: false,
			checkJSON: map[string]interface{}{
				"time":  "2023-01-02T03:04:05.000Z",
				"level": "INFO",
				"msg":   "test message",
			},
		},
		{
			name: "with custom attributes",
			opts: &slog.HandlerOptions{
				Level: slog.LevelWarn,
			},
			record: func() slog.Record {
				r := slog.NewRecord(fixedTime, slog.LevelWarn, "warning", 0)
				r.AddAttrs(slog.String("user", "john"), slog.Int("count", 42))
				return r
			}(),
			pretty: true,
			checkJSON: map[string]interface{}{
				"time":  "2023-01-02T03:04:05.000Z",
				"level": "WARN",
				"msg":   "warning",
				"user":  "john",
				"count": float64(42), // JSON numbers are parsed as float64
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			handler := NewJSONHandler(&buf, tt.opts, tt.pretty)

			err := handler.Handle(context.Background(), tt.record)
			if err != nil {
				t.Errorf("JSONHandler.Handle() error = %v", err)
				return
			}

			got := buf.String()

			// Check JSON formatting based on pretty setting
			if tt.pretty && !strings.Contains(got, "\n") {
				t.Errorf("Expected pretty JSON with line breaks, got: %s", got)
			}

			// Parse the JSON output
			var parsed map[string]interface{}
			if err := json.Unmarshal([]byte(got), &parsed); err != nil {
				t.Fatalf("Failed to parse JSON output: %v", err)
			}

			// Check required fields
			for k, v := range tt.checkJSON {
				if parsed[k] != v {
					t.Errorf("JSON field %q = %v, want %v", k, parsed[k], v)
				}
			}
		})
	}
}

// TestJSONHandler_Handle_WithSource tests the Handle method with source information.
// It verifies that source information is correctly included in the JSON output when enabled.
//
// Parameters:
//   - t: The testing instance used for assertions and test control
func TestJSONHandler_Handle_WithSource(t *testing.T) {
	// Suppress log output for this test
	defer SuppressLogOutput(t)()

	var buf bytes.Buffer
	var levelVar slog.LevelVar
	levelVar.Set(slog.LevelInfo)

	opts := &slog.HandlerOptions{
		Level:     &levelVar,
		AddSource: true,
	}

	handler := NewJSONHandler(&buf, opts, false)

	// Create a record with a valid PC value from the current call site
	pc, _, _, _ := runtime.Caller(0)
	r := slog.NewRecord(time.Now(), slog.LevelInfo, "test with source", pc)

	// When the record has a valid PC value, Handle should include source information
	err := handler.Handle(context.Background(), r)
	if err != nil {
		t.Fatalf("Handle() error = %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &parsed); err != nil {
		t.Fatalf("Failed to parse JSON output: %v", err)
	}

	if _, exists := parsed["source"]; !exists {
		t.Errorf("Output should contain 'source' field when AddSource=true, got: %v", parsed)
	}
}

// TestJSONHandler_WithAttrs tests the WithAttrs method of JSONHandler.
// It verifies that the method correctly creates a new handler with the provided attributes.
//
// Parameters:
//   - t: The testing instance used for assertions and test control
func TestJSONHandler_WithAttrs(t *testing.T) {
	// Suppress log output for this test
	defer SuppressLogOutput(t)()

	var buf bytes.Buffer
	var levelVar slog.LevelVar
	levelVar.Set(slog.LevelInfo)

	opts := &slog.HandlerOptions{
		Level: &levelVar,
	}

	handler := NewJSONHandler(&buf, opts, false)

	newAttrs := []slog.Attr{
		slog.String("key1", "value1"),
		slog.Int("key2", 42),
	}

	newHandler := handler.WithAttrs(newAttrs)
	if newHandler == nil {
		t.Fatal("WithAttrs() returned nil")
	}

	// Verify it's a new handler instance
	if newHandler == handler {
		t.Error("WithAttrs() should return a new handler instance")
	}
}

// TestJSONHandler_WithGroup tests the WithGroup method of JSONHandler.
// It verifies the behavior of grouping in the JSON handler implementation.
//
// Parameters:
//   - t: The testing instance used for assertions and test control
func TestJSONHandler_WithGroup(t *testing.T) {
	// Suppress log output for this test
	defer SuppressLogOutput(t)()

	var buf bytes.Buffer
	var levelVar slog.LevelVar
	levelVar.Set(slog.LevelInfo)

	opts := &slog.HandlerOptions{
		Level: &levelVar,
	}

	handler := NewJSONHandler(&buf, opts, false)

	newHandler := handler.WithGroup("test_group")

	// Since WithGroup is a no-op in this implementation, it should return the same handler
	if newHandler != handler {
		t.Error("WithGroup() should return the same handler for this implementation")
	}
}

// TestHandlerConsistency tests the consistency between JSONHandler and CustomTextHandler.
// It verifies that both handlers process the same log record in a consistent manner.
//
// Parameters:
//   - t: The testing instance used for assertions and test control
func TestHandlerConsistency(t *testing.T) {
	// Suppress log output for this test
	defer SuppressLogOutput(t)()

	fixedTime := time.Date(2023, 1, 2, 3, 4, 5, 0, time.UTC)

	var textBuf, jsonBuf bytes.Buffer
	var levelVar slog.LevelVar
	levelVar.Set(slog.LevelInfo)

	opts := &slog.HandlerOptions{
		Level: &levelVar,
	}

	textHandler := NewCustomTextHandler(&textBuf, opts)
	jsonHandler := NewJSONHandler(&jsonBuf, opts, false)

	// Create a test record
	r := slog.NewRecord(fixedTime, slog.LevelInfo, "test consistency", 0)
	r.AddAttrs(slog.String("attr1", "value1"), slog.Int("attr2", 42))

	// Process with both handlers
	if err := textHandler.Handle(context.Background(), r); err != nil {
		t.Fatalf("TextHandler.Handle() error = %v", err)
	}

	if err := jsonHandler.Handle(context.Background(), r); err != nil {
		t.Fatalf("JSONHandler.Handle() error = %v", err)
	}

	// Parse JSON output
	var parsed map[string]interface{}
	if err := json.Unmarshal(jsonBuf.Bytes(), &parsed); err != nil {
		t.Fatalf("Failed to parse JSON output: %v", err)
	}

	// Check that both outputs contain the same information
	textOut := textBuf.String()

	if !strings.Contains(textOut, "level=INFO") || parsed["level"] != "INFO" {
		t.Errorf("Inconsistent level formatting")
	}

	if !strings.Contains(textOut, "msg=test consistency") || parsed["msg"] != "test consistency" {
		t.Errorf("Inconsistent message formatting")
	}

	if !strings.Contains(textOut, "attr1=value1") || parsed["attr1"] != "value1" {
		t.Errorf("Inconsistent attribute1 formatting")
	}

	if !strings.Contains(textOut, "attr2=42") || parsed["attr2"] != float64(42) {
		t.Errorf("Inconsistent attribute2 formatting")
	}
}
