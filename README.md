# go-logo
A golang logging library that allows logging to the console, to a channel or both.


## Features
- Multiple log levels (TRACE, DEBUG, INFO, WARN, ERROR, FATAL)
- Multiple output formats (text, JSON, pretty JSON)
- Multiple output destinations (console, file, channel)
- Colorized console output
- Structured logging with attributes
- Source code location information
- File rotation with size and age limits
- Context-aware logging
- Channel-based logging for asynchronous processing

## Usage
```bash
/workspace/bin ./go-logo.elf
Initialized logger
[18:56:51] time=2025-07-30T18:56:51.722Z level=DEBUG msg=This is a debug message source=/workspace/cmd/main.go:23
[18:56:51] time=2025-07-30T18:56:51.722Z level=INFO msg=This is an info message source=/workspace/cmd/main.go:24
[18:56:51] time=2025-07-30T18:56:51.723Z level=WARN msg=This is a warning message source=/workspace/cmd/main.go:25
[18:56:51] time=2025-07-30T18:56:51.723Z level=ERROR msg=This is an error message source=/workspace/cmd/main.go:26
[18:56:51] time=2025-07-30T18:56:51.723Z level=FATAL msg=This is a fatal message - It will exit the program source=/workspace/cmd/main.go:27
```

### Colored output
![Colorful logs](assets/colored_logs.png)

## Example
Example source code is given in the [examples](examples/) directory.

### Logging to console
```golang
package main

import (
	"fmt"
	"log/slog"
	logger "github.com/aN0mad/go-logo/logo"
)

func main() {

    // Initialize the logger with desired options
	logger.Init(
		logger.AddSource(),                                    // Log source caller information
        logger.SetLevel(slog.LevelDebug),                      // Set minimum log level
        logger.AddFileOutput("logs/app.log", 10, 3, 30, true), // Log to file with rotation
	)

    defer logger.Close() // Ensure all file writers are closed properly

    // Get the logger instance
	log := logger.L()
	fmt.Println("Initialized logger")

    // Log messages at different levels
	log.Trace("This is a trace message - It will not show up")
	log.Debug("This is a debug message")
	log.Info("This is an info message")
	log.Warn("This is a warning message")
	log.Error("This is an error message")
	log.Fatal("This is a fatal message - It will exit the program")
}
```

## Configuration Options
### Log levels
```golang
logger.Init(
    logger.SetLevel(slog.LevelDebug), // Standard log levels: DEBUG, INFO, WARN, ERROR
    logger.EnableTrace(),             // Enable TRACE level (below DEBUG)
)
```

### Output Formats
```golang
// Text format (default)
logger.Init()

// JSON format
logger.Init(
    logger.UseJSON(false) // Compact JSON
)

// Pretty JSON format
logger.Init(
    logger.UseJSON(true) // Pretty-printed JSON
)
```

### Output Destinations
```golang
// Console output (default)
logger.Init()

// File output with rotation
logger.Init(
    logger.AddFileOutput("logs/app.log", 10, 3, 30, true)
    // Parameters: path, maxSize (MB), backups, maxAge (days), compress
)

// Channel output
logger.Init(
    logger.AddChannelOutput(logChan) // logChan is a chan string
)

// Multiple outputs
logger.Init(
    logger.AddFileOutput("logs/app.log", 10, 3, 30, true),
    logger.AddChannelOutput(logChan)
)

// Disable console output when using other outputs
logger.Init(
    logger.DisableConsole(),
    logger.AddFileOutput("logs/app.log", 10, 3, 30, true)
)
```

### Additional Features
```golang
// Add source file and line information
logger.Init(
    logger.AddSource()
)

// Disable colored output
logger.Init(
    logger.DisableColors()
)

// Use custom handler
logger.Init(
    logger.UseCustomHandler(myCustomHandler)
)

// Context-aware logging
ctx := context.WithValue(context.Background(), "request_id", "req-123")
requestLogger := logger.WithContext(ctx)
requestLogger.Info("Processing request") // Includes request_id automatically
```

## Building
The included Makefile makes building and managing the project easy:

```bash
# List available commands
make help

# Format code and tidy dependencies
make tidy

# Build for current platform
make build

# Build for Windows
make buildwin

# Build for both platforms
make all

# Clean build artifacts
make clean
```

## Examples
Check the examples directory for more detailed examples:

- Text logging
- JSON logging
- Pretty JSON logging
- File logging without console
- Channel-based logging
- Context-aware logging
- Comprehensive logger configuration

### Logging with context
```golang
// No WithContext in core library
// To use with web contexts, extend the logger in your application:

// In your web app's logger package
func LoggerFromContext(ctx context.Context) *logo.Logger {
    logger := logo.L()
    
    // Add application-specific context values
    if requestID, ok := ctx.Value("request_id").(string); ok {
        logger = logger.With("request_id", requestID)
    }
    if userID, ok := ctx.Value("user_id").(string); ok {
        logger = logger.With("user_id", userID)
    }
    
    return logger
}
```

## Running Tests
### Run all unit tests:
```bash
go test ./logo/
```

### Check test coverage
```bash
go test -v -coverprofile coverage.out
```

#### Generate HTML report
```bash
go tool cover -html coverage.out -o coverage.html
```

## Enhancements
- Migrate from [Lumberjack](https://github.com/aN0mad/lumberjack) to [timberjack](https://github.com/DeRuina/timberjack/) for enhanced log file control

## Acknowledgments
- [natefinch](https://github.com/natefinch) for the [Lumberjack](https://github.com/natefinch/lumberjack) package to handle log files and rotations
- [charmbracelet](https://github.com/charmbracelet) for the [lipgloss](https://github.com/charmbracelet/lipgloss) package for ANSI coloring
- Docs generated with [gomarkdoc](https://github.com/princjef/gomarkdoc) using `/go/bin/gomarkdoc . > doc.md` within the module directory (`./logo`)
