package logo

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewLumberjackWriter(t *testing.T) {
	// Suppress log output for this test
	defer SuppressLogOutput(t)()

	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "logger-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	testCases := []struct {
		name     string
		filename string
		maxSize  int
		backups  int
		maxAge   int
		compress bool
	}{
		{
			name:     "basic config",
			filename: filepath.Join(tempDir, "test.log"),
			maxSize:  10,
			backups:  3,
			maxAge:   30,
			compress: false,
		},
		{
			name:     "with compression",
			filename: filepath.Join(tempDir, "compressed.log"),
			maxSize:  5,
			backups:  5,
			maxAge:   14,
			compress: true,
		},
		{
			name:     "minimal config",
			filename: filepath.Join(tempDir, "minimal.log"),
			maxSize:  1,
			backups:  0,
			maxAge:   1,
			compress: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			writer := NewLumberjackWriter(tc.filename, tc.maxSize, tc.backups, tc.maxAge, tc.compress)

			if writer == nil {
				t.Fatal("NewLumberjackWriter returned nil")
			}

			// Check if the writer was configured correctly
			if writer.Filename != tc.filename {
				t.Errorf("Incorrect filename, got: %s, want: %s", writer.Filename, tc.filename)
			}

			if writer.MaxSize != tc.maxSize {
				t.Errorf("Incorrect maxSize, got: %d, want: %d", writer.MaxSize, tc.maxSize)
			}

			if writer.MaxBackups != tc.backups {
				t.Errorf("Incorrect backups, got: %d, want: %d", writer.MaxBackups, tc.backups)
			}

			if writer.MaxAge != tc.maxAge {
				t.Errorf("Incorrect maxAge, got: %d, want: %d", writer.MaxAge, tc.maxAge)
			}

			if writer.Compress != tc.compress {
				t.Errorf("Incorrect compress setting, got: %v, want: %v", writer.Compress, tc.compress)
			}
		})
	}
}

func TestLumberjackWriter_Write(t *testing.T) {
	// Suppress log output for this test
	defer SuppressLogOutput(t)()

	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "logger-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	logFile := filepath.Join(tempDir, "test.log")
	writer := NewLumberjackWriter(logFile, 10, 3, 30, false)

	// Write some test data
	testData := "Test log message"
	n, err := writer.Write([]byte(testData))
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	if n != len(testData) {
		t.Errorf("Write returned %d bytes written, want %d", n, len(testData))
	}

	// Check if file was created and contains the test data
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		t.Errorf("Log file was not created at %s", logFile)
	}

	fileContent, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	if string(fileContent) != testData {
		t.Errorf("Log file content = %q, want %q", string(fileContent), testData)
	}

	// Close the writer to ensure file handles are released
	err = writer.Close()
	if err != nil {
		t.Errorf("Failed to close writer: %v", err)
	}
}

func TestLumberjackWriter_Rotation(t *testing.T) {
	// Suppress log output for this test
	defer SuppressLogOutput(t)()

	// Skip this test in CI environments where we might not have write permissions
	// or where file operations might be restricted
	if os.Getenv("CI") != "" {
		t.Skip("Skipping rotation test in CI environment")
	}

	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "logger-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	logFile := filepath.Join(tempDir, "rotation.log")

	// Create a writer with a very small max size to trigger rotation
	writer := NewLumberjackWriter(logFile, 1, 3, 30, false)

	// Write enough data to trigger at least one rotation
	data := make([]byte, 1024*1024) // 1MB of data
	for i := range data {
		data[i] = 'A' // Fill with 'A' characters
	}

	// Write the data multiple times to ensure rotation occurs
	for i := 0; i < 3; i++ {
		_, err = writer.Write(data)
		if err != nil {
			t.Fatalf("Write %d failed: %v", i, err)
		}

		// Force rotation by closing and reopening
		if i < 2 {
			writer.Close()
			writer.Rotate()
		}
	}

	// Close the writer to ensure all data is flushed
	err = writer.Close()
	if err != nil {
		t.Errorf("Failed to close writer: %v", err)
	}

	// Check if the log file exists
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		t.Errorf("Main log file was not created at %s", logFile)
	}

	// List all files in the directory
	entries, err := os.ReadDir(tempDir)
	if err != nil {
		t.Fatalf("Error reading directory: %v", err)
	}

	// Count rotated files (any file that starts with "rotation-")
	var rotatedFiles []string
	for _, entry := range entries {
		name := entry.Name()
		if strings.HasPrefix(name, "rotation-") {
			rotatedFiles = append(rotatedFiles, name)
		}
	}

	// We should have at least one rotated file
	if len(rotatedFiles) == 0 {
		t.Error("No rotated log files found")

		t.Log("Directory contents:")
		for _, entry := range entries {
			info, _ := entry.Info()
			t.Logf("  %s (%d bytes)", entry.Name(), info.Size())
		}
	} else {
		t.Logf("Found %d rotated log files: %v", len(rotatedFiles), rotatedFiles)
	}
}

func TestIntegrationWithLogger(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "logger-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	logFile := filepath.Join(tempDir, "integration.log")

	// Initialize the logger with file output only
	Init(
		SetLevel(LevelTrace),
		DisableConsole(), // Turn off console output
		AddFileOutput(logFile, 10, 3, 30, false),
	)

	// Log some test messages
	L().Info("This is an info message")
	L().Debug("This is a debug message")
	L().Error("This is an error message", "error_code", 123)

	// Close the logger to ensure all data is flushed
	err = Close()
	if err != nil {
		t.Errorf("Failed to close logger: %v", err)
	}

	// Check if the log file exists
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		t.Fatalf("Log file was not created at %s", logFile)
	}

	// Read the log file content
	fileContent, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	content := string(fileContent)

	// Check for the presence of log messages in the file
	for _, expected := range []string{"INFO", "DEBUG", "ERROR", "This is an info message", "This is a debug message", "This is an error message"} {
		if !strings.Contains(content, expected) {
			t.Errorf("Log file doesn't contain expected text %q", expected)
		}
	}
}
