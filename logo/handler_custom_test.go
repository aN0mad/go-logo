package logger

import (
	"bytes"
	"context"
	"log/slog"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestNewCustomTextHandler(t *testing.T) {
	// Suppress log output for this test
	defer SuppressLogOutput(t)()

	var buf bytes.Buffer
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	handler := NewCustomTextHandler(&buf, opts)

	if handler == nil {
		t.Fatal("NewCustomTextHandler returned nil")
	}

	customHandler, ok := handler.(*CustomTextHandler)
	if !ok {
		t.Fatal("Handler is not a CustomTextHandler")
	}

	if customHandler.out != &buf {
		t.Error("Handler output writer not set correctly")
	}

	if customHandler.opts != opts {
		t.Error("Handler options not set correctly")
	}
}

func TestCustomTextHandler_Enabled(t *testing.T) {
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

			handler := &CustomTextHandler{
				out:       &buf,
				opts:      opts,
				attrOrder: attrOrder,
			}

			if got := handler.Enabled(context.Background(), tt.checkLevel); got != tt.want {
				t.Errorf("CustomTextHandler.Enabled() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCustomTextHandler_Handle(t *testing.T) {
	// Suppress log output for this test
	defer SuppressLogOutput(t)()

	fixedTime := time.Date(2023, 1, 2, 3, 4, 5, 0, time.UTC)

	tests := []struct {
		name         string
		opts         *slog.HandlerOptions
		record       slog.Record
		wantContains []string
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
			wantContains: []string{
				"time=2023-01-02T03:04:05.000Z",
				"level=INFO",
				"msg=test message",
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
			wantContains: []string{
				"level=WARN",
				"msg=warning",
				"user=john",
				"count=42",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			handler := &CustomTextHandler{
				out:       &buf,
				opts:      tt.opts,
				attrOrder: attrOrder,
			}

			err := handler.Handle(context.Background(), tt.record)
			if err != nil {
				t.Errorf("CustomTextHandler.Handle() error = %v", err)
				return
			}

			got := buf.String()
			for _, want := range tt.wantContains {
				if !strings.Contains(got, want) {
					t.Errorf("Output %q should contain %q", got, want)
				}
			}
		})
	}
}

func TestCustomTextHandler_Handle_WithSource(t *testing.T) {
	// Suppress log output for this test
	defer SuppressLogOutput(t)()

	var buf bytes.Buffer
	var levelVar slog.LevelVar
	levelVar.Set(slog.LevelInfo)

	opts := &slog.HandlerOptions{
		Level:     &levelVar,
		AddSource: true,
	}

	handler := &CustomTextHandler{
		out:       &buf,
		opts:      opts,
		attrOrder: []string{"time", "level", "msg", "source"},
	}

	// Create a record with a valid PC value from the current call site
	pc, _, _, _ := runtime.Caller(0)
	r := slog.NewRecord(time.Now(), slog.LevelInfo, "test with source", pc)

	// When the record has a valid PC value, Handle should include source information
	err := handler.Handle(context.Background(), r)
	if err != nil {
		t.Fatalf("Handle() error = %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "source=") {
		t.Errorf("Output should contain source information when AddSource=true, got: %q", output)
	}
}

func TestCustomTextHandler_WithAttrs(t *testing.T) {
	// Suppress log output for this test
	defer SuppressLogOutput(t)()

	var buf bytes.Buffer
	var levelVar slog.LevelVar
	levelVar.Set(slog.LevelInfo)

	opts := &slog.HandlerOptions{
		Level: &levelVar,
	}

	handler := &CustomTextHandler{
		out:       &buf,
		opts:      opts,
		attrOrder: attrOrder,
	}

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

func TestCustomTextHandler_WithGroup(t *testing.T) {
	// Suppress log output for this test
	defer SuppressLogOutput(t)()

	var buf bytes.Buffer
	var levelVar slog.LevelVar
	levelVar.Set(slog.LevelInfo)

	opts := &slog.HandlerOptions{
		Level: &levelVar,
	}

	handler := &CustomTextHandler{
		out:       &buf,
		opts:      opts,
		attrOrder: attrOrder,
		attrs:     []slog.Attr{}, // Initialize empty attributes
		groups:    []string{},    // Initialize empty groups
	}

	newHandler := handler.WithGroup("test_group")

	// The WithGroup implementation now correctly creates a new handler
	if newHandler == handler {
		t.Error("WithGroup() should return a new handler with the group information")
	}

	// Check that the new handler has the group information
	customHandler, ok := newHandler.(*CustomTextHandler)
	if !ok {
		t.Fatal("WithGroup() didn't return a CustomTextHandler")
	}

	if len(customHandler.groups) != 1 || customHandler.groups[0] != "test_group" {
		t.Errorf("WithGroup() didn't correctly add the group name, groups: %v", customHandler.groups)
	}
}

func TestLevelToString(t *testing.T) {
	// Suppress log output for this test
	defer SuppressLogOutput(t)()

	tests := []struct {
		level slog.Level
		want  string
	}{
		{LevelTrace, "TRACE"},
		{slog.LevelDebug, "DEBUG"},
		{slog.LevelInfo, "INFO"},
		{slog.LevelWarn, "WARN"},
		{slog.LevelError, "ERROR"},
		{LevelFatal, "FATAL"},
		{slog.Level(999), ""}, // Unknown level
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := levelToString(tt.level); got != tt.want {
				t.Errorf("levelToString(%v) = %q, want %q", tt.level, got, tt.want)
			}
		})
	}
}
