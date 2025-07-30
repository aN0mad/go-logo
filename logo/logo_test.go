package logger

import (
	"bytes"
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// Mock for os.Exit function to test fatal logging
var osExit = os.Exit

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
		if logLevel != slog.LevelInfo {
			t.Errorf("Default log level is %v, want %v", logLevel, slog.LevelInfo)
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
		if logLevel != slog.LevelDebug {
			t.Errorf("Log level is %v, want %v", logLevel, slog.LevelDebug)
		}

		// Verify stack traces are enabled
		if !includeStackTraces {
			t.Error("Trace should be enabled")
		}

		// Verify source is enabled
		if !includeSource {
			t.Error("Source should be enabled")
		}

		// Verify colors are disabled
		if colorEnabled {
			t.Error("Colors should be disabled")
		}

		// Test the logger works by writing to our buffer
		L().Info("Test message")

		if !strings.Contains(buf.String(), "Test message") {
			t.Errorf("Log output doesn't contain expected message: %q", buf.String())
		}
	})
}

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

	// Initialize the logger with file output
	Init(
		AddFileOutput(logFile, 10, 3, 30, false),
	)

	// Log some test messages
	L().Info("This is a file test message", "key", "value")

	// Close the logger to ensure all file handles are closed and flushed
	err = Close()
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

func TestFatal(t *testing.T) {
	// Suppress log output for this test
	defer SuppressLogOutput(t)()

	// Save the original osExit and restore after test
	originalOsExit := osExit
	defer func() { osExit = originalOsExit }()

	exitCode := 0
	osExit = func(code int) {
		exitCode = code
		// Don't actually exit
	}

	var buf bytes.Buffer
	Init(SetConsoleOutput(&buf))

	// Call Fatal which should call osExit(1)
	L().Fatal("Fatal message")

	// Check that osExit was called with code 1
	if exitCode != 1 {
		t.Errorf("Fatal() called osExit with code %d, want 1", exitCode)
	}

	// Check message was logged
	if !strings.Contains(buf.String(), "Fatal message") {
		t.Errorf("Fatal message not logged: %q", buf.String())
	}
}

func TestLoggerOptions(t *testing.T) {
	// Test each logger option individually

	t.Run("SetLevel", func(t *testing.T) {
		// Reset logger state
		logLevel = slog.LevelInfo

		// Apply the option
		opt := SetLevel(slog.LevelError)
		opt()

		if logLevel != slog.LevelError {
			t.Errorf("SetLevel() didn't set the level correctly, got %v, want %v", logLevel, slog.LevelError)
		}
	})

	t.Run("DisableColors", func(t *testing.T) {
		// Reset logger state
		colorEnabled = true

		// Apply the option
		opt := DisableColors()
		opt()

		if colorEnabled {
			t.Error("DisableColors() didn't disable colors")
		}
	})

	t.Run("EnableStackTraces", func(t *testing.T) {
		// Reset logger state
		includeStackTraces = false

		// Apply the option
		opt := EnableStackTraces()
		opt()

		if !includeStackTraces {
			t.Error("EnableStackTraces() didn't enable stack traces")
		}
	})

	t.Run("AddSource", func(t *testing.T) {
		// Reset logger state
		includeSource = false

		// Apply the option
		opt := AddSource()
		opt()

		if !includeSource {
			t.Error("AddSource() didn't enable source")
		}
	})

	t.Run("UseJSON", func(t *testing.T) {
		// Reset logger state
		useJSONFormat = false
		jsonPretty = false

		// Apply the option with pretty printing
		opt := UseJSON(true)
		opt()

		if !useJSONFormat {
			t.Error("UseJSON() didn't enable JSON format")
		}

		if !jsonPretty {
			t.Error("UseJSON(true) didn't enable pretty printing")
		}

		// Test without pretty printing
		useJSONFormat = false
		jsonPretty = false

		opt = UseJSON(false)
		opt()

		if !useJSONFormat {
			t.Error("UseJSON() didn't enable JSON format")
		}

		if jsonPretty {
			t.Error("UseJSON(false) didn't disable pretty printing")
		}
	})
}
