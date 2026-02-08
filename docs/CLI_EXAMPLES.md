# CLI Configuration Examples

This directory contains examples demonstrating Phase 5 - CLI Configuration & Profiles.

## Examples

### 1. cli_basic - Basic CLI Flags

Demonstrates basic CLI flag parsing for logger configuration.

```bash
# Run with default settings (info level, pretty format)
go run main.go

# Run with debug level
go run main.go -log.level=debug

# Run with text format instead of pretty
go run main.go -log.console.format=text

# Disable console colors
go run main.go -log.console.color=never

# Enable source location tracking
go run main.go -log.source=true
```

Available flags:
- `-log.level` - Minimum log level (debug, info, warn, error)
- `-log.console` - Enable console output (default: true)
- `-log.console.format` - Console format (pretty, text)
- `-log.console.color` - Color mode (auto, always, never)
- `-log.file.path` - Log file path
- `-log.file.level` - File log level
- `-log.file.format` - File format (json, text)
- `-log.source` - Include source location
- `-log.rotation.*` - File rotation settings

### 2. cli_profiles - Configuration Profiles

Demonstrates predefined profiles (dev/prod) and environment variables.

```bash
# Run in development mode (debug level, colors, source info)
go run main.go -log.profile=dev

# Run in production mode (info level, no colors, no source)
go run main.go -log.profile=prod

# Use environment variable
LOG_PROFILE=dev go run main.go

# Override profile with explicit flags
go run main.go -log.profile=prod -log.level=debug
```

Profile configurations:
- **dev** - Debug level, pretty format, colors auto, source enabled
- **prod** - Info level, text format, no colors, source disabled

### 3. cli_full - Complete Example

Comprehensive example showing all features working together.

```bash
# Production mode with structured logging
go run main.go -log.profile=prod

# Development mode with debug output
go run main.go -log.profile=dev

# Custom configuration
go run main.go \
  -log.level=debug \
  -log.console.format=text \
  -log.console.color=always \
  -log.source=true

# Log to file with rotation
go run main.go \
  -log.file.path=/tmp/app.log \
  -log.file.format=json \
  -log.rotation.maxsize=50 \
  -log.rotation.maxage=7 \
  -log.rotation.compress=true
```

## Environment Variables

All flags can be configured via environment variables using the `LOG_*` prefix:

```bash
export LOG_LEVEL=debug
export LOG_CONSOLE_FORMAT=text
export LOG_CONSOLE_COLOR=never
export LOG_SOURCE=true
export LOG_PROFILE=prod

go run main.go
```

Environment variable mapping:
- `LOG_LEVEL` → `-log.level`
- `LOG_CONSOLE` → `-log.console`
- `LOG_CONSOLE_FORMAT` → `-log.console.format`
- `LOG_CONSOLE_COLOR` → `-log.console.color`
- `LOG_FILE_PATH` → `-log.file.path`
- `LOG_FILE_LEVEL` → `-log.file.level`
- `LOG_FILE_FORMAT` → `-log.file.format`
- `LOG_PROFILE` → `-log.profile`
- `LOG_SOURCE` → `-log.source`

## Configuration Priority

Configuration sources are applied in this order (later overrides earlier):

1. Default values
2. Environment variables
3. CLI flags
4. Profile (if specified)

Note: When a profile is set, it's applied first, but explicit CLI flags will override individual profile settings.

## Integration Pattern

```go
package main

import (
    "flag"
    "fmt"
    logo "github.com/aN0mad/go-logo/v2"
)

func main() {
    // 1. Create config with defaults
    cfg := logo.NewConfig()
    
    // 2. Load environment variables
    cfg.LoadFromEnv()
    
    // 3. Register CLI flags
    cfg.RegisterFlags(nil)
    
    // 4. Parse command line
    flag.Parse()
    
    // 5. Build logger (validates config)
    logger, err := cfg.BuildLogger()
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    defer logger.Close()
    
    // 6. Use logger
    logger.Info("Application started")
}
```

## Validation

The config system validates all settings before building the logger:

- Log levels must be: debug, info, warn, error
- Console format must be: pretty, text
- Console color must be: auto, always, never
- File format must be: json, text
- Rotation settings must be valid positive numbers

Invalid configurations will return a clear error message.
