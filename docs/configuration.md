# Configuration Guide

go-logo v2 provides flexible configuration through CLI flags, environment variables, and predefined profiles.

For complete API documentation, see the [API Reference](doc.md).

## Configuration Methods

### 1. Programmatic Configuration

```go
log := logo.New(
    logo.WithConsole(),
    logo.WithJSONFile("app.log"),
    logo.WithSource(true),
)
```

### 2. CLI Flags

```go
cfg := logo.NewConfig()
cfg.RegisterFlags(nil)  // Uses flag.CommandLine
flag.Parse()

logger, err := cfg.BuildLogger()
```

Run with:
```bash
./app -log.level=debug -log.console.format=pretty
```

### 3. Environment Variables

```go
cfg := logo.NewConfig()
cfg.LoadFromEnv()
flag.Parse()

logger, err := cfg.BuildLogger()
```

Set variables:
```bash
export LOG_LEVEL=debug
export LOG_CONSOLE_FORMAT=pretty
./app
```

### 4. Profiles

```bash
./app -log.profile=dev
./app -log.profile=prod
```

## Available Flags

| Flag | Environment | Default | Description |
|------|------------|---------|-------------|
| `-log.level` | `LOG_LEVEL` | `info` | Minimum log level (trace, debug, info, warn, error, fatal) |
| `-log.console` | `LOG_CONSOLE` | `true` | Enable console output |
| `-log.console.format` | `LOG_CONSOLE_FORMAT` | `pretty` | Console format (pretty, text) |
| `-log.console.color` | `LOG_CONSOLE_COLOR` | `auto` | Color mode (auto, always, never) |
| `-log.file.path` | `LOG_FILE_PATH` | `""` | Log file path (empty = disabled) |
| `-log.file.level` | `LOG_FILE_LEVEL` | `info` | Minimum file log level |
| `-log.file.format` | `LOG_FILE_FORMAT` | `json` | File format (json, text) |
| `-log.profile` | `LOG_PROFILE` | `""` | Profile (dev, prod) |
| `-log.source` | `LOG_SOURCE` | `false` | Include source location (file:line) |
| `-log.rotation.maxsize` | - | `100` | Max size in MB before rotation |
| `-log.rotation.maxage` | - | `30` | Max days to retain old files |
| `-log.rotation.maxbackups` | - | `10` | Max number of old files to keep |
| `-log.rotation.compress` | - | `true` | Compress rotated files with gzip |

## Profiles

### Development Profile (`dev`)

Optimized for developer experience:

```bash
./app -log.profile=dev
```

Configuration:
- **Level**: `debug`
- **Console**: Enabled with `pretty` format
- **Colors**: Auto-detected
- **Source**: Enabled (shows file:line)
- **File**: Disabled by default

Use when:
- Local development
- Debugging issues
- Running tests

### Production Profile (`prod`)

Optimized for production systems:

```bash
./app -log.profile=prod -log.file.path=/var/log/app.log
```

Configuration:
- **Level**: `info`
- **Console**: Enabled with `text` format (no colors)
- **Source**: Disabled
- **File**: JSON format with rotation (when path specified)
- **Rotation**: 100MB files, 30 day retention, compressed

Use when:
- Production deployments
- Docker containers
- CI/CD pipelines
- Log aggregation systems

## Configuration Priority

Configuration sources are applied in this order (later overrides earlier):

1. Default values
2. Environment variables
3. Profile application
4. CLI flags

Example:
```bash
# Profile sets level=info, but flag overrides to debug
export LOG_PROFILE=prod
./app -log.level=debug
```

## Validation

The configuration is validated before building a logger:

```go
logger, err := cfg.BuildLogger()
if err != nil {
    // Handle validation error
}
```

Common validation errors:
- Invalid log level: `"invalid log level: invalid (valid: trace, debug, info, warn, error, fatal)"`
- Invalid format: `"invalid console format: xml (valid: pretty, text)"`
- Invalid color mode: `"invalid color mode: blue (valid: auto, always, never)"`
- Invalid rotation: `"rotation max size must be at least 1 MB"`

## Examples

### Console Only (Default)

```bash
./app
```

Output to stdout with pretty formatting.

### File Only

```bash
./app -log.console=false -log.file.path=app.log
```

JSON logs to file only.

### Dual Output

```bash
./app -log.file.path=app.log
```

Pretty console + JSON file.

### Different Levels for Console and File

```bash
./app -log.level=debug -log.file.path=app.log -log.file.level=warn
```

Debug logs to console, only warnings and errors to file.

### Custom Rotation

```bash
./app -log.file.path=app.log \
  -log.rotation.maxsize=10 \
  -log.rotation.maxage=7 \
  -log.rotation.maxbackups=5
```

Rotate at 10MB, keep for 7 days, max 5 old files.

### Production with Environment Variables

```bash
export LOG_PROFILE=prod
export LOG_FILE_PATH=/var/log/myapp/app.log
export LOG_LEVEL=warn
./app
```

## Custom FlagSet

Use a custom flag set instead of the global one:

```go
fs := flag.NewFlagSet("myapp", flag.ExitOnError)
cfg := logo.NewConfig()
cfg.RegisterFlags(fs)
fs.Parse(os.Args[1:])

logger, _ := cfg.BuildLogger()
```

## Programmatic Configuration

For full control, build the logger manually:

```go
consoleHandler := logo.NewConsoleHandler(
    logo.WithConsoleLevel(logo.LevelInfo),
)

formatter := logo.NewJSONFormatter(logo.WithPrettyJSON(false))
fileHandler, _ := logo.NewFileHandler("app.log", formatter,
    logo.WithFileLevel(logo.LevelWarn),
    logo.WithRotation(&logo.RotationConfig{
        MaxSize:    50,
        MaxAge:     14,
        MaxBackups: 3,
        Compress:   true,
    }),
)

router := logo.NewRouter(consoleHandler, fileHandler)
logger := logo.New(
    logo.WithRouter(router),
    logo.WithSource(true),
    logo.WithName("myapp"),
)
```

## Best Practices

1. **Use profiles as a base** - Start with `dev` or `prod` and override specific settings
2. **Environment variables for secrets** - Don't put sensitive paths in code
3. **Validate early** - Call `cfg.Validate()` before building if you want custom error handling
4. **Log configuration on startup** - Log the active configuration for debugging
5. **Use JSON in production** - Structured logs are easier to parse and aggregate
6. **Enable rotation** - Prevent disk space issues with automatic rotation
7. **Set appropriate levels** - Use `debug` in dev, `info` or `warn` in production

