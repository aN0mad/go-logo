package logo

import (
	"io"
	"os"
	"runtime"
)

// ANSI color codes
const (
	colorReset   = "\033[0m"
	colorRed     = "\033[31m"
	colorYellow  = "\033[33m"
	colorBlue    = "\033[34m"
	colorGray    = "\033[90m"
	colorCyan    = "\033[36m"
	colorGreen   = "\033[32m"
	colorMagenta = "\033[35m"

	colorBoldRed     = "\033[1;31m"
	colorBoldYellow  = "\033[1;33m"
	colorBoldBlue    = "\033[1;34m"
	colorBoldGray    = "\033[1;90m"
	colorBoldMagenta = "\033[1;35m"
	colorBoldCyan    = "\033[1;36m"
)

// ColorMode determines how colors are applied
type ColorMode int

const (
	// ColorAuto automatically detects terminal color support
	ColorAuto ColorMode = iota
	// ColorAlways always uses colors
	ColorAlways
	// ColorNever never uses colors
	ColorNever
)

// ColorSupport determines if the writer supports colors
type ColorSupport struct {
	enabled bool
}

// NewColorSupport creates a new color support detector
func NewColorSupport(writer io.Writer, mode ColorMode) *ColorSupport {
	switch mode {
	case ColorAlways:
		return &ColorSupport{enabled: true}
	case ColorNever:
		return &ColorSupport{enabled: false}
	case ColorAuto:
		return &ColorSupport{enabled: detectColorSupport(writer)}
	default:
		return &ColorSupport{enabled: false}
	}
}

// Enabled returns whether colors are enabled
func (c *ColorSupport) Enabled() bool {
	return c.enabled
}

// Colorize wraps text in ANSI color codes if colors are enabled
func (c *ColorSupport) Colorize(text, color string) string {
	if !c.enabled {
		return text
	}
	return color + text + colorReset
}

// detectColorSupport attempts to detect if the writer supports colors
func detectColorSupport(writer io.Writer) bool {
	// Check if writer is a file
	file, ok := writer.(*os.File)
	if !ok {
		// Not a file, assume no color support
		return false
	}

	// Check if it's stdout or stderr
	if file != os.Stdout && file != os.Stderr {
		return false
	}

	// Check for NO_COLOR environment variable
	if os.Getenv("NO_COLOR") != "" {
		return false
	}

	// Check for TERM environment variable
	term := os.Getenv("TERM")
	if term == "" || term == "dumb" {
		return false
	}

	// Check for CI environments that typically don't support colors well
	if os.Getenv("CI") != "" {
		// Some CI systems do support colors
		// Check for specific CI systems that support colors
		if os.Getenv("GITHUB_ACTIONS") != "" ||
			os.Getenv("GITLAB_CI") != "" {
			return true
		}
		// Default to no colors in CI
		return false
	}

	// Check if running on Windows
	if runtime.GOOS == "windows" {
		// Modern Windows Terminal supports ANSI colors
		// Check for Windows Terminal
		if os.Getenv("WT_SESSION") != "" {
			return true
		}
		// Otherwise, be conservative on Windows
		return false
	}

	// On Unix-like systems with a proper TERM, assume color support
	return true
}

// LevelColor returns the color for a log level
func LevelColor(level Level) string {
	switch level {
	case LevelTrace:
		return colorCyan
	case LevelDebug:
		return colorGray
	case LevelInfo:
		return colorBlue
	case LevelWarn:
		return colorYellow
	case LevelError:
		return colorRed
	case LevelFatal:
		return colorMagenta
	default:
		return colorReset
	}
}

// LevelColorBold returns the bold color for a log level
func LevelColorBold(level Level) string {
	switch level {
	case LevelTrace:
		return colorBoldCyan
	case LevelDebug:
		return colorBoldGray
	case LevelInfo:
		return colorBoldBlue
	case LevelWarn:
		return colorBoldYellow
	case LevelError:
		return colorBoldRed
	case LevelFatal:
		return colorBoldMagenta
	default:
		return colorReset
	}
}
