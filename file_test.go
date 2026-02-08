package logo

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestDefaultRotationConfig(t *testing.T) {
	config := DefaultRotationConfig()

	if config.MaxSize != 100 {
		t.Errorf("MaxSize = %d, want 100", config.MaxSize)
	}

	if config.MaxAge != 30 {
		t.Errorf("MaxAge = %d, want 30", config.MaxAge)
	}

	if config.MaxBackups != 10 {
		t.Errorf("MaxBackups = %d, want 10", config.MaxBackups)
	}

	if !config.Compress {
		t.Error("Compress should be true")
	}

	if config.LocalTime {
		t.Error("LocalTime should be false (UTC)")
	}
}

func TestNewFileHandler(t *testing.T) {
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "test.log")

	formatter := NewJSONFormatter()
	handler, err := NewFileHandler(logFile, formatter)
	if err != nil {
		t.Fatalf("failed to create file handler: %v", err)
	}
	defer handler.Close()

	// Write a log entry
	entry := NewEntry(LevelInfo, "test message", "test", map[string]interface{}{
		"user": "test-user",
	})

	handled := handler.Handle(entry)
	if !handled {
		t.Error("handler should have handled the entry")
	}

	// Close to flush
	handler.Close()

	// Read the file
	data, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("failed to read log file: %v", err)
	}

	if len(data) == 0 {
		t.Fatal("log file is empty")
	}

	// Verify JSON format
	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("log file does not contain valid JSON: %v", err)
	}

	if parsed["msg"] != "test message" {
		t.Errorf("msg = %v, want 'test message'", parsed["msg"])
	}
}

func TestFileHandler_WithRotation(t *testing.T) {
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "test.log")

	formatter := NewJSONFormatter()
	config := &RotationConfig{
		MaxSize:    1, // 1 MB
		MaxAge:     7,
		MaxBackups: 3,
		Compress:   false, // Don't compress for easier testing
		LocalTime:  true,
	}

	handler, err := NewFileHandler(logFile, formatter, WithRotation(config))
	if err != nil {
		t.Fatalf("failed to create file handler: %v", err)
	}
	defer handler.Close()

	// Verify config was applied
	if handler.logger.MaxSize != 1 {
		t.Errorf("MaxSize = %d, want 1", handler.logger.MaxSize)
	}

	if handler.logger.MaxAge != 7 {
		t.Errorf("MaxAge = %d, want 7", handler.logger.MaxAge)
	}

	if handler.logger.MaxBackups != 3 {
		t.Errorf("MaxBackups = %d, want 3", handler.logger.MaxBackups)
	}
}

func TestFileHandler_LevelFiltering(t *testing.T) {
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "test.log")

	formatter := NewJSONFormatter()
	handler, err := NewFileHandler(logFile, formatter, WithFileLevel(LevelWarn))
	if err != nil {
		t.Fatalf("failed to create file handler: %v", err)
	}
	defer handler.Close()

	// Info should be filtered
	infoEntry := NewEntry(LevelInfo, "info message", "test", nil)
	handled := handler.Handle(infoEntry)

	if handled {
		t.Error("info message should have been filtered")
	}

	// Warn should pass
	warnEntry := NewEntry(LevelWarn, "warn message", "test", nil)
	handled = handler.Handle(warnEntry)

	if !handled {
		t.Error("warn message should have been handled")
	}

	handler.Close()

	// Read file
	data, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("failed to read log file: %v", err)
	}

	content := string(data)

	if strings.Contains(content, "info message") {
		t.Error("file should not contain filtered info message")
	}

	if !strings.Contains(content, "warn message") {
		t.Error("file should contain warn message")
	}
}

func TestFileHandler_Rotate(t *testing.T) {
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "test.log")

	formatter := NewJSONFormatter()
	handler, err := NewFileHandler(logFile, formatter)
	if err != nil {
		t.Fatalf("failed to create file handler: %v", err)
	}
	defer handler.Close()

	// Write first entry
	entry1 := NewEntry(LevelInfo, "before rotation", "test", nil)
	handler.Handle(entry1)

	// Manually trigger rotation
	if err := handler.Rotate(); err != nil {
		t.Errorf("failed to rotate: %v", err)
	}

	// Write second entry
	entry2 := NewEntry(LevelInfo, "after rotation", "test", nil)
	handler.Handle(entry2)

	handler.Close()

	// The current file should contain only the second entry
	data, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("failed to read log file: %v", err)
	}

	content := string(data)

	if !strings.Contains(content, "after rotation") {
		t.Error("current file should contain entry after rotation")
	}
}

func TestWithJSONFile(t *testing.T) {
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "app.log")

	log := New(WithJSONFile(logFile))
	defer log.Close()

	if log.router == nil {
		t.Fatal("router should be set")
	}

	handlers := log.router.Handlers()
	if len(handlers) == 0 {
		t.Fatal("should have at least one handler")
	}

	// Log a message
	log.Info("test message", F("key", "value"))

	// Close to flush
	log.Close()

	// Give it a moment to flush
	time.Sleep(10 * time.Millisecond)

	// Verify file exists and contains JSON
	data, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("failed to read log file: %v", err)
	}

	if len(data) == 0 {
		t.Fatal("log file is empty")
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("log file does not contain valid JSON: %v", err)
	}
}

func TestFileHandler_MultipleEntries(t *testing.T) {
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "test.log")

	formatter := NewJSONFormatter()
	handler, err := NewFileHandler(logFile, formatter)
	if err != nil {
		t.Fatalf("failed to create file handler: %v", err)
	}
	defer handler.Close()

	// Write multiple entries
	for i := 0; i < 5; i++ {
		entry := NewEntry(LevelInfo, "message", "test", map[string]interface{}{
			"count": i,
		})
		handler.Handle(entry)
	}

	handler.Close()

	// Read file
	data, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("failed to read log file: %v", err)
	}

	// Should have 5 lines (one per entry)
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) != 5 {
		t.Errorf("expected 5 log lines, got %d", len(lines))
	}

	// Each line should be valid JSON
	for i, line := range lines {
		var parsed map[string]interface{}
		if err := json.Unmarshal([]byte(line), &parsed); err != nil {
			t.Errorf("line %d is not valid JSON: %v", i, err)
		}
	}
}

func TestOpenFile(t *testing.T) {
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "test.log")

	file, err := OpenFile(logFile, 0644)
	if err != nil {
		t.Fatalf("failed to open file: %v", err)
	}
	defer file.Close()

	// Write some data
	_, err = file.WriteString("test\n")
	if err != nil {
		t.Fatalf("failed to write: %v", err)
	}

	file.Close()

	// Verify file exists and has correct permissions
	info, err := os.Stat(logFile)
	if err != nil {
		t.Fatalf("failed to stat file: %v", err)
	}

	if info.Mode().Perm() != 0644 {
		t.Errorf("file permissions = %o, want 0644", info.Mode().Perm())
	}
}
