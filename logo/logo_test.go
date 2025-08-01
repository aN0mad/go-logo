package logo

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestInit tests the initialization of the global logger with various options.
// It verifies that the logger is properly configured with the provided options.
//
// Parameters:
//   - t: The testing instance used for assertions and test control
func TestInit(t *testing.T) {
	// Suppress log output for this test
	defer SuppressLogOutput(t)()

	// Test with default settings
	t.Run("default settings", func(t *testing.T) {
		Init()

		// Verify the logger is initialized
		if logger == nil {
			t.Error("Logger was not initialized")
		}

		// Verify default level is INFO
		if logger.ctx.logLevel != slog.LevelInfo {
			t.Errorf("Default log level is %v, want %v", logger.ctx.logLevel, slog.LevelInfo)
		}
	})

	// Test with custom settings
	t.Run("custom settings", func(t *testing.T) {
		var buf bytes.Buffer

		Init(
			EnableStackTraces(),
			SetLevel(slog.LevelDebug),
			DisableColors(),
			SetConsoleOutput(&buf),
			AddSource(),
		)

		// Verify the logger is initialized
		if logger == nil {
			t.Error("Logger was not initialized")
		}

		// Verify custom level is set
		if logger.ctx.logLevel != slog.LevelDebug {
			t.Errorf("Log level is %v, want %v", logger.ctx.logLevel, slog.LevelDebug)
		}

		// Verify stack traces are enabled
		if !logger.ctx.includeStackTraces {
			t.Error("Stack traces should be enabled")
		}

		// Verify source is enabled
		if !logger.ctx.includeSource {
			t.Error("Source should be enabled")
		}

		// Verify colors are disabled
		if logger.ctx.colorEnabled {
			t.Error("Colors should be disabled")
		}

		// Test the logger works by writing to our buffer
		L().Info("Test message")

		if !strings.Contains(buf.String(), "Test message") {
			t.Errorf("Log output doesn't contain expected message: %q", buf.String())
		}
	})
}

// TestLoggerMethods tests the various logging methods provided by the Logger.
// It verifies that each method correctly formats and outputs messages at the appropriate level.
//
// Parameters:
//   - t: The testing instance used for assertions and test control
func TestLoggerMethods(t *testing.T) {
	// Suppress log output for this test
	defer SuppressLogOutput(t)()

	var buf bytes.Buffer

	Init(
		SetLevel(LevelTrace), // Set to lowest level to capture everything
		SetConsoleOutput(&buf),
		DisableColors(), // For easier testing of output
	)

	// Test different log levels
	testCases := []struct {
		name    string
		logFunc func(string, ...any)
		level   string
		message string
	}{
		{
			name: "trace level",
			logFunc: func(msg string, args ...any) {
				L().Trace(msg, args...)
			},
			level:   "TRACE",
			message: "trace message",
		},
		{
			name: "debug level",
			logFunc: func(msg string, args ...any) {
				L().Debug(msg, args...)
			},
			level:   "DEBUG",
			message: "debug message",
		},
		{
			name: "info level",
			logFunc: func(msg string, args ...any) {
				L().Info(msg, args...)
			},
			level:   "INFO",
			message: "info message",
		},
		{
			name: "warn level",
			logFunc: func(msg string, args ...any) {
				L().Warn(msg, args...)
			},
			level:   "WARN",
			message: "warn message",
		},
		{
			name: "error level",
			logFunc: func(msg string, args ...any) {
				L().Error(msg, args...)
			},
			level:   "ERROR",
			message: "error message",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			buf.Reset()
			tc.logFunc(tc.message)

			output := buf.String()
			if !strings.Contains(output, tc.level) {
				t.Errorf("Output %q should contain level %q", output, tc.level)
			}

			if !strings.Contains(output, tc.message) {
				t.Errorf("Output %q should contain message %q", output, tc.message)
			}
		})
	}

	// Test with attributes
	t.Run("with attributes", func(t *testing.T) {
		buf.Reset()
		L().Info("message with attributes", "key1", "value1", "key2", 42)

		output := buf.String()
		if !strings.Contains(output, "key1=value1") {
			t.Errorf("Output %q should contain attribute key1=value1", output)
		}

		if !strings.Contains(output, "key2=42") {
			t.Errorf("Output %q should contain attribute key2=42", output)
		}
	})
}

// TestContextLogger tests creating and using loggers with context information.
// It verifies that context values are correctly included in log output.
//
// Parameters:
//   - t: The testing instance used for assertions and test control
func TestContextLogger(t *testing.T) {
	// Suppress log output for this test
	defer SuppressLogOutput(t)()

	var buf bytes.Buffer

	// Initialize logger with settings that match the example
	Init(
		SetConsoleOutput(&buf),
		DisableColors(), // For easier testing of output
		AddSource(),     // Add this to match the example setup
	)

	// Create context with request_id
	ctx := context.WithValue(context.Background(), "request_id", "test-id-123")

	// Create a test-specific context logger (similar to the example)
	loggerFromContext := func(ctx context.Context) *Logger {
		baseLogger := L()

		// Add application-specific context values
		if requestID, ok := ctx.Value("request_id").(string); ok {
			baseLogger = &Logger{
				Logger: baseLogger.Logger.With("request_id", requestID),
			}
		}

		return baseLogger
	}

	// Get logger with context
	logWithCtx := loggerFromContext(ctx)

	// Log a message with the context logger
	logWithCtx.Info("test message")

	// Debug the output to see what's actually being logged
	output := buf.String()
	t.Logf("Logged output: %q", output)

	// Check for the request_id attribute in the output
	if !strings.Contains(output, "request_id=test-id-123") {
		t.Errorf("Output %q should contain request_id attribute", output)
	}

	// Test with context that doesn't have request_id
	buf.Reset()
	emptyCtx := context.Background()
	loggerFromContext(emptyCtx).Info("another message")

	// Should log normally without request_id
	output = buf.String()
	if strings.Contains(output, "request_id") {
		t.Errorf("Output %q should not contain request_id attribute", output)
	}
}

// TestFileOutput tests logging to file outputs.
// It verifies that log messages are correctly written to the configured log file.
//
// Parameters:
//   - t: The testing instance used for assertions and test control
func TestFileOutput(t *testing.T) {
	// Suppress log output for this test
	defer SuppressLogOutput(t)()

	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "logger-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	logFile := filepath.Join(tempDir, "test.log")

	// Create a logger with file output
	fileLogger := NewLogger(
		AddFileOutput(logFile, 10, 3, 30, false),
	)

	// Log some test messages
	fileLogger.Info("This is a file test message", "key", "value")

	// Close the logger to ensure all file handles are closed and flushed
	err = fileLogger.Close()
	if err != nil {
		t.Errorf("Failed to close logger: %v", err)
	}

	// Verify the log file exists and contains the message
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		t.Errorf("Log file was not created at %s", logFile)
	}

	fileContent, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	content := string(fileContent)
	if !strings.Contains(content, "This is a file test message") {
		t.Errorf("Log file doesn't contain expected message: %q", content)
	}

	if !strings.Contains(content, "key=value") {
		t.Errorf("Log file doesn't contain expected attribute: %q", content)
	}
}

// TestJSONOutput tests JSON formatted logging.
// It verifies that log messages are correctly formatted as JSON with the expected structure.
//
// Parameters:
//   - t: The testing instance used for assertions and test control
func TestJSONOutput(t *testing.T) {
	// Suppress log output for this test
	defer SuppressLogOutput(t)()

	var buf bytes.Buffer

	// Initialize logger with JSON output
	Init(
		UseJSON(false), // Not pretty-printed
		SetConsoleOutput(&buf),
	)

	// Log a test message
	L().Info("JSON test message", "key", "value")

	// Check for JSON format indicators
	output := buf.String()
	if !strings.Contains(output, `"msg":"JSON test message"`) {
		t.Errorf("Output %q should be in JSON format with msg field", output)
	}

	if !strings.Contains(output, `"key":"value"`) {
		t.Errorf("Output %q should contain JSON attribute", output)
	}

	// Test pretty-printed JSON
	buf.Reset()
	Init(
		UseJSON(true), // Pretty-printed
		SetConsoleOutput(&buf),
	)

	L().Info("Pretty JSON message", "key", "value")

	output = buf.String()
	// Pretty JSON should have line breaks and indentation
	if !strings.Contains(output, "\n") || !strings.Contains(output, "  ") {
		t.Errorf("Output %q should be in pretty-printed JSON format", output)
	}
}

// TestChannelOutput tests logging to channel outputs.
// It verifies that log messages are correctly sent to the configured channel.
//
// Parameters:
//   - t: The testing instance used for assertions and test control
func TestChannelOutput(t *testing.T) {
	// Suppress log output for this test
	defer SuppressLogOutput(t)()

	// Create a channel to receive log messages
	ch := make(chan string, 5)

	// Initialize the logger with channel output
	Init(
		AddChannelOutput(ch),
	)

	// Log a test message
	L().Info("Channel test message")

	// Check if the message was received on the channel
	select {
	case msg := <-ch:
		if !strings.Contains(msg, "Channel test message") {
			t.Errorf("Channel message %q doesn't contain expected text", msg)
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("Timeout waiting for message on channel")
	}
}

// TestFatal tests the Fatal logging method.
// It verifies that the method correctly logs the message and calls os.Exit with code 1.
//
// Parameters:
//   - t: The testing instance used for assertions and test control
func TestFatal(t *testing.T) {
	// Suppress log output for this test
	defer SuppressLogOutput(t)()

	// Save the original osExit function
	originalOsExit := osExit

	// Create a local variable to capture the exit code
	var exitCode int

	// Replace osExit with our test version before the test
	osExit = func(code int) {
		exitCode = code
		// Don't actually exit - just record the exit code
	}

	// Restore the original function when the test completes
	defer func() {
		osExit = originalOsExit
	}()

	var buf bytes.Buffer
	Init(
		SetConsoleOutput(&buf),
		DisableColors(),
	)

	// Call Fatal which should trigger our mocked osExit
	L().Fatal("Fatal message")

	// Check that our mock was called with code 1
	if exitCode != 1 {
		t.Errorf("Fatal() called osExit with code %d, want 1", exitCode)
	}

	// Check message was logged
	if !strings.Contains(buf.String(), "Fatal message") {
		t.Errorf("Fatal message not logged: %q", buf.String())
	}
}

// TestLoggerOptions tests the functional options used to configure the logger.
// It verifies that each option correctly modifies the logger's configuration.
//
// Parameters:
//   - t: The testing instance used for assertions and test control
func TestLoggerOptions(t *testing.T) {
	// Test each logger option with the context-based approach

	t.Run("SetLevel", func(t *testing.T) {
		ctx := &loggerContext{
			logLevel: slog.LevelInfo,
		}

		// Apply the option
		opt := SetLevel(slog.LevelError)
		opt(ctx)

		if ctx.logLevel != slog.LevelError {
			t.Errorf("SetLevel() didn't set the level correctly, got %v, want %v",
				ctx.logLevel, slog.LevelError)
		}
	})

	t.Run("DisableColors", func(t *testing.T) {
		ctx := &loggerContext{
			colorEnabled: true,
		}

		// Apply the option
		opt := DisableColors()
		opt(ctx)

		if ctx.colorEnabled {
			t.Error("DisableColors() didn't disable colors")
		}
	})

	t.Run("EnableStackTraces", func(t *testing.T) {
		ctx := &loggerContext{
			includeStackTraces: false,
		}

		// Apply the option
		opt := EnableStackTraces()
		opt(ctx)

		if !ctx.includeStackTraces {
			t.Error("EnableStackTraces() didn't enable stack traces")
		}
	})

	t.Run("AddSource", func(t *testing.T) {
		ctx := &loggerContext{
			includeSource: false,
		}

		// Apply the option
		opt := AddSource()
		opt(ctx)

		if !ctx.includeSource {
			t.Error("AddSource() didn't enable source")
		}
	})

	t.Run("UseJSON", func(t *testing.T) {
		ctx := &loggerContext{
			useJSONFormat: false,
			jsonPretty:    false,
		}

		// Apply the option with pretty printing
		opt := UseJSON(true)
		opt(ctx)

		if !ctx.useJSONFormat {
			t.Error("UseJSON() didn't enable JSON format")
		}

		if !ctx.jsonPretty {
			t.Error("UseJSON(true) didn't enable pretty printing")
		}

		// Test without pretty printing
		ctx = &loggerContext{
			useJSONFormat: false,
			jsonPretty:    false,
		}

		opt = UseJSON(false)
		opt(ctx)

		if !ctx.useJSONFormat {
			t.Error("UseJSON() didn't enable JSON format")
		}

		if ctx.jsonPretty {
			t.Error("UseJSON(false) didn't disable pretty printing")
		}
	})
}

// TestEnableLogLevelTrace tests the EnableLogLevelTrace option.
// It verifies that this option correctly sets the logger's level to LevelTrace.
//
// Parameters:
//   - t: The testing instance used for assertions and test control
func TestEnableLogLevelTrace(t *testing.T) {
	// Create a context with default level
	ctx := &loggerContext{
		logLevel: slog.LevelInfo,
	}

	// Apply the option to the context
	opt := EnableLogLevelTrace()
	opt(ctx)

	// Check that the context's log level was updated
	if ctx.logLevel != LevelTrace {
		t.Errorf("EnableLogLevelTrace() didn't set the level correctly, got %v, want %v",
			ctx.logLevel, LevelTrace)
	}
}

// TestNormalizeAttrs tests the normalizeAttrs function.
// It verifies that the function correctly normalizes various attribute formats.
//
// Parameters:
//   - t: The testing instance used for assertions and test control
func TestNormalizeAttrs(t *testing.T) {
	tests := []struct {
		name     string
		args     []any
		expected int // Expected number of attributes
	}{
		{
			name:     "empty args",
			args:     []any{},
			expected: 0,
		},
		{
			name:     "string key-value pairs",
			args:     []any{"key1", "value1", "key2", 42},
			expected: 2,
		},
		{
			name:     "slog.Attr arguments",
			args:     []any{slog.String("key", "value")},
			expected: 1,
		},
		{
			name:     "mixed argument types",
			args:     []any{slog.Int("count", 5), "message", "hello"},
			expected: 2,
		},
		{
			name:     "malformed arguments (odd count)",
			args:     []any{"key1", "value1", "orphan_key"},
			expected: 2, // Should include the orphan key with (MISSING) value
		},
		{
			name:     "non-string keys",
			args:     []any{123, "value"},
			expected: 0, // Should skip non-string keys
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attrs := normalizeAttrs(tt.args...)
			if len(attrs) != tt.expected {
				t.Errorf("normalizeAttrs() returned %d attributes, want %d",
					len(attrs), tt.expected)
			}
		})
	}
}

// TestSetFileHandlerForTesting tests the SetFileHandlerForTesting option.
// It verifies that this option correctly configures a writer for testing.
//
// Parameters:
//   - t: The testing instance used for assertions and test control
func TestSetFileHandlerForTesting(t *testing.T) {
	// Suppress log output
	defer SuppressLogOutput(t)()

	var buf bytes.Buffer

	// Create a context to apply the option to
	ctx := &loggerContext{}

	// Apply the option to add the buffer to the context's outputs
	opt := SetFileHandlerForTesting(&buf)
	opt(ctx)

	// Check that the buffer was added to the outputs
	if len(ctx.outputs) != 1 {
		t.Errorf("SetFileHandlerForTesting() didn't add output, got %d outputs", len(ctx.outputs))
	}

	// Create a logger with the test buffer
	testLogger := NewLogger(
		SetFileHandlerForTesting(&buf),
	)

	// Test by writing a log message
	testLogger.Info("Test message")

	// Check if the message was logged to our buffer
	if !strings.Contains(buf.String(), "Test message") {
		t.Errorf("Message not logged to test writer: %q", buf.String())
	}
}

// TestFileAndConsoleLogger tests logging to both file and console simultaneously.
// It verifies that log messages are correctly written to both outputs.
//
// Parameters:
//   - t: The testing instance used for assertions and test control
func TestFileAndConsoleLogger(t *testing.T) {
	// Suppress log output for this test
	defer SuppressLogOutput(t)()

	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "logger-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	logFile := filepath.Join(tempDir, "multi.log")
	t.Logf("Using log file path: %s", logFile)

	// Create buffer for console output
	var buf bytes.Buffer

	// Create a logger with BOTH console and file output properly configured
	multiLogger := NewLogger(
		AddFileOutput(logFile, 10, 3, 30, false),
		SetLevel(slog.LevelDebug),
		SetConsoleOutput(&buf), // Use the helper function to set console output
	)

	// Write multiple messages to increase chances of flushing
	multiLogger.Info("Test both outputs", "type", "multi")
	multiLogger.Debug("Another test message")
	multiLogger.Error("Error message for testing")

	// Force sync by closing the logger
	err = multiLogger.Close()
	if err != nil {
		t.Errorf("Failed to close logger: %v", err)
	}

	// Add delay to ensure filesystem operations complete
	time.Sleep(200 * time.Millisecond)

	// Check if the file exists and get its size
	fileInfo, err := os.Stat(logFile)
	if err != nil {
		if os.IsNotExist(err) {
			t.Fatalf("Log file doesn't exist: %s", logFile)
		}
		t.Fatalf("Error checking log file: %v", err)
	}
	t.Logf("Log file exists with size: %d bytes", fileInfo.Size())

	// List directory contents for debugging
	files, _ := os.ReadDir(tempDir)
	var fileNames []string
	for _, f := range files {
		info, _ := f.Info()
		fileNames = append(fileNames, fmt.Sprintf("%s (%d bytes)", f.Name(), info.Size()))
	}
	t.Logf("Directory contents: %v", fileNames)

	// Check console output
	consoleOutput := buf.String()
	t.Logf("Console output length: %d bytes", len(consoleOutput))
	t.Logf("Console output content: %q", consoleOutput)

	if !strings.Contains(consoleOutput, "Test both outputs") {
		t.Errorf("Console output doesn't contain the message, got: %q", consoleOutput)
	}

	// Check file output
	fileContent, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	fileOutput := string(fileContent)
	t.Logf("File content (%d bytes): %q", len(fileOutput), fileOutput)

	if !strings.Contains(fileOutput, "Test both outputs") {
		t.Error("File output doesn't contain the message")
	}
	if !strings.Contains(fileOutput, "type=multi") {
		t.Error("File output doesn't contain the attribute")
	}
}

// TestMultipleLoggers tests creating and using multiple independent logger instances.
// It verifies that each logger maintains its own configuration.
//
// Parameters:
//   - t: The testing instance used for assertions and test control
func TestMultipleLoggers(t *testing.T) {
	// Suppress log output for this test
	defer SuppressLogOutput(t)()

	// Create buffers for capturing output
	var buf1, buf2 bytes.Buffer

	// Create first logger with debug level and direct console output to buffer
	logger1 := NewLogger(
		SetLevel(slog.LevelDebug),
		DisableColors(),
		SetConsoleOutput(&buf1), // Use helper function instead of manual replacement
	)

	// Create second logger with error level and direct console output to buffer
	logger2 := NewLogger(
		SetLevel(slog.LevelError),
		DisableColors(),
		SetConsoleOutput(&buf2), // Use helper function instead of manual replacement
	)

	// Test that levels work independently
	logger1.Debug("Debug message") // Should appear in buf1
	logger1.Error("Error message") // Should appear in buf1

	logger2.Debug("Debug message") // Should NOT appear in buf2
	logger2.Error("Error message") // Should appear in buf2

	// Check logger1 output
	output1 := buf1.String()
	t.Logf("Logger1 output: %q", output1) // Add debugging output
	if !strings.Contains(output1, "Debug message") {
		t.Error("Logger1 should log debug messages")
	}
	if !strings.Contains(output1, "Error message") {
		t.Error("Logger1 should log error messages")
	}

	// Check logger2 output
	output2 := buf2.String()
	t.Logf("Logger2 output: %q", output2) // Add debugging output
	if strings.Contains(output2, "Debug message") {
		t.Error("Logger2 should NOT log debug messages")
	}
	if !strings.Contains(output2, "Error message") {
		t.Error("Logger2 should log error messages")
	}

	// Test that modifying one logger doesn't affect the other
	SetLoggerLevel(logger1, slog.LevelError)

	// Reset buffers
	buf1.Reset()
	buf2.Reset()

	logger1.Debug("Another debug message") // Should NOT appear in buf1 now
	logger2.Debug("Another debug message") // Should NOT appear in buf2

	// Add some messages that should appear
	logger1.Error("Error after level change")
	logger2.Error("Error from logger2")

	t.Logf("Logger1 output after level change: %q", buf1.String())

	if strings.Contains(buf1.String(), "Another debug message") {
		t.Error("Logger1 should NOT log debug messages after level change")
	}

	if !strings.Contains(buf1.String(), "Error after level change") {
		t.Error("Logger1 should still log error messages after level change")
	}
}
