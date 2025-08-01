// Package logo provides functionality for structured logging.
//
// This file contains utility functions for testing to help with
// common tasks like suppressing log output during tests.
package logo

import (
	"io"
	"os"
	"testing"
)

// SuppressLogOutput temporarily redirects stdout and stderr to /dev/null
// during test execution and restores them afterwards. This prevents log
// messages from cluttering the test output.
//
// Parameters:
//   - t: The testing.T instance for test context
//
// Returns:
//   - func(): A cleanup function that should be deferred to restore the original output
func SuppressLogOutput(t *testing.T) func() {
	t.Helper()

	// Save the original stdout and stderr
	oldStdout := os.Stdout
	oldStderr := os.Stderr

	// Create null device for discarding output
	devNull, err := os.Open(os.DevNull)
	if err != nil {
		t.Fatalf("Failed to open %s: %v", os.DevNull, err)
	}

	// Redirect stdout and stderr to null device
	os.Stdout = devNull
	os.Stderr = devNull

	// Return a cleanup function to restore original stdout and stderr
	return func() {
		// Restore original stdout and stderr
		os.Stdout = oldStdout
		os.Stderr = oldStderr

		// Close the null device
		devNull.Close()
	}
}

// SetConsoleOutput sets a custom writer for the console output.
// This is useful for testing to redirect console logs to a null device or buffer.
//
// Parameters:
//   - w: The io.Writer to use for console output
//
// Returns:
//   - LoggerOption: A logger option function that can be passed to Init() or NewLogger()
func SetConsoleOutput(w io.Writer) LoggerOption {
	return func(ctx *loggerContext) {
		// Find and replace any existing console writers
		// Look for StyledConsoleWriter or default os.Stdout/os.Stderr writers
		hasConsoleWriter := false
		for i, output := range ctx.outputs {
			if _, ok := output.(*StyledConsoleWriter); ok {
				// Replace existing styled writer
				if ctx.useJSONFormat {
					ctx.outputs[i] = w
				} else {
					ctx.outputs[i] = NewStyledConsoleWriter(w, ctx)
				}
				hasConsoleWriter = true
				break
			}

			// Also check for direct stdout/stderr references
			if output == os.Stdout || output == os.Stderr {
				if ctx.useJSONFormat {
					ctx.outputs[i] = w
				} else {
					ctx.outputs[i] = NewStyledConsoleWriter(w, ctx)
				}
				hasConsoleWriter = true
				break
			}
		}

		// If no console writer found, add a new one
		if !hasConsoleWriter {
			if ctx.useJSONFormat {
				ctx.outputs = append(ctx.outputs, w)
			} else {
				ctx.outputs = append(ctx.outputs, NewStyledConsoleWriter(w, ctx))
			}
		}

		// Mark that we've handled console output
		ctx.consoleOn = false
	}
}
