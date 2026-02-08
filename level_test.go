package logo

import (
	"testing"
)

func TestLevel_String(t *testing.T) {
	tests := []struct {
		level    Level
		expected string
	}{
		{LevelTrace, "TRACE"},
		{LevelDebug, "DEBUG"},
		{LevelInfo, "INFO"},
		{LevelWarn, "WARN"},
		{LevelError, "ERROR"},
		{LevelFatal, "FATAL"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.level.String(); got != tt.expected {
				t.Errorf("Level.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestParseLevel(t *testing.T) {
	tests := []struct {
		input    string
		expected Level
	}{
		{"TRACE", LevelTrace},
		{"trace", LevelTrace},
		{"DEBUG", LevelDebug},
		{"debug", LevelDebug},
		{"INFO", LevelInfo},
		{"info", LevelInfo},
		{"WARN", LevelWarn},
		{"warn", LevelWarn},
		{"WARNING", LevelWarn},
		{"ERROR", LevelError},
		{"error", LevelError},
		{"FATAL", LevelFatal},
		{"fatal", LevelFatal},
		{"unknown", LevelInfo}, // defaults to info
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := ParseLevel(tt.input); got != tt.expected {
				t.Errorf("ParseLevel(%v) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}
