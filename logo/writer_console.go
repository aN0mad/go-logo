// Package logo provides functionality for structured logging.
//
// This file contains the console writer implementation with support for colorized
// and styled output to make log messages more readable in terminal environments.
package logo

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
	ctx *loggerContext // Reference to the logger context for configuration
}

// NewStyledConsoleWriter creates a new console writer that applies styling to log output.
// It applies colors and formatting based on log level and configuration.
//
// Parameters:
//   - out: The output writer where log messages will be written
//   - ctx: Logger context containing configuration options like color settings
//
// Returns:
//   - *StyledConsoleWriter: A new console writer with styling capabilities
func NewStyledConsoleWriter(w io.Writer, ctx *loggerContext) *StyledConsoleWriter {
	return &StyledConsoleWriter{
		out: w,
		ctx: ctx,
	}
}

// NewDefaultStyledConsoleWriter creates a styled console writer with default settings.
// It creates a writer that uses global configuration rather than a specific context.
//
// Parameters:
//   - out: The output writer where log messages will be written
//
// Returns:
//   - *StyledConsoleWriter: A new console writer with default styling configuration
func NewDefaultStyledConsoleWriter(w io.Writer) *StyledConsoleWriter {
	return &StyledConsoleWriter{
		out: w,
		ctx: nil, // Will use global defaults
	}
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

	// Determine if colors should be enabled
	colorEnabledForThisWriter := COLORENABLED
	if cw.ctx != nil {
		colorEnabledForThisWriter = cw.ctx.colorEnabled
	}

	// If colors are disabled, use a simpler rendering
	if !colorEnabledForThisWriter {
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

// detectLevel extracts the log level from a log message.
// It parses the message string to find the level indicator.
//
// Parameters:
//   - message: The log message to parse
//
// Returns:
//   - string: The detected log level, or empty string if none found
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

// contextWithCaller creates a context that includes caller information.
// It captures the current stack frame information for source code location.
//
// Returns:
//   - context.Context: A context containing caller information
func contextWithCaller() context.Context {
	_, file, line, _ := runtime.Caller(2)
	return context.WithValue(context.Background(), "caller", fmt.Sprintf("%s:%d", file, line))
}

// extractSource extracts source file and line information from a log message.
// It parses the message string to find source code location information.
//
// Parameters:
//   - message: The log message to parse
//
// Returns:
//   - string: The extracted source information, or empty string if none found
func extractSource(s string) string {
	re := regexp.MustCompile(`\bsource=([^ ]+)`)
	matches := re.FindStringSubmatch(s)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

// removeSourceFromMessage removes source information from a log message.
// This is useful when the source should be displayed separately from the message.
//
// Parameters:
//   - message: The log message containing source information
//
// Returns:
//   - string: The message with source information removed
func removeSourceFromMessage(msg string) string {
	// Remove the source=file.go:line pattern from the message
	re := regexp.MustCompile(`(source=[^ ]+)`)
	return re.ReplaceAllString(msg, "")
}
