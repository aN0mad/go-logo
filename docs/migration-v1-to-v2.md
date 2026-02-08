# Migration Guide: v1 → v2

This guide helps you migrate from go-logo v1 to v2. The v2 release brings significant architectural improvements, better performance, and more flexibility while maintaining the core logging functionality you rely on.

## Table of Contents
- [Quick Start](#quick-start)
- [Import Path Changes](#import-path-changes)
- [Breaking Changes](#breaking-changes)
- [API Migration Table](#api-migration-table)
- [Common Migration Patterns](#common-migration-patterns)
- [Feature Comparison](#feature-comparison)
- [Upgrade Checklist](#upgrade-checklist)

---

## Quick Start

The fastest way to migrate is to understand the core architectural change: v1 used a global singleton logger with `Init()` and `L()`, while v2 uses explicit logger instances with a builder pattern.

### v1 Pattern
```go
logo.Init(
    logo.SetLevel(slog.LevelInfo),
    logo.AddSource(),
)
log := logo.L()
log.Info("hello")
```

### v2 Pattern
```go
logger := logo.New(
    logo.WithConsole(),
    logo.WithLevel(logo.LevelInfo),
    logo.WithCaller(),
)
logger.Info("hello")
```

---

## Import Path Changes

### v1
```go
import (
    "log/slog"
    logger "github.com/aN0mad/go-logo/logo"
)
```

### v2
```go
import "github.com/aN0mad/go-logo/v2"

// Or with alias
import logo "github.com/aN0mad/go-logo/v2"
```

**Key Changes:**
- `/v2` suffix required in import path
- No longer imports `log/slog` (v2 has its own types)
- Simpler package structure (no nested `/logo`)

---

## Breaking Changes

### 1. **No Global Singleton (Primary Breaking Change)**

**v1:** Used `Init()` to configure a global logger accessed via `L()`
```go
logo.Init(logo.SetLevel(slog.LevelDebug))
log := logo.L()
```

**v2:** Explicit logger instances
```go
logger := logo.New(logo.WithLevel(logo.LevelDebug))
```

**Why:** Explicit instances are more testable, thread-safe, and allow multiple independent loggers in the same application.

### 2. **Different Builder Options**

v1 used functional options like `SetLevel()`, `AddSource()`, `UseJSON()`.  
v2 uses `With*` prefixed options like `WithLevel()`, `WithCaller()`, `WithJSON()`.

### 3. **Handler/Writer Architecture**

**v1:** Outputs configured via options like `AddFileOutput()`, `AddChannelOutput()`
```go
logo.Init(
    logo.AddFileOutput("app.log", 10, 3, 30, true),
)
```

**v2:** Explicit handler/writer pattern
```go
logger := logo.New(
    logo.WithFile("app.log",
        logo.WithMaxSize(10),
        logo.WithMaxBackups(3),
        logo.WithMaxAge(30),
        logo.WithCompress(true),
    ),
)
```

### 4. **Log Levels**

**v1:** Used `slog.Level` types
```go
logo.SetLevel(slog.LevelDebug)
logo.SetLevel(logo.LevelTrace) // Custom level
```

**v2:** Self-contained level system
```go
logo.WithLevel(logo.LevelDebug)
logo.WithLevel(logo.LevelTrace)
```

### 5. **Field/Attribute Syntax**

**v1:** Alternating key-value pairs
```go
log.Info("message", "key1", "value1", "key2", "value2")
```

**v2:** Explicit `Field` constructor
```go
logger.Info("message", 
    logo.F("key1", "value1"),
    logo.F("key2", "value2"),
)
```

**Why:** Type safety, better IDE autocomplete, less error-prone.

### 6. **Configuration Profiles**

**v2 Only Feature:** Built-in CLI configuration and profiles
```go
cfg := logo.NewConfig()
cfg.ApplyProfile(logo.ProfileDev)
logger, _ := cfg.BuildLogger()
```

This has no equivalent in v1.

---

## API Migration Table

| v1 API | v2 API | Notes |
|--------|--------|-------|
| `logo.Init(opts...)` | `logo.New(opts...)` | No global state in v2 |
| `logo.L()` | Use logger instance | No global accessor |
| `logo.SetLevel(lvl)` | `logo.WithLevel(lvl)` | Option renamed |
| `logo.AddSource()` | `logo.WithCaller()` | Renamed for clarity |
| `logo.UseJSON(pretty)` | `logo.WithJSON()` | Simplified |
| `logo.DisableColors()` | `logo.WithNoColor()` | Renamed |
| `logo.EnableStackTraces()` | Not in v2 core | Removed (trace level shows stack) |
| `logo.AddFileOutput(...)` | `logo.WithFile(...)` | Different parameters |
| `logo.AddConsoleOutput()` | `logo.WithConsole()` | Renamed |
| `logo.AddChannelOutput(ch)` | `logo.WithChannel(ch)` | Renamed |
| `logo.DisableConsole()` | Don't use `WithConsole()` | Implicit |
| `logo.Close()` | `logger.Close()` | Instance method |
| `log.Trace(msg, ...)` | `logger.Trace(msg, ...)` | Same |
| `log.Debug(msg, ...)` | `logger.Debug(msg, ...)` | Same |
| `log.Info(msg, ...)` | `logger.Info(msg, ...)` | Same |
| `log.Warn(msg, ...)` | `logger.Warn(msg, ...)` | Same |
| `log.Error(msg, ...)` | `logger.Error(msg, ...)` | Same |
| `log.Fatal(msg, ...)` | `logger.Fatal(msg, ...)` | Same |
| `slog.LevelDebug` | `logo.LevelDebug` | Self-contained |
| `logo.LevelTrace` | `logo.LevelTrace` | Self-contained |
| `logo.LevelFatal` | `logo.LevelFatal` | Self-contained |

---

## Common Migration Patterns

### Pattern 1: Basic Console Logging

**v1:**
```go
logger.Init(
    logger.SetLevel(slog.LevelInfo),
)
log := logger.L()
log.Info("Application started", "version", "1.0.0")
```

**v2:**
```go
logger := logo.New(
    logo.WithConsole(),
    logo.WithLevel(logo.LevelInfo),
)
logger.Info("Application started", logo.F("version", "1.0.0"))
```

### Pattern 2: JSON Logging with File Output

**v1:**
```go
logger.Init(
    logger.SetLevel(slog.LevelDebug),
    logger.UseJSON(false),
    logger.AddSource(),
    logger.AddFileOutput("logs/app.log", 10, 3, 30, true),
)
defer logger.Close()

log := logger.L()
log.Debug("Debug information", "user_id", 123)
```

**v2:**
```go
logger := logo.New(
    logo.WithJSON(),
    logo.WithLevel(logo.LevelDebug),
    logo.WithCaller(),
    logo.WithFile("logs/app.log",
        logo.WithMaxSize(10),
        logo.WithMaxBackups(3),
        logo.WithMaxAge(30),
        logo.WithCompress(true),
    ),
)
defer logger.Close()

logger.Debug("Debug information", logo.F("user_id", 123))
```

### Pattern 3: Multiple Outputs (Console + File)

**v1:**
```go
logger.Init(
    logger.SetLevel(slog.LevelInfo),
    logger.AddConsoleOutput(),
    logger.AddFileOutput("logs/app.log", 10, 3, 30, false),
)
```

**v2:**
```go
logger := logo.New(
    logo.WithConsole(),
    logo.WithFile("logs/app.log",
        logo.WithMaxSize(10),
        logo.WithMaxBackups(3),
        logo.WithMaxAge(30),
    ),
    logo.WithLevel(logo.LevelInfo),
)
```

### Pattern 4: Trace Logging

**v1:**
```go
logger.Init(
    logger.EnableLogLevelTrace(),
    logger.EnableStackTraces(),
)
log := logger.L()
log.Trace("Detailed trace information")
```

**v2:**
```go
logger := logo.New(
    logo.WithConsole(),
    logo.WithLevel(logo.LevelTrace),
)
logger.Trace("Detailed trace information")
```

### Pattern 5: Structured Logging with Complex Data

**v1:**
```go
log.Info("User action",
    "user_id", 12345,
    "action", "login",
    "ip", "192.168.1.1",
    "timestamp", time.Now(),
)
```

**v2:**
```go
logger.Info("User action",
    logo.F("user_id", 12345),
    logo.F("action", "login"),
    logo.F("ip", "192.168.1.1"),
    logo.F("timestamp", time.Now()),
)
```

### Pattern 6: Channel Output

**v1:**
```go
ch := make(chan string, 100)
logger.Init(
    logger.AddChannelOutput(ch),
)
```

**v2:**
```go
ch := make(chan string, 100)
logger := logo.New(
    logo.WithChannel(ch),
)
```

### Pattern 7: CLI-Driven Configuration (v2 Only)

**New in v2:**
```go
cfg := logo.NewConfig()
cfg.RegisterFlags(nil) // Uses flag.CommandLine
flag.Parse()

logger, err := cfg.BuildLogger()
if err != nil {
    panic(err)
}
defer logger.Close()
```

```bash
# Run with flags
./app --log-level=debug --log-format=json --log-file=app.log
```

### Pattern 8: Context-Aware HTTP Logging (v2 Only)

**New in v2 with Fiber:**
```go
import (
    "time"
    "github.com/aN0mad/go-logo/v2"
    "github.com/gofiber/fiber/v2"
)

logger := logo.New(logo.WithConsole())
app := fiber.New()

// Custom middleware - simple and framework-agnostic
app.Use(func(c *fiber.Ctx) error {
    start := time.Now()
    requestID := logo.GenerateUUIDv4()
    ctx := logo.WithRequestID(c.Context(), requestID)
    c.SetUserContext(ctx)
    
    logger.WithFields(map[string]interface{}{
        "request_id": requestID,
        "method": c.Method(),
        "path": c.Path(),
    }).Info("request started")
    
    err := c.Next()
    
    logger.WithFields(map[string]interface{}{
        "request_id": requestID,
        "duration_ms": time.Since(start).Milliseconds(),
    }).Info("request completed")
    
    return err
})

app.Get("/", func(c *fiber.Ctx) error {
    log := logo.FromContext(c.Context())
    log.Info("Request handled")
    return c.SendString("OK")
})
```

**Benefits:** No middleware dependencies, full customization control.

---

## Feature Comparison

| Feature | v1 | v2 | Notes |
|---------|----|----|-------|
| Console output | ✅ | ✅ | Improved colors in v2 |
| JSON output | ✅ | ✅ | Stable schema in v2 |
| File rotation | ✅ | ✅ | Same lumberjack backend |
| Channel output | ✅ | ✅ | Same functionality |
| Trace level | ✅ | ✅ | Same |
| Fatal level | ✅ | ✅ | Same |
| Source/caller info | ✅ | ✅ | Renamed option |
| Color support | ✅ | ✅ | Better auto-detection in v2 |
| Multiple outputs | ✅ | ✅ | Cleaner API in v2 |
| Stack traces | ✅ | ⚠️ | Only on Trace level in v2 |
| CLI configuration | ❌ | ✅ | New in v2 |
| Profiles (dev/prod) | ❌ | ✅ | New in v2 |
| HTTP middleware | ❌ | ✅ | New in v2 (Fiber/Gin/Echo) |
| Context integration | ❌ | ✅ | New in v2 |
| Request ID tracking | ❌ | ✅ | New in v2 |
| Global singleton | ✅ | ❌ | Removed for testability |
| Custom handlers | ✅ | ✅ | Different API |

---

## Upgrade Checklist

Use this checklist to ensure a smooth migration:

### Pre-Migration
- [ ] Review all `logo.Init()` calls in your codebase
- [ ] Identify all `logo.L()` usages
- [ ] Document custom handler implementations (if any)
- [ ] Note any channel output configurations
- [ ] Verify file rotation settings

### Code Changes
- [ ] Update import path from `/logo` to `/v2`
- [ ] Replace `logo.Init()` with `logo.New()`
- [ ] Replace all `logo.L()` with explicit logger instances
- [ ] Update option names (`SetLevel` → `WithLevel`, etc.)
- [ ] Change field syntax to use `logo.F()`
- [ ] Update log level constants (remove `slog.` prefix)
- [ ] Replace `logo.Close()` with `logger.Close()`
- [ ] Update file output configuration if used

### Testing
- [ ] Run all tests to catch compilation errors
- [ ] Verify log output format (console and file)
- [ ] Test file rotation behavior
- [ ] Verify JSON schema if using structured logging
- [ ] Test CLI flags if using config profiles
- [ ] Validate middleware integration if using HTTP features

### Optional Enhancements
- [ ] Consider using CLI configuration for easier deployment
- [ ] Add HTTP middleware for request-scoped logging
- [ ] Use dev/prod profiles for environment-specific configs
- [ ] Implement request ID tracking for distributed tracing
- [ ] Refactor to use context-aware logging helpers

### Deployment
- [ ] Update documentation
- [ ] Update CI/CD pipelines if they parse logs
- [ ] Test in staging environment
- [ ] Monitor log output after deployment
- [ ] Verify file paths and rotation settings in production

---

## Troubleshooting

### "undefined: slog.LevelDebug"
**Solution:** Use `logo.LevelDebug` instead of `slog.LevelDebug`

### "undefined: logo.L"
**Solution:** Replace global `logo.L()` with your logger instance

### "too many arguments to logo.Info"
**Solution:** Wrap fields with `logo.F()`: `logo.F("key", value)`

### "cannot use logo.SetLevel as type logo.Option"
**Solution:** Rename to `logo.WithLevel`

### File rotation not working
**Solution:** Verify file writer options use `WithMaxSize()`, `WithMaxBackups()`, etc.

### Colors not showing in terminal
**Solution:** v2 auto-detects TTY. Use `logo.WithColor()` to force enable.

---

## Need Help?

- **GitHub Issues:** https://github.com/aN0mad/go-logo/issues
- **Examples:** See `/examples` directory for working code
- **API Docs:** https://pkg.go.dev/github.com/aN0mad/go-logo/v2

---

## Benefits of Upgrading

1. **Better Testability** - Explicit instances make unit testing easier
2. **CLI Configuration** - No code changes needed for log level/format changes
3. **HTTP Integration** - Built-in middleware for Fiber, Gin, Echo
4. **Request Tracking** - Automatic request ID generation and propagation
5. **Type Safety** - `logo.F()` provides better IDE support
6. **Performance** - Improved routing and zero-allocation paths
7. **Stability** - Locked JSON schema guarantees
8. **Flexibility** - Multiple independent loggers in same app
