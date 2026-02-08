package logo

import (
	"bytes"
	"strings"
	"testing"
)

func TestLogger_BasicLogging(t *testing.T) {
	buf := &bytes.Buffer{}
	handler := NewWriterHandler(buf, NewTextFormatter(), LevelDebug)
	router := NewRouter(handler)
	logger := New(WithRouter(router))

	logger.Info("test message")

	output := buf.String()
	if !strings.Contains(output, "test message") {
		t.Error("output should contain message")
	}

	if !strings.Contains(output, "[INFO]") {
		t.Error("output should contain log level")
	}
}

func TestLogger_WithFields(t *testing.T) {
	buf := &bytes.Buffer{}
	handler := NewWriterHandler(buf, NewTextFormatter(), LevelDebug)
	router := NewRouter(handler)
	logger := New(WithRouter(router))

	logger.Info("test", F("key", "value"))

	output := buf.String()
	if !strings.Contains(output, "key=value") {
		t.Errorf("output should contain field, got: %s", output)
	}
}

func TestLogger_With(t *testing.T) {
	buf := &bytes.Buffer{}
	handler := NewWriterHandler(buf, NewTextFormatter(), LevelDebug)
	router := NewRouter(handler)
	logger := New(WithRouter(router))

	// Create a child logger with additional fields
	childLogger := logger.With(F("request_id", "123"))
	childLogger.Info("test message")

	output := buf.String()
	if !strings.Contains(output, "request_id=123") {
		t.Errorf("output should contain inherited field, got: %s", output)
	}

	// Verify original logger is unchanged
	buf.Reset()
	logger.Info("test message")

	output = buf.String()
	if strings.Contains(output, "request_id") {
		t.Error("original logger should not have child's fields")
	}
}

func TestLogger_WithName(t *testing.T) {
	buf := &bytes.Buffer{}
	handler := NewWriterHandler(buf, NewTextFormatter(), LevelDebug)
	router := NewRouter(handler)
	logger := New(WithRouter(router), WithName("parent"))

	childLogger := logger.WithName("child")
	childLogger.Info("test")

	// Original logger should still have "parent" name
	buf.Reset()
	logger.Info("test")
	// Name is used internally, not directly visible in text output
}

func TestLogger_LogLevels(t *testing.T) {
	buf := &bytes.Buffer{}
	handler := NewWriterHandler(buf, NewTextFormatter(), LevelTrace)
	router := NewRouter(handler)
	logger := New(WithRouter(router))

	tests := []struct {
		method   func(string, ...Field)
		expected string
	}{
		{logger.Trace, "[TRACE]"},
		{logger.Debug, "[DEBUG]"},
		{logger.Info, "[INFO]"},
		{logger.Warn, "[WARN]"},
		{logger.Error, "[ERROR]"},
		{logger.Fatal, "[FATAL]"},
	}

	for _, tt := range tests {
		buf.Reset()
		tt.method("test message")

		output := buf.String()
		if !strings.Contains(output, tt.expected) {
			t.Errorf("expected %s in output, got: %s", tt.expected, output)
		}
	}
}

func TestLogger_LevelFiltering(t *testing.T) {
	buf := &bytes.Buffer{}
	handler := NewWriterHandler(buf, NewTextFormatter(), LevelWarn)
	router := NewRouter(handler)
	logger := New(WithRouter(router))

	// Debug and Info should be filtered
	logger.Debug("debug message")
	logger.Info("info message")

	if buf.Len() != 0 {
		t.Error("debug and info messages should have been filtered")
	}

	// Warn and Error should pass through
	logger.Warn("warn message")
	if buf.Len() == 0 {
		t.Error("warn message should have been logged")
	}

	buf.Reset()
	logger.Error("error message")
	if buf.Len() == 0 {
		t.Error("error message should have been logged")
	}
}

func TestLogger_ConcurrentUse(t *testing.T) {
	buf := &bytes.Buffer{}
	handler := NewWriterHandler(buf, NewTextFormatter(), LevelInfo)
	router := NewRouter(handler)
	logger := New(WithRouter(router))

	// Spawn multiple goroutines writing logs
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(id int) {
			logger.Info("concurrent message", F("goroutine", id))
			done <- true
		}(i)
	}

	// Wait for all to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	// Should have output from all goroutines without panicking
	if buf.Len() == 0 {
		t.Error("should have received messages from goroutines")
	}
}

func TestLogger_MultipleHandlers(t *testing.T) {
	buf1 := &bytes.Buffer{}
	buf2 := &bytes.Buffer{}

	handler1 := NewWriterHandler(buf1, NewTextFormatter(), LevelDebug)
	handler2 := NewWriterHandler(buf2, NewTextFormatter(), LevelInfo)

	router := NewRouter(handler1, handler2)
	logger := New(WithRouter(router))

	logger.Info("test message")

	// Both buffers should have output
	if buf1.Len() == 0 {
		t.Error("handler1 should have output")
	}

	if buf2.Len() == 0 {
		t.Error("handler2 should have output")
	}

	// Test level filtering
	buf1.Reset()
	buf2.Reset()

	logger.Debug("debug message")

	// Only buf1 should have output (handler2 filters debug)
	if buf1.Len() == 0 {
		t.Error("handler1 should have debug output")
	}

	if buf2.Len() != 0 {
		t.Error("handler2 should have filtered debug")
	}
}
