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

// StyledConsoleWriter is a custom writer that formats log messages with colors and styles.
// It implements the io.Writer interface, allowing it to be used with standard logging functions.
// This is useful for applications that require visually distinct log messages in the console.
type StyledConsoleWriter struct {
	out io.Writer
}

// NewStyledConsoleWriter creates a new StyledConsoleWriter instance.
//
// Parameters:
//   - w: The underlying io.Writer where formatted output will be written (typically os.Stdout)
//
// Returns a StyledConsoleWriter that implements io.Writer
func NewStyledConsoleWriter(w io.Writer) *StyledConsoleWriter {
	return &StyledConsoleWriter{out: w}
}

// logLevelStyles maps log levels to their corresponding lipgloss styles.
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
// Returns the number of bytes written and any error encountered
func (cw *StyledConsoleWriter) Write(p []byte) (int, error) {
	msg := string(p)
	level := detectLevel(msg)

	// Remove the source information from the message to avoid duplication
	// This will remove "source=file.go:line" from the displayed message
	cleanedMsg := removeSourceFromMessage(msg)

	// If colors are disabled, use a simpler rendering
	if !colorEnabled {
		timestamp := time.Now().Format("15:04:05")

		// Extract source information if present
		sourceInfo := extractSource(cleanedMsg)
		sourceText := ""
		if sourceInfo != "" {
			sourceText = " " + sourceInfo
		}

		line := fmt.Sprintf("[%s]%s %s", timestamp, sourceText, strings.TrimSpace(msg))
		return fmt.Fprintln(cw.out, line)
	}

	// Color output below (existing code)
	style, ok := logLevelStyles[level]
	if !ok {
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("7")).Bold(true)
	}

	timestamp := lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render(time.Now().Format("15:04:05"))

	// Extract source information if present
	sourceInfo := extractSource(msg)
	sourceText := ""
	if sourceInfo != "" {
		sourceText = lipgloss.NewStyle().Foreground(lipgloss.Color("99")).Render(sourceInfo)
		sourceText = " " + sourceText
	}

	styled := style.Render(strings.TrimSpace(msg))
	line := fmt.Sprintf("[%s]%s %s", timestamp, sourceText, styled)
	return fmt.Fprintln(cw.out, line)
}

// detectLevel extracts the log level from a log message using the LEVEL=<LOG_LEVEL> pattern.
// It dynamically matches against the known levels from the logLevelStyles map.
//
// Parameters:
//   - s: The log message to analyze
//
// Returns the detected log level as a string, or empty string if not detected
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
// Returns a context.Context with the caller information stored as a value
func contextWithCaller() context.Context {
	_, file, line, _ := runtime.Caller(2)
	return context.WithValue(context.Background(), "caller", fmt.Sprintf("%s:%d", file, line))
}

// extractSource extracts source file information from a log message.
//
// Parameters:
//   - s: The log message to analyze
//
// Returns the source file information as a string, or empty string if not found
func extractSource(s string) string {
	re := regexp.MustCompile(`\bsource=([^ ]+)`)
	matches := re.FindStringSubmatch(s)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

// removeSourceFromMessage removes the source=file.go:line from the message
//
// Parameters:
//   - msg: The log message to clean
//
// Returns the log message with source information removed
func removeSourceFromMessage(msg string) string {
	// Remove the source=file.go:line pattern from the message
	re := regexp.MustCompile(`source=[^ ]+\s*`)
	return re.ReplaceAllString(msg, "")
}
