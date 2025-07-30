// Package logger provides functionality for structured logging.
//
// This file contains the console writer implementation with support for colorized
// and styled output to make log messages more readable in terminal environments.
package logger

import (
	"context"
	"fmt"
	"io"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// StyledConsoleWriter is an io.Writer that formats log messages with styles and colors.
// It detects log levels and applies appropriate styling to make logs more readable.
type StyledConsoleWriter struct {
	out io.Writer
}

// NewStyledConsoleWriter creates a new StyledConsoleWriter instance.
//
// Parameters:
//   - w: The underlying io.Writer where formatted output will be written (typically os.Stdout)
//
// Returns:
//   - *StyledConsoleWriter: A new styled console writer that implements io.Writer
func NewStyledConsoleWriter(w io.Writer) *StyledConsoleWriter {
	return &StyledConsoleWriter{out: w}
}

// logLevelStyles defines the styling for each log level in console output.
// Each log level has a distinct color and formatting to make it easily identifiable
// in console output.
var logLevelStyles = map[string]lipgloss.Style{
	"TRACE": lipgloss.NewStyle().Foreground(lipgloss.Color("245")).Italic(true),
	"DEBUG": lipgloss.NewStyle().Foreground(lipgloss.Color("243")).Bold(true),
	"INFO":  lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true),
	"WARN":  lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Bold(true),
	"ERROR": lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true),
	"FATAL": lipgloss.NewStyle().Foreground(lipgloss.Color("160")).Bold(true).Underline(true),
}

// Write implements the io.Writer interface for StyledConsoleWriter.
// It formats the log message with appropriate styles and colors based on the log level.
//
// Parameters:
//   - p: The byte slice containing the log message to write
//
// Returns:
//   - int: The number of bytes written
//   - error: Any error encountered during writing
func (cw *StyledConsoleWriter) Write(p []byte) (int, error) {
	msg := string(p)
	level := detectLevel(msg)

	// If colors are disabled, use a simpler rendering
	if !colorEnabled {
		timestamp := time.Now().Format("15:04:05")
		line := fmt.Sprintf("[%s] %s", timestamp, strings.TrimSpace(msg))
		return fmt.Fprintln(cw.out, line)
	}

	// Color output below
	style, ok := logLevelStyles[level]
	if !ok {
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("7")).Bold(true)
	}

	timestamp := lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render(time.Now().Format("15:04:05"))

	styled := style.Render(strings.TrimSpace(msg))
	line := fmt.Sprintf("[%s] %s", timestamp, styled)
	return fmt.Fprintln(cw.out, line)
}

// detectLevel extracts the log level from a log message using the LEVEL=<LOG_LEVEL> pattern.
// It dynamically matches against the known levels from the logLevelStyles map.
//
// Parameters:
//   - s: The log message to analyze
//
// Returns:
//   - string: The detected log level, or empty string if not detected
func detectLevel(s string) string {
	s = strings.ToUpper(s)
	for level := range logLevelStyles {
		re := regexp.MustCompile(`\bLEVEL=` + level + `\b`)
		if re.MatchString(s) {
			return level
		}
	}
	return ""
}

// contextWithCaller creates a context with the caller's file and line number.
// This is useful for logging and debugging purposes, providing context about where the log message originated.
//
// Returns:
//   - context.Context: A context with the caller information stored as a value
func contextWithCaller() context.Context {
	_, file, line, _ := runtime.Caller(2)
	return context.WithValue(context.Background(), "caller", fmt.Sprintf("%s:%d", file, line))
}

// extractSource extracts source file information from a log message.
//
// Parameters:
//   - s: The log message to analyze
//
// Returns:
//   - string: The source file information, or empty string if not found
func extractSource(s string) string {
	re := regexp.MustCompile(`\bsource=([^ ]+)`)
	matches := re.FindStringSubmatch(s)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

// removeSourceFromMessage removes the source information from a log message.
// This helps prevent duplication when the source is already displayed separately.
//
// Parameters:
//   - msg: The log message to process
//
// Returns:
//   - string: The message with source information removed
func removeSourceFromMessage(msg string) string {
	// Remove the source=file.go:line pattern from the message
	re := regexp.MustCompile(`(source=[^ ]+)`)
	return re.ReplaceAllString(msg, "")
}
