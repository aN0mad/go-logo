package logo

import (
	"testing"
	"time"
)

func TestNewEntry(t *testing.T) {
	fields := map[string]interface{}{
		"key1": "value1",
		"key2": 42,
	}

	entry := NewEntry(LevelInfo, "test message", "test-logger", fields)

	if entry.Level != LevelInfo {
		t.Errorf("expected level INFO, got %v", entry.Level)
	}

	if entry.Message != "test message" {
		t.Errorf("expected message 'test message', got %v", entry.Message)
	}

	if entry.Logger != "test-logger" {
		t.Errorf("expected logger 'test-logger', got %v", entry.Logger)
	}

	if len(entry.Fields) != 2 {
		t.Errorf("expected 2 fields, got %d", len(entry.Fields))
	}

	if entry.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}

	// Check timestamp is recent (within 1 second)
	now := time.Now().UTC()
	if now.Sub(entry.Timestamp) > time.Second {
		t.Errorf("timestamp is too old: %v", entry.Timestamp)
	}
}

func TestEntry_WithCaller(t *testing.T) {
	entry := NewEntry(LevelInfo, "test", "logger", nil)
	caller := &Caller{
		File:     "test.go",
		Line:     42,
		Function: "TestFunc",
	}

	entry = entry.WithCaller(caller)

	if entry.Caller == nil {
		t.Fatal("expected caller to be set")
	}

	if entry.Caller.File != "test.go" {
		t.Errorf("expected file 'test.go', got %v", entry.Caller.File)
	}

	if entry.Caller.Line != 42 {
		t.Errorf("expected line 42, got %v", entry.Caller.Line)
	}
}

func TestEntry_Clone(t *testing.T) {
	fields := map[string]interface{}{
		"key1": "value1",
	}

	original := NewEntry(LevelInfo, "test", "logger", fields)
	clone := original.Clone()

	// Verify fields match
	if clone.Level != original.Level {
		t.Error("cloned level doesn't match")
	}

	if clone.Message != original.Message {
		t.Error("cloned message doesn't match")
	}

	// Verify it's a deep copy
	clone.Fields["key1"] = "modified"
	if original.Fields["key1"] == "modified" {
		t.Error("modifying clone affected original")
	}
}
