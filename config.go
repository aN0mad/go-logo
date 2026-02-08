package logo

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

// Config holds the logging configuration that can be set via CLI flags or environment variables.
type Config struct {
	// Level is the minimum log level (debug, info, warn, error)
	Level string

	// Console enables console output
	Console bool

	// ConsoleFormat is the console output format (pretty, text)
	ConsoleFormat string

	// ConsoleColor controls color output (auto, always, never)
	ConsoleColor string

	// FilePath is the path to the log file
	FilePath string

	// FileLevel is the minimum level for file output
	FileLevel string

	// FileFormat is the file output format (json, text)
	FileFormat string

	// Profile is a predefined configuration profile (dev, prod)
	Profile string

	// Source enables source location in logs
	Source bool

	// Rotation settings for file output
	RotationMaxSize    int
	RotationMaxAge     int
	RotationMaxBackups int
	RotationCompress   bool
}

// NewConfig creates a new configuration with defaults.
func NewConfig() *Config {
	return &Config{
		Level:              "info",
		Console:            true,
		ConsoleFormat:      "pretty",
		ConsoleColor:       "auto",
		FilePath:           "",
		FileLevel:          "info",
		FileFormat:         "json",
		Profile:            "",
		Source:             false,
		RotationMaxSize:    100,
		RotationMaxAge:     30,
		RotationMaxBackups: 10,
		RotationCompress:   true,
	}
}

// RegisterFlags registers CLI flags for the configuration.
// Call this before flag.Parse().
func (c *Config) RegisterFlags(fs *flag.FlagSet) {
	if fs == nil {
		fs = flag.CommandLine
	}

	fs.StringVar(&c.Level, "log.level", c.Level, "Minimum log level (debug, info, warn, error)")
	fs.BoolVar(&c.Console, "log.console", c.Console, "Enable console output")
	fs.StringVar(&c.ConsoleFormat, "log.console.format", c.ConsoleFormat, "Console format (pretty, text)")
	fs.StringVar(&c.ConsoleColor, "log.console.color", c.ConsoleColor, "Console color (auto, always, never)")
	fs.StringVar(&c.FilePath, "log.file.path", c.FilePath, "Log file path (empty disables file logging)")
	fs.StringVar(&c.FileLevel, "log.file.level", c.FileLevel, "Minimum level for file output")
	fs.StringVar(&c.FileFormat, "log.file.format", c.FileFormat, "File format (json, text)")
	fs.StringVar(&c.Profile, "log.profile", c.Profile, "Configuration profile (dev, prod)")
	fs.BoolVar(&c.Source, "log.source", c.Source, "Include source location in logs")
	fs.IntVar(&c.RotationMaxSize, "log.rotation.maxsize", c.RotationMaxSize, "Max size in MB before rotation")
	fs.IntVar(&c.RotationMaxAge, "log.rotation.maxage", c.RotationMaxAge, "Max days to retain old files")
	fs.IntVar(&c.RotationMaxBackups, "log.rotation.maxbackups", c.RotationMaxBackups, "Max number of old files to keep")
	fs.BoolVar(&c.RotationCompress, "log.rotation.compress", c.RotationCompress, "Compress rotated files")
}

// LoadFromEnv loads configuration from environment variables.
// Environment variables follow the pattern: LOG_<FLAG_NAME> (e.g., LOG_LEVEL, LOG_CONSOLE_FORMAT)
func (c *Config) LoadFromEnv() {
	if val := os.Getenv("LOG_LEVEL"); val != "" {
		c.Level = val
	}
	if val := os.Getenv("LOG_CONSOLE"); val != "" {
		c.Console = val == "true" || val == "1"
	}
	if val := os.Getenv("LOG_CONSOLE_FORMAT"); val != "" {
		c.ConsoleFormat = val
	}
	if val := os.Getenv("LOG_CONSOLE_COLOR"); val != "" {
		c.ConsoleColor = val
	}
	if val := os.Getenv("LOG_FILE_PATH"); val != "" {
		c.FilePath = val
	}
	if val := os.Getenv("LOG_FILE_LEVEL"); val != "" {
		c.FileLevel = val
	}
	if val := os.Getenv("LOG_FILE_FORMAT"); val != "" {
		c.FileFormat = val
	}
	if val := os.Getenv("LOG_PROFILE"); val != "" {
		c.Profile = val
	}
	if val := os.Getenv("LOG_SOURCE"); val != "" {
		c.Source = val == "true" || val == "1"
	}
}

// ApplyProfile applies a predefined configuration profile.
// Profiles are applied first, then individual flags override.
func (c *Config) ApplyProfile(profile string) error {
	switch strings.ToLower(profile) {
	case "dev", "development":
		c.Level = "trace"
		c.Console = true
		c.ConsoleFormat = "pretty"
		c.ConsoleColor = "auto"
		c.Source = true
		c.FileFormat = "json"
		return nil

	case "prod", "production":
		c.Level = "info"
		c.Console = true
		c.ConsoleFormat = "text"
		c.ConsoleColor = "never"
		c.Source = false
		c.FileFormat = "json"
		c.RotationCompress = true
		return nil

	case "":
		// No profile, use defaults or explicit flags
		return nil

	default:
		return fmt.Errorf("unknown profile: %s (valid: dev, prod)", profile)
	}
}

// Validate checks the configuration for invalid values.
func (c *Config) Validate() error {
	// Validate level
	level := strings.ToLower(c.Level)
	if level != "trace" && level != "debug" && level != "info" && level != "warn" && level != "error" && level != "fatal" {
		return fmt.Errorf("invalid log level: %s (valid: trace, debug, info, warn, error, fatal)", c.Level)
	}

	// Validate file level
	fileLevel := strings.ToLower(c.FileLevel)
	if fileLevel != "trace" && fileLevel != "debug" && fileLevel != "info" && fileLevel != "warn" && fileLevel != "error" && fileLevel != "fatal" {
		return fmt.Errorf("invalid file level: %s (valid: trace, debug, info, warn, error, fatal)", c.FileLevel)
	}

	// Validate console format
	consoleFormat := strings.ToLower(c.ConsoleFormat)
	if consoleFormat != "pretty" && consoleFormat != "text" {
		return fmt.Errorf("invalid console format: %s (valid: pretty, text)", c.ConsoleFormat)
	}

	// Validate console color
	consoleColor := strings.ToLower(c.ConsoleColor)
	if consoleColor != "auto" && consoleColor != "always" && consoleColor != "never" {
		return fmt.Errorf("invalid console color: %s (valid: auto, always, never)", c.ConsoleColor)
	}

	// Validate file format
	fileFormat := strings.ToLower(c.FileFormat)
	if fileFormat != "json" && fileFormat != "text" {
		return fmt.Errorf("invalid file format: %s (valid: json, text)", c.FileFormat)
	}

	// Validate rotation settings
	if c.RotationMaxSize < 1 {
		return fmt.Errorf("rotation max size must be at least 1 MB")
	}

	if c.RotationMaxAge < 0 {
		return fmt.Errorf("rotation max age cannot be negative")
	}

	if c.RotationMaxBackups < 0 {
		return fmt.Errorf("rotation max backups cannot be negative")
	}

	return nil
}

// BuildLogger creates a logger from the configuration.
func (c *Config) BuildLogger() (*Logger, error) {
	// Apply profile first
	if c.Profile != "" {
		if err := c.ApplyProfile(c.Profile); err != nil {
			return nil, err
		}
	}

	// Validate configuration
	if err := c.Validate(); err != nil {
		return nil, err
	}

	var handlers []Handler

	// Parse log level
	level := ParseLevel(c.Level)

	// Add console handler if enabled
	if c.Console {
		var colorMode ColorMode
		switch strings.ToLower(c.ConsoleColor) {
		case "always":
			colorMode = ColorAlways
		case "never":
			colorMode = ColorNever
		default:
			colorMode = ColorAuto
		}

		var formatter Formatter
		if strings.ToLower(c.ConsoleFormat) == "pretty" {
			colorSupport := NewColorSupport(os.Stdout, colorMode)
			opts := []PrettyFormatterOption{}
			if c.Source {
				opts = append(opts, WithSourceInfo(true))
			}
			formatter = NewPrettyFormatter(colorSupport, opts...)
		} else {
			formatter = NewTextFormatter()
		}

		handler := NewWriterHandler(os.Stdout, formatter, level)
		handlers = append(handlers, handler)
	}

	// Add file handler if path is specified
	if c.FilePath != "" {
		fileLevel := ParseLevel(c.FileLevel)

		var formatter Formatter
		if strings.ToLower(c.FileFormat) == "json" {
			opts := []JSONFormatterOption{}
			if c.Source {
				opts = append(opts, WithJSONSource(true))
			}
			formatter = NewJSONFormatter(opts...)
		} else {
			formatter = NewTextFormatter()
		}

		rotationConfig := &RotationConfig{
			MaxSize:    c.RotationMaxSize,
			MaxAge:     c.RotationMaxAge,
			MaxBackups: c.RotationMaxBackups,
			Compress:   c.RotationCompress,
		}

		fileHandler, err := NewFileHandler(c.FilePath, formatter,
			WithRotation(rotationConfig),
			WithFileLevel(fileLevel),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create file handler: %w", err)
		}

		handlers = append(handlers, fileHandler)
	}

	// Create router with handlers
	router := NewRouter(handlers...)

	// Create logger
	opts := []Option{
		WithRouter(router),
	}
	if c.Source {
		opts = append(opts, WithSource(true))
	}

	return New(opts...), nil
}
