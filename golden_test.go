package logo

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestGoldenJSON validates that the JSON output schema remains stable.
// This test uses golden files to detect any breaking changes to the JSON schema.
func TestGoldenJSON(t *testing.T) {
	// Fixed timestamp for deterministic output
	fixedTime := time.Date(2026, 2, 6, 12, 30, 45, 123456789, time.UTC)

	tests := []struct {
		name   string
		entry  *Entry
		pretty bool
		source bool
	}{
		{
			name: "basic_info",
			entry: &Entry{
				Timestamp: fixedTime,
				Level:     LevelInfo,
				Message:   "application started",
				Logger:    "main",
				Fields:    map[string]interface{}{},
			},
			pretty: false,
			source: false,
		},
		{
			name: "with_fields",
			entry: &Entry{
				Timestamp: fixedTime,
				Level:     LevelWarn,
				Message:   "high latency detected",
				Logger:    "http",
				Fields: map[string]interface{}{
					"duration_ms": 150,
					"endpoint":    "/api/users",
					"method":      "POST",
				},
			},
			pretty: false,
			source: false,
		},
		{
			name: "with_source",
			entry: &Entry{
				Timestamp: fixedTime,
				Level:     LevelError,
				Message:   "connection failed",
				Logger:    "database",
				Fields: map[string]interface{}{
					"error": "timeout",
					"retry": 3,
				},
				Caller: &Caller{
					File: "db.go",
					Line: 42,
				},
			},
			pretty: false,
			source: true,
		},
		{
			name: "pretty_formatted",
			entry: &Entry{
				Timestamp: fixedTime,
				Level:     LevelDebug,
				Message:   "cache miss",
				Logger:    "cache",
				Fields: map[string]interface{}{
					"key":   "user:123",
					"ttl":   300,
					"store": "redis",
				},
			},
			pretty: true,
			source: false,
		},
		{
			name: "all_levels",
			entry: &Entry{
				Timestamp: fixedTime,
				Level:     LevelError,
				Message:   "test message",
				Logger:    "test",
				Fields:    map[string]interface{}{"level_test": true},
			},
			pretty: false,
			source: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create formatter with options
			opts := []JSONFormatterOption{}
			if tt.pretty {
				opts = append(opts, WithPrettyJSON(true))
			}
			if tt.source {
				opts = append(opts, WithJSONSource(true))
			}

			formatter := NewJSONFormatter(opts...)
			output, err := formatter.Format(tt.entry)
			if err != nil {
				t.Fatalf("failed to format: %v", err)
			}

			// Validate it's valid JSON
			var parsed map[string]interface{}
			if err := json.Unmarshal(output, &parsed); err != nil {
				t.Fatalf("output is not valid JSON: %v", err)
			}

			// Check golden file
			goldenFile := filepath.Join("testdata", "golden", tt.name+".json")

			// Update mode: set environment variable UPDATE_GOLDEN=1 to update golden files
			if os.Getenv("UPDATE_GOLDEN") == "1" {
				// Ensure directory exists
				dir := filepath.Dir(goldenFile)
				if err := os.MkdirAll(dir, 0755); err != nil {
					t.Fatalf("failed to create golden directory: %v", err)
				}

				if err := os.WriteFile(goldenFile, output, 0644); err != nil {
					t.Fatalf("failed to write golden file: %v", err)
				}
				t.Logf("Updated golden file: %s", goldenFile)
				return
			}

			// Normal test mode: compare against golden file
			expected, err := os.ReadFile(goldenFile)
			if err != nil {
				t.Fatalf("failed to read golden file: %v (run with UPDATE_GOLDEN=1 to create)", err)
			}

			if !bytes.Equal(output, expected) {
				t.Errorf("output does not match golden file\nGot:\n%s\nExpected:\n%s", output, expected)
			}
		})
	}
}

// TestSchemaStability validates that required fields are always present.
func TestSchemaStability(t *testing.T) {
	formatter := NewJSONFormatter()

	entry := NewEntry(LevelInfo, "test", "test", map[string]interface{}{
		"custom": "value",
	})

	output, err := formatter.Format(entry)
	if err != nil {
		t.Fatalf("failed to format: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(output, &parsed); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	// All required fields must be present
	requiredFields := []string{"ts", "level", "msg", "logger"}
	for _, field := range requiredFields {
		if _, ok := parsed[field]; !ok {
			t.Errorf("required field %q is missing", field)
		}
	}

	// Check field types
	if _, ok := parsed["ts"].(string); !ok {
		t.Error("ts must be a string")
	}

	if _, ok := parsed["level"].(string); !ok {
		t.Error("level must be a string")
	}

	if _, ok := parsed["msg"].(string); !ok {
		t.Error("msg must be a string")
	}

	if _, ok := parsed["logger"].(string); !ok {
		t.Error("logger must be a string")
	}
}
