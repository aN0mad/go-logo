package logo

import (
	"bytes"
	"strings"
	"testing"
)

func TestRouter_Route(t *testing.T) {
	buf1 := &bytes.Buffer{}
	buf2 := &bytes.Buffer{}

	handler1 := NewWriterHandler(buf1, NewTextFormatter(), LevelDebug)
	handler2 := NewWriterHandler(buf2, NewTextFormatter(), LevelInfo)

	router := NewRouter(handler1, handler2)

	// Debug message - should go to handler1 only
	entry := NewEntry(LevelDebug, "debug message", "test", nil)
	router.Route(entry)

	if buf1.Len() == 0 {
		t.Error("handler1 should have received debug message")
	}

	if buf2.Len() != 0 {
		t.Error("handler2 should not have received debug message")
	}

	// Info message - should go to both
	buf1.Reset()
	buf2.Reset()

	entry = NewEntry(LevelInfo, "info message", "test", nil)
	router.Route(entry)

	if buf1.Len() == 0 {
		t.Error("handler1 should have received info message")
	}

	if buf2.Len() == 0 {
		t.Error("handler2 should have received info message")
	}
}

func TestRouter_AddHandler(t *testing.T) {
	router := NewRouter()

	if len(router.Handlers()) != 0 {
		t.Error("expected empty router")
	}

	buf := &bytes.Buffer{}
	handler := NewWriterHandler(buf, NewTextFormatter(), LevelInfo)
	router.AddHandler(handler)

	if len(router.Handlers()) != 1 {
		t.Error("expected 1 handler")
	}

	// Test that added handler receives messages
	entry := NewEntry(LevelInfo, "test", "test", nil)
	router.Route(entry)

	if buf.Len() == 0 {
		t.Error("handler should have received message")
	}
}

func TestRouter_Close(t *testing.T) {
	buf := &bytes.Buffer{}
	handler := NewWriterHandler(buf, NewTextFormatter(), LevelInfo)
	router := NewRouter(handler)

	err := router.Close()
	if err != nil {
		t.Errorf("unexpected error closing router: %v", err)
	}
}

func TestTextFormatter_Format(t *testing.T) {
	formatter := NewTextFormatter()
	fields := map[string]interface{}{
		"user": "test-user",
		"id":   123,
	}

	entry := NewEntry(LevelInfo, "test message", "test-logger", fields)

	output, err := formatter.Format(entry)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	str := string(output)

	// Check that output contains expected elements
	if !strings.Contains(str, "[INFO]") {
		t.Error("output should contain [INFO]")
	}

	if !strings.Contains(str, "test message") {
		t.Error("output should contain message")
	}

	if !strings.Contains(str, "user=test-user") {
		t.Error("output should contain user field")
	}

	if !strings.Contains(str, "id=123") {
		t.Error("output should contain id field")
	}

	// Should end with newline
	if !strings.HasSuffix(str, "\n") {
		t.Error("output should end with newline")
	}
}

func TestWriterHandler_Handle(t *testing.T) {
	buf := &bytes.Buffer{}
	handler := NewWriterHandler(buf, NewTextFormatter(), LevelInfo)

	// Info message should be handled
	entry := NewEntry(LevelInfo, "info message", "test", nil)
	handled := handler.Handle(entry)

	if !handled {
		t.Error("info message should have been handled")
	}

	if buf.Len() == 0 {
		t.Error("buffer should contain output")
	}

	// Debug message should be filtered
	buf.Reset()
	entry = NewEntry(LevelDebug, "debug message", "test", nil)
	handled = handler.Handle(entry)

	if handled {
		t.Error("debug message should have been filtered")
	}

	if buf.Len() != 0 {
		t.Error("buffer should be empty")
	}
}

func TestBaseHandler_Enabled(t *testing.T) {
	handler := NewBaseHandler(LevelWarn)

	tests := []struct {
		level   Level
		enabled bool
	}{
		{LevelDebug, false},
		{LevelInfo, false},
		{LevelWarn, true},
		{LevelError, true},
	}

	for _, tt := range tests {
		t.Run(tt.level.String(), func(t *testing.T) {
			if got := handler.Enabled(tt.level); got != tt.enabled {
				t.Errorf("Enabled(%v) = %v, want %v", tt.level, got, tt.enabled)
			}
		})
	}
}
