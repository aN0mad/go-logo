# Migration Examples: v1 → v2

This directory contains side-by-side examples showing how to migrate from go-logo v1 to v2.

## Files

### Basic Migration
- **v1_style.go.example** - Typical v1 code (reference only, won't compile)
- **v2_style.go** - Same functionality in v2

### Specific Use Cases
- **json_migration.go** - JSON logging migration
- **cli_config.go** - NEW in v2: CLI configuration
- **http_middleware.go** - NEW in v2: HTTP request-scoped logging

## Running the Examples

### Basic v2 Style
```bash
cd examples/migration_v1_to_v2
go mod init example/migration
go get github.com/aN0mad/go-logo/v2
go run v2_style.go
```

### JSON Migration
```bash
go run json_migration.go
cat logs/app-json.log
```

### CLI Configuration
```bash
# Default console output
go run cli_config.go

# Debug level with JSON
go run cli_config.go --log-level=debug --log-format=json

# Write to file with rotation
go run cli_config.go \
  --log-file=logs/app.log \
  --log-max-size=10 \
  --log-max-backups=3 \
  --log-compress=true

# Use built-in dev profile
go run cli_config.go --log-profile=dev

# Use environment variables
LOG_LEVEL=debug LOG_FORMAT=json go run cli_config.go
```

### HTTP Middleware (Fiber)
```bash
# Install dependencies
go get github.com/gofiber/fiber/v2

# Run server (includes inline middleware example)
go run http_middleware.go

# In another terminal, test endpoints
curl http://localhost:3000/users/123
curl -X POST http://localhost:3000/orders
curl http://localhost:3000/health
```

**Note:** v2 doesn't include middleware packages. The example shows how to write custom middleware (simple and flexible).

## Key Differences to Notice

### 1. Import Path
```diff
- import logger "github.com/aN0mad/go-logo/logo"
+ import logo "github.com/aN0mad/go-logo/v2"
```

### 2. Logger Creation
```diff
- logger.Init(logger.SetLevel(slog.LevelDebug))
- log := logger.L()
+ logger := logo.New(logo.WithLevel(logo.LevelDebug))
```

### 3. Field Syntax
```diff
- log.Info("message", "key", "value")
+ logger.Info("message", logo.F("key", "value"))
```

### 4. Options Naming
```diff
- logger.SetLevel(slog.LevelDebug)
+ logo.WithLevel(logo.LevelDebug)

- logger.AddSource()
+ logo.WithCaller()

- logger.UseJSON(false)
+ logo.WithJSON()

- logger.AddFileOutput("app.log", 10, 3, 30, true)
+ logo.WithFile("app.log", 
+   logo.WithMaxSize(10),
+   logo.WithMaxBackups(3),
+   logo.WithMaxAge(30),
+   logo.WithCompress(true))
```

### 5. Log Levels
```diff
- slog.LevelDebug
+ logo.LevelDebug

- logo.LevelTrace  // Still works in v1
+ logo.LevelTrace  // Self-contained in v2
```

## New Features in v2

### CLI Configuration
```go
cfg := logo.NewConfig()
cfg.RegisterFlags(nil)
flag.Parse()
logger, _ := cfg.BuildLogger()
```

Run with flags:
```bash
./app --log-level=debug --log-format=json --log-file=app.log
```

### Profiles
```go
cfg := logo.NewConfig()
cfg.ApplyProfile(logo.ProfileDev)  // or ProfileProd
logger, _ := cfg.BuildLogger()
```

### HTTP Middleware
```go
app.Use(fiberlogo.Middleware(logger))

// In handlers
log := fiberlogo.Logger(c)
requestID := fiberlogo.RequestID(c)
```

### Context Fields
```go
log := logo.ContextWithFields(ctx,
    logo.F("user_id", 123),
    logo.F("tenant", "acme"),
)
```

## Migration Checklist

For each file using v1:

- [ ] Update import path to `/v2`
- [ ] Replace `logger.Init()` with `logo.New()`
- [ ] Remove `logger.L()` calls, use logger instance
- [ ] Update option names (`SetLevel` → `WithLevel`)
- [ ] Wrap fields with `logo.F()`
- [ ] Update log level constants
- [ ] Change `logger.Close()` to `logger.Close()` on instance
- [ ] Test output format

## Tips

1. **Start with basic console logging** - Get the simplest case working first
2. **Use CLI config for flexibility** - Easier than hardcoding
3. **Leverage profiles** - Use dev/prod profiles for environment-specific configs
4. **Add HTTP middleware** - Get request-scoped logging for free
5. **Use context helpers** - Propagate fields throughout request lifecycle

## Need Help?

See the main migration guide: `/docs/migration-v1-to-v2.md`
