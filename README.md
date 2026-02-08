# go-logo

[![Go Reference](https://pkg.go.dev/badge/github.com/aN0mad/go-logo/v2.svg)](https://pkg.go.dev/github.com/aN0mad/go-logo/v2)
[![Go Report Card](https://goreportcard.com/badge/github.com/aN0mad/go-logo/v2)](https://goreportcard.com/report/github.com/aN0mad/go-logo/v2)

**go-logo** is a fast, zero-dependency, structured logging library for Go with beautiful console output and production-ready JSON logging.

## Features

- **Zero Dependencies** - No external packages required (except optional file rotation)
- **Beautiful Console Output** - Color-coded, human-readable logs with smart terminal detection
- **Structured Logging** - First-class support for structured fields with stable JSON schema
- **File Rotation** - Built-in log rotation with size and age limits
- **CLI Configuration** - Configure everything via flags, environment variables, or profiles
- **Context Integration** - Pass loggers through request contexts seamlessly
- **HTTP Middleware** - Ready-to-use examples for Fiber, Echo, Gin, and more
- **Thread-Safe** - Concurrent logging without data races
- **Type-Safe** - Compile-time safety with clear, idiomatic Go APIs

## Installation

```bash
go get github.com/aN0mad/go-logo/v2
```

> **Note:** v1 of go-logo has been moved to the `v1` branch. This is v2 which includes breaking changes and new features. See [migration guide](docs/migration-v1-to-v2.md) for upgrade instructions.

## Quick Start

### Basic Usage

```go
package main

import (
    "github.com/aN0mad/go-logo/v2"
)

func main() {
    // Create a logger with console output
    log := logo.New(logo.WithConsole())
    defer log.Close()
    
    log.Info("Hello, world!", logo.F("version", "2.0.0"))
    log.Debug("Debug message", logo.F("count", 42))
    log.Error("Something went wrong", logo.F("error", "timeout"))
}
```

### With CLI Configuration

```go
package main

import (
    "flag"
    "github.com/aN0mad/go-logo/v2"
)

func main() {
    // Configure via flags and environment variables
    cfg := logo.NewConfig()
    cfg.RegisterFlags(nil)
    flag.Parse()
    
    logger, err := cfg.BuildLogger()
    if err != nil {
        panic(err)
    }
    defer logger.Close()
    
    logger.Info("Application started")
}
```

Run with different configurations:
```bash
# Development mode
./app -log.level=debug -log.console.format=pretty

# Production mode  
./app -log.profile=prod -log.file.path=app.log

# Via environment variables
LOG_LEVEL=debug LOG_CONSOLE_COLOR=always ./app
```

### With File Logging

```go
log := logo.New(
    logo.WithConsole(),
    logo.WithJSONFile("app.log"),
)
defer log.Close()

log.Info("Logged to both console and file")
```

### Context-Aware Logging

```go
ctx := logo.WithLogger(context.Background(), log)
ctx = logo.WithRequestID(ctx, "req-123")

// Later in your code
logger := logo.FromContext(ctx)
logger.Info("Processing request") // Includes request ID
```

## Documentation

- **[API Reference](docs/doc.md)** - Complete API documentation
- **[Configuration Guide](docs/configuration.md)** - CLI flags, environment variables, and profiles
- **[HTTP Integration](docs/http.md)** - Middleware examples for web frameworks
- **[Architecture Overview](docs/architecture.md)** - Design principles and patterns
- **[Schema Reference](docs/schema.md)** - JSON logging schema specification
- **[File Rotation](docs/rotation.md)** - Log rotation configuration
- **[Migration from v1](docs/migration-v1-to-v2.md)** - Upgrade guide
- **[Changelog](CHANGELOG.md)** - Version history and release notes

## Examples

Explore the [examples/](examples/) directory for complete, runnable examples:

### Core Examples
- **[basic](examples/basic/)** - Simple console logging
- **[console_pretty](examples/console_pretty/)** - Pretty-formatted console output
- **[console_colors](examples/console_colors/)** - Color mode demonstrations
- **[json_logging](examples/json_logging/)** - Structured JSON logging
- **[file_logging](examples/file_logging/)** - File output with rotation
- **[dual_output](examples/dual_output/)** - Console + file simultaneously
- **[multi_handler](examples/multi_handler/)** - Custom routing patterns

### CLI Examples
- **[cli_basic](examples/cli_basic/)** - CLI flag usage
- **[cli_profiles](examples/cli_profiles/)** - Development and production profiles
- **[cli_full](examples/cli_full/)** - Comprehensive CLI configuration

### HTTP Framework Examples
- **[http_fiber_v2](examples/http_fiber_v2/)** - Fiber v2 middleware
- **[http_fiber_v3](examples/http_fiber_v3/)** - Fiber v3 middleware (beta)
- **[http_echo](examples/http_echo/)** - Echo framework middleware
- **[http_gin](examples/http_gin/)** - Gin framework middleware

See [HTTP Examples documentation](docs/HTTP_EXAMPLES.md) for detailed HTTP integration patterns.

## Common Patterns

### Structured Fields

```go
log.Info("User login",
    logo.F("user_id", 12345),
    logo.F("ip", "192.168.1.1"),
    logo.F("duration_ms", 42),
)
```

### Child Loggers

```go
baseLogger := logo.New(logo.WithConsole())
requestLogger := baseLogger.With(
    logo.F("request_id", "abc-123"),
    logo.F("user", "john"),
)

requestLogger.Info("Processing") // Includes request_id and user
```

### Multiple Handlers

```go
consoleHandler := logo.NewConsoleHandler(
    logo.WithConsoleLevel(logo.LevelInfo),
)

fileHandler, _ := logo.NewFileHandler("app.log", 
    logo.NewJSONFormatter(),
    logo.WithFileLevel(logo.LevelDebug),
)

router := logo.NewRouter(consoleHandler, fileHandler)
log := logo.New(logo.WithRouter(router))

log.Debug("Detailed info")  // Only to file
log.Info("General info")    // To both
```

## Configuration Profiles

### Development Profile
```bash
./app -log.profile=dev
```
- Debug level logging
- Pretty console format with colors  
- Source location tracking enabled
- Optimized for developer experience

### Production Profile
```bash
./app -log.profile=prod -log.file.path=app.log
```
- Info level logging
- JSON file output with rotation
- Source location disabled
- Optimized for performance

## Contributing

Contributions are welcome! Please feel free to submit issues and pull requests.

## License

MIT License - see [LICENSE](LICENSE) for details.

## Version History

- **v2.0.0** - Major rewrite with zero dependencies, improved API, and production features
- **v1.x** - Legacy version (available on `v1` branch)
