package logo

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestColorSupport_Auto(t *testing.T) {
	// Test with stdout
	cs := NewColorSupport(os.Stdout, ColorAuto)
	// Can't assert specific value since it depends on environment
	// Just verify it doesn't panic
	_ = cs.Enabled()
}

func TestColorSupport_Always(t *testing.T) {
	cs := NewColorSupport(os.Stdout, ColorAlways)
	if !cs.Enabled() {
		t.Error("ColorAlways should always enable colors")
	}
}

func TestColorSupport_Never(t *testing.T) {
	cs := NewColorSupport(os.Stdout, ColorNever)
	if cs.Enabled() {
		t.Error("ColorNever should never enable colors")
	}
}

func TestColorSupport_Colorize(t *testing.T) {
	// With colors enabled
	cs := NewColorSupport(os.Stdout, ColorAlways)
	result := cs.Colorize("test", colorRed)
	if !strings.Contains(result, "test") {
		t.Error("colorized text should contain original text")
	}
	if !strings.Contains(result, colorRed) {
		t.Error("colorized text should contain color code")
	}
	if !strings.Contains(result, colorReset) {
		t.Error("colorized text should contain reset code")
	}

	// With colors disabled
	cs = NewColorSupport(os.Stdout, ColorNever)
	result = cs.Colorize("test", colorRed)
	if result != "test" {
		t.Errorf("without colors, text should be unchanged, got %q", result)
	}
}

func TestLevelColor(t *testing.T) {
	tests := []struct {
		level Level
		want  string
	}{
		{LevelTrace, colorCyan},
		{LevelDebug, colorGray},
		{LevelInfo, colorBlue},
		{LevelWarn, colorYellow},
		{LevelError, colorRed},
		{LevelFatal, colorMagenta},
	}

	for _, tt := range tests {
		t.Run(tt.level.String(), func(t *testing.T) {
			got := LevelColor(tt.level)
			if got != tt.want {
				t.Errorf("LevelColor(%v) = %v, want %v", tt.level, got, tt.want)
			}
		})
	}
}

func TestPrettyFormatter_Format(t *testing.T) {
	cs := NewColorSupport(os.Stdout, ColorNever) // Disable colors for testing
	formatter := NewPrettyFormatter(cs)

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

	// Check for message
	if !strings.Contains(str, "test message") {
		t.Error("output should contain message")
	}

	// Check for level
	if !strings.Contains(str, "INFO") {
		t.Error("output should contain level")
	}

	// Check for fields (sorted alphabetically)
	if !strings.Contains(str, "id=123") {
		t.Error("output should contain id field")
	}
	if !strings.Contains(str, "user=test-user") {
		t.Error("output should contain user field")
	}

	// Should end with newline
	if !strings.HasSuffix(str, "\n") {
		t.Error("output should end with newline")
	}
}

func TestPrettyFormatter_WithColors(t *testing.T) {
	cs := NewColorSupport(os.Stdout, ColorAlways)
	formatter := NewPrettyFormatter(cs)

	entry := NewEntry(LevelError, "error message", "test", nil)

	output, err := formatter.Format(entry)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	str := string(output)

	// Should contain ANSI codes when colors are enabled
	if !strings.Contains(str, "\033[") {
		t.Error("output with colors should contain ANSI escape codes")
	}
}

func TestPrettyFormatter_WithSource(t *testing.T) {
	cs := NewColorSupport(os.Stdout, ColorNever)
	formatter := NewPrettyFormatter(cs, WithSourceInfo(true))

	entry := NewEntry(LevelInfo, "test", "test", nil)
	entry.WithCaller(&Caller{
		File: "test.go",
		Line: 42,
	})

	output, err := formatter.Format(entry)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	str := string(output)

	if !strings.Contains(str, "test.go:42") {
		t.Errorf("output should contain source location, got: %s", str)
	}
}

func TestPrettyFormatter_FieldFormatting(t *testing.T) {
	cs := NewColorSupport(os.Stdout, ColorNever)
	formatter := NewPrettyFormatter(cs)

	tests := []struct {
		name   string
		fields map[string]interface{}
		want   string
	}{
		{
			name:   "string with spaces",
			fields: map[string]interface{}{"msg": "hello world"},
			want:   `msg="hello world"`,
		},
		{
			name:   "string without spaces",
			fields: map[string]interface{}{"name": "john"},
			want:   "name=john",
		},
		{
			name:   "number",
			fields: map[string]interface{}{"count": 42},
			want:   "count=42",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entry := NewEntry(LevelInfo, "test", "test", tt.fields)
			output, err := formatter.Format(entry)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			str := string(output)
			if !strings.Contains(str, tt.want) {
				t.Errorf("output should contain %q, got: %s", tt.want, str)
			}
		})
	}
}

func TestConsoleHandler(t *testing.T) {
	buf := &bytes.Buffer{}
	handler := NewConsoleHandlerWithWriter(buf, ColorNever, LevelInfo)

	entry := NewEntry(LevelInfo, "test message", "test", nil)
	handled := handler.Handle(entry)

	if !handled {
		t.Error("handler should have handled the entry")
	}

	if buf.Len() == 0 {
		t.Error("handler should have written output")
	}

	output := buf.String()
	if !strings.Contains(output, "test message") {
		t.Error("output should contain message")
	}
}

func TestConsoleHandler_LevelFiltering(t *testing.T) {
	buf := &bytes.Buffer{}
	handler := NewConsoleHandlerWithWriter(buf, ColorNever, LevelWarn)

	// Info should be filtered
	entry := NewEntry(LevelInfo, "info message", "test", nil)
	handled := handler.Handle(entry)

	if handled {
		t.Error("info message should have been filtered")
	}

	if buf.Len() != 0 {
		t.Error("buffer should be empty for filtered message")
	}

	// Warn should pass through
	entry = NewEntry(LevelWarn, "warn message", "test", nil)
	handled = handler.Handle(entry)

	if !handled {
		t.Error("warn message should have been handled")
	}

	if buf.Len() == 0 {
		t.Error("buffer should contain warn message")
	}
}

func TestWithConsole(t *testing.T) {
	log := New(WithConsole())

	if log.router == nil {
		t.Fatal("router should be set")
	}

	handlers := log.router.Handlers()
	if len(handlers) == 0 {
		t.Error("router should have at least one handler")
	}

	// Verify it's a console handler
	_, ok := handlers[0].(*ConsoleHandler)
	if !ok {
		t.Error("handler should be a ConsoleHandler")
	}
}
