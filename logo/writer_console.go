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
func NewStyledConsoleWriter(w io.Writer) *StyledConsoleWriter {
	return &StyledConsoleWriter{out: w}
}

// logLevelStyles maps log levels to their corresponding lipgloss styles.
var logLevelStyles = map[string]lipgloss.Style{
	"TRACE": lipgloss.NewStyle().Foreground(lipgloss.Color("245")).Italic(true),
	"DEBUG": lipgloss.NewStyle().Foreground(lipgloss.Color("243")).Bold(true),
	"INFO":  lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true),
	"WARN":  lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Bold(true),
	"ERROR": lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true),
	"FATAL": lipgloss.NewStyle().Foreground(lipgloss.Color("160")).Bold(true).Underline(true),
}

// Write implements the io.Writer interface for StyledConsoleWriter.
func (cw *StyledConsoleWriter) Write(p []byte) (int, error) {
	msg := string(p)
	level := detectLevel(msg)

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
	line := fmt.Sprintf("[%s] %s", timestamp, styled)
	return fmt.Fprintln(cw.out, line)
}

// detectLevel extracts the log level from a log message using the LEVEL=<LOG_LEVEL> pattern.
// It dynamically matches against the known levels from the logLevelStyles map.
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
func contextWithCaller() context.Context {
	_, file, line, _ := runtime.Caller(2)
	return context.WithValue(context.Background(), "caller", fmt.Sprintf("%s:%d", file, line))
}

// Add this function to extract source info from log message
func extractSource(s string) string {
	re := regexp.MustCompile(`\bsource=([^ ]+)`)
	matches := re.FindStringSubmatch(s)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}
