package logo

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestIsReservedKey(t *testing.T) {
	tests := []struct {
		key      string
		reserved bool
	}{
		{"ts", true},
		{"level", true},
		{"msg", true},
		{"logger", true},
		{"source", true},
		{"caller", true},
		{"user_id", false},
		{"custom_field", false},
		{"timestamp", false},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			got := IsReservedKey(tt.key)
			if got != tt.reserved {
				t.Errorf("IsReservedKey(%q) = %v, want %v", tt.key, got, tt.reserved)
			}
		})
	}
}

func TestValidateFields(t *testing.T) {
	tests := []struct {
		name    string
		fields  map[string]interface{}
		wantErr bool
	}{
		{
			name:    "valid fields",
			fields:  map[string]interface{}{"user_id": 123, "name": "test"},
			wantErr: false,
		},
		{
			name:    "reserved key ts",
			fields:  map[string]interface{}{"ts": "value"},
			wantErr: true,
		},
		{
			name:    "reserved key level",
			fields:  map[string]interface{}{"level": "INFO"},
			wantErr: true,
		},
		{
			name:    "reserved key msg",
			fields:  map[string]interface{}{"msg": "message"},
			wantErr: true,
		},
		{
			name:    "empty fields",
			fields:  map[string]interface{}{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFields(tt.fields)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateFields() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFilterReservedKeys(t *testing.T) {
	fields := map[string]interface{}{
		"ts":      "should be removed",
		"level":   "should be removed",
		"msg":     "should be removed",
		"user_id": 123,
		"name":    "test",
	}

	filtered := FilterReservedKeys(fields)

	// Should only have user fields
	if len(filtered) != 2 {
		t.Errorf("expected 2 fields after filtering, got %d", len(filtered))
	}

	if _, ok := filtered["ts"]; ok {
		t.Error("ts should be filtered out")
	}

	if _, ok := filtered["user_id"]; !ok {
		t.Error("user_id should be present")
	}

	if _, ok := filtered["name"]; !ok {
		t.Error("name should be present")
	}
}

func TestJSONFormatter_Format(t *testing.T) {
	formatter := NewJSONFormatter()

	fields := map[string]interface{}{
		"user_id": 123,
		"name":    "test",
	}

	entry := NewEntry(LevelInfo, "test message", "test-logger", fields)

	output, err := formatter.Format(entry)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Parse the JSON to validate structure
	var parsed map[string]interface{}
	if err := json.Unmarshal(output, &parsed); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	// Validate required fields
	if _, ok := parsed["ts"]; !ok {
		t.Error("missing required field: ts")
	}

	if parsed["level"] != "INFO" {
		t.Errorf("level = %v, want INFO", parsed["level"])
	}

	if parsed["msg"] != "test message" {
		t.Errorf("msg = %v, want 'test message'", parsed["msg"])
	}

	if parsed["logger"] != "test-logger" {
		t.Errorf("logger = %v, want 'test-logger'", parsed["logger"])
	}

	// Validate user fields
	if parsed["user_id"] != float64(123) { // JSON numbers are float64
		t.Errorf("user_id = %v, want 123", parsed["user_id"])
	}

	if parsed["name"] != "test" {
		t.Errorf("name = %v, want 'test'", parsed["name"])
	}
}

func TestJSONFormatter_ReservedKeyFiltering(t *testing.T) {
	formatter := NewJSONFormatter()

	// Try to override reserved keys
	fields := map[string]interface{}{
		"ts":      "user-provided-timestamp",
		"level":   "user-provided-level",
		"msg":     "user-provided-msg",
		"user_id": 123,
	}

	entry := NewEntry(LevelWarn, "actual message", "test-logger", fields)

	output, err := formatter.Format(entry)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(output, &parsed); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	// Schema fields should NOT be overridden
	if parsed["level"] != "WARN" {
		t.Error("level should be WARN, not user-provided value")
	}

	if parsed["msg"] != "actual message" {
		t.Error("msg should be actual message, not user-provided value")
	}

	// User field should still be present
	if parsed["user_id"] != float64(123) {
		t.Error("user_id should be preserved")
	}
}

func TestJSONFormatter_PrettyPrint(t *testing.T) {
	formatter := NewJSONFormatter(WithPrettyJSON(true))

	entry := NewEntry(LevelInfo, "test", "test", map[string]interface{}{"key": "value"})

	output, err := formatter.Format(entry)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	str := string(output)

	// Pretty-printed JSON should have newlines and indentation
	if !strings.Contains(str, "\n") {
		t.Error("pretty JSON should contain newlines")
	}

	if !strings.Contains(str, "  ") {
		t.Error("pretty JSON should contain indentation")
	}
}

func TestJSONFormatter_WithSource(t *testing.T) {
	formatter := NewJSONFormatter(WithJSONSource(true))

	entry := NewEntry(LevelInfo, "test", "test", nil)
	entry.WithCaller(&Caller{
		File: "test.go",
		Line: 42,
	})

	output, err := formatter.Format(entry)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(output, &parsed); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	source, ok := parsed["source"]
	if !ok {
		t.Fatal("source field should be present")
	}

	if source != "test.go:42" {
		t.Errorf("source = %v, want 'test.go:42'", source)
	}
}

func TestJSONFormatter_TimestampFormat(t *testing.T) {
	formatter := NewJSONFormatter()

	now := time.Date(2026, 2, 6, 12, 30, 45, 123456789, time.UTC)
	entry := &Entry{
		Timestamp: now,
		Level:     LevelInfo,
		Message:   "test",
		Logger:    "test",
		Fields:    map[string]interface{}{},
	}

	output, err := formatter.Format(entry)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(output, &parsed); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	ts := parsed["ts"].(string)

	// Should be RFC3339Nano format
	_, err = time.Parse(time.RFC3339Nano, ts)
	if err != nil {
		t.Errorf("timestamp format invalid: %v", err)
	}

	// Should be in UTC
	if !strings.HasSuffix(ts, "Z") {
		t.Error("timestamp should be in UTC (end with Z)")
	}
}

func TestJSONHandler(t *testing.T) {
	buf := &bytes.Buffer{}
	handler := NewJSONHandler(buf, LevelInfo)

	entry := NewEntry(LevelInfo, "test message", "test", map[string]interface{}{"user": "john"})
	handled := handler.Handle(entry)

	if !handled {
		t.Error("handler should have handled the entry")
	}

	if buf.Len() == 0 {
		t.Error("handler should have written output")
	}

	// Validate JSON structure
	var parsed map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &parsed); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	if parsed["msg"] != "test message" {
		t.Error("message not in output")
	}
}

func TestJSONHandler_LevelFiltering(t *testing.T) {
	buf := &bytes.Buffer{}
	handler := NewJSONHandler(buf, LevelWarn)

	// Info should be filtered
	entry := NewEntry(LevelInfo, "info", "test", nil)
	handled := handler.Handle(entry)

	if handled {
		t.Error("info message should have been filtered")
	}

	if buf.Len() != 0 {
		t.Error("buffer should be empty")
	}

	// Warn should pass
	entry = NewEntry(LevelWarn, "warning", "test", nil)
	handled = handler.Handle(entry)

	if !handled {
		t.Error("warn message should have been handled")
	}

	if buf.Len() == 0 {
		t.Error("buffer should contain output")
	}
}
