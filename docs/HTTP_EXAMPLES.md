# HTTP Integration Examples

These examples demonstrate how to integrate go-logo v2 with popular Go web frameworks.

## Framework-Agnostic Design

go-logo v2 has **zero web framework dependencies**. The examples below show simple middleware implementations that you can customize for your needs.

## Common Patterns

All examples demonstrate:
- **Request ID Generation** - Using `logo.GenerateUUIDv4()` or custom generators
- **Context-Aware Logging** - Using `logo.WithRequestID()` and `logo.FromContext()`
- **Request/Response Logging** - Logging method, path, status, duration
- **Structured Fields** - Adding custom fields to log entries

## Writing Your Own Middleware

Each framework has slightly different middleware APIs, but the pattern is similar:

1. Generate or extract a request ID
2. Add the request ID to the context using `logo.WithRequestID(ctx, requestID)`
3. Log the request start
4. Call the next handler
5. Log the request completion with timing and status

## Fiber v2 Example

```bash
cd examples/http_fiber
go mod tidy
go run main.go
```

**Custom Middleware Pattern:**
```go
func loggerMiddleware(logger *logo.Logger) fiber.Handler {
    return func(c *fiber.Ctx) error {
        start := time.Now()
        
        // Generate request ID
        requestID := logo.GenerateUUIDv4()
        c.Locals("request_id", requestID)
        
        // Add to context
        ctx := logo.WithRequestID(c.Context(), requestID)
        c.SetUserContext(ctx)
        
        // Log request
        logger.WithFields(map[string]interface{}{
            "request_id": requestID,
            "method": c.Method(),
            "path": c.Path(),
        }).Info("request started")
        
        err := c.Next()
        
        // Log response
        logger.WithFields(map[string]interface{}{
            "request_id": requestID,
            "status": c.Response().StatusCode(),
            "duration_ms": time.Since(start).Milliseconds(),
        }).Info("request completed")
        
        return err
    }
}

app.Use(loggerMiddleware(logger))
```

**Usage in Handlers:**
```go
// Get logger from context
log := logo.FromContext(c.Context())
log.Info("handling request")

// Add custom fields
log.WithFields(map[string]interface{}{
    "user_id": c.Params("id"),
}).Info("fetching user")
```

**Test:**
```bash
curl http://localhost:3000/
curl http://localhost:3000/user/123
curl -X POST http://localhost:3000/orders
```

## Echo Example

```bash
cd examples/http_echo
go mod init example/http_echo
go get github.com/aN0mad/go-logo/v2
go get github.com/labstack/echo/v4
go run main.go
```

**Custom Middleware Pattern:**
```go
func loggerMiddleware(logger *logo.Logger) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            start := time.Now()
            
            // Generate or extract request ID
            requestID := c.Request().Header.Get("X-Request-ID")
            if requestID == "" {
                requestID = logo.GenerateUUIDv4()
            }
            c.Response().Header().Set("X-Request-ID", requestID)
            
            // Add to context
            req := c.Request()
            ctx := logo.WithRequestID(req.Context(), requestID)
            c.SetRequest(req.WithContext(ctx))
            
            // Log request
            logger.WithFields(map[string]interface{}{
                "request_id": requestID,
                "method": c.Request().Method,
                "path": c.Request().URL.Path,
            }).Info("request started")
            
            err := next(c)
            
            // Log response
            logger.WithFields(map[string]interface{}{
                "request_id": requestID,
                "status": c.Response().Status,
                "duration_ms": time.Since(start).Milliseconds(),
            }).Info("request completed")
            
            return err
        }
    }
}

e.Use(loggerMiddleware(logger))
```

**Usage in Handlers:**
```go
// Get logger from context
log := logo.FromContext(c.Request().Context())
log.Info("handling request")

// Add custom fields
log.WithFields(map[string]interface{}{
    "user_id": c.Param("id"),
}).Info("fetching user")
```

**Test:**
```bash
curl http://localhost:3002/
curl http://localhost:3002/user/456
curl -X POST http://localhost:3002/orders
```

## Gin Example

```bash
cd examples/http_gin
go mod tidy
go run main.go
```

**Custom Middleware Pattern:**
```go
func loggerMiddleware(logger *logo.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        
        // Generate or extract request ID
        requestID := c.GetHeader("X-Request-ID")
        if requestID == "" {
            requestID = logo.GenerateUUIDv4()
        }
        c.Header("X-Request-ID", requestID)
        
        // Add to context
        ctx := logo.WithRequestID(c.Request.Context(), requestID)
        c.Request = c.Request.WithContext(ctx)
        
        // Log request
        logger.WithFields(map[string]interface{}{
            "request_id": requestID,
            "method": c.Request.Method,
            "path": c.Request.URL.Path,
        }).Info("request started")
        
        c.Next()
        
        // Log response
        logger.WithFields(map[string]interface{}{
            "request_id": requestID,
            "status": c.Writer.Status(),
            "duration_ms": time.Since(start).Milliseconds(),
        }).Info("request completed")
    }
}

router.Use(loggerMiddleware(logger))
```

**Usage in Handlers:**
```go
// Get logger from context
log := logo.FromContext(c.Request.Context())
log.Info("handling request")

// Add custom fields
log.WithFields(map[string]interface{}{
    "user_id": c.Param("id"),
}).Info("fetching user")
```

**Test:**
```bash
curl http://localhost:3001/
curl http://localhost:3001/user/789
curl -X POST http://localhost:3001/orders
```

## Fiber v3 Example

```bash
cd examples/http_fiber_v3
go mod tidy
go run main.go
```

**Custom Middleware Pattern:**
```go
func loggerMiddleware(logger *logo.Logger) fiber.Handler {
    return func(c fiber.Ctx) error {
        start := time.Now()
        
        // Generate request ID
        requestID := logo.GenerateUUIDv4()
        c.Locals("request_id", requestID)
        
        // Add to context
        ctx := logo.WithRequestID(c.Context(), requestID)
        c.SetUserContext(ctx)
        
        // Log request
        logger.WithFields(map[string]interface{}{
            "request_id": requestID,
            "method": c.Method(),
            "path": c.Path(),
        }).Info("request started")
        
        err := c.Next()
        
        // Log response
        logger.WithFields(map[string]interface{}{
            "request_id": requestID,
            "status": c.Response().StatusCode(),
            "duration_ms": time.Since(start).Milliseconds(),
        }).Info("request completed")
        
        return err
    }
}

app.Use(loggerMiddleware(logger))
```

**Note:** Fiber v3 has minor API differences from v2 (mainly `fiber.Ctx` interface changes). The middleware pattern remains the same.

**Test:**
```bash
curl http://localhost:3003/
curl http://localhost:3003/user/321
curl -X POST http://localhost:3003/orders
```

## Additional Frameworks

The same pattern works with any Go web framework:

### Standard net/http

```go
func loggerMiddleware(logger *logo.Logger) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()
            
            // Generate request ID
            requestID := logo.GenerateUUIDv4()
            w.Header().Set("X-Request-ID", requestID)
            
            // Add to context
            ctx := logo.WithRequestID(r.Context(), requestID)
            r = r.WithContext(ctx)
            
            // Log request
            logger.WithFields(map[string]interface{}{
                "request_id": requestID,
                "method": r.Method,
                "path": r.URL.Path,
            }).Info("request started")
            
            // Wrap response writer to capture status
            wrapped := &statusWriter{ResponseWriter: w, status: 200}
            next.ServeHTTP(wrapped, r)
            
            // Log response
            logger.WithFields(map[string]interface{}{
                "request_id": requestID,
                "status": wrapped.status,
                "duration_ms": time.Since(start).Milliseconds(),
            }).Info("request completed")
        })
    }
}

type statusWriter struct {
    http.ResponseWriter
    status int
}

func (w *statusWriter) WriteHeader(status int) {
    w.status = status
    w.ResponseWriter.WriteHeader(status)
}
```

### Chi Router

```go
func loggerMiddleware(logger *logo.Logger) func(http.Handler) http.Handler {
    // Same pattern as net/http above
    return func(next http.Handler) http.Handler {
        // ... (same implementation)
    }
}

r := chi.NewRouter()
r.Use(loggerMiddleware(logger))
```

## Key Concepts

### Context-Aware Logging

The logger can be attached to the request context for easy access:

```go
// Add logger to context
ctx := logo.WithRequestID(r.Context(), requestID)

// Later, retrieve logger
log := logo.FromContext(ctx)
log.Info("processing request")
```

### Structured Fields

Add contextual information to log entries:

```go
log.WithFields(map[string]interface{}{
    "user_id": "123",
    "action": "update_profile",
    "ip": r.RemoteAddr,
}).Info("user action")
```

### Request ID Tracking

Generate and propagate request IDs:

```go
// Generate
requestID := logo.GenerateUUIDv4()
// Or: logo.GenerateShortID()
// Or: logo.GenerateSequentialID("REQ-")

// Add to context
ctx := logo.WithRequestID(r.Context(), requestID)

// Retrieve
requestID := logo.GetRequestID(ctx)
```

## Output Examples

**Console Output (dev mode):**
```
2024-01-15 14:32:10 INFO  starting server                          port=3000
2024-01-15 14:32:15 INFO  request started                          method=GET path=/ request_id=550e8400-e29b-41d4
2024-01-15 14:32:15 INFO  handling root request                    request_id=550e8400-e29b-41d4
2024-01-15 14:32:15 INFO  request completed                        request_id=550e8400-e29b-41d4 status=200 duration_ms=1
```

**JSON Output (prod mode):**
```json
{"ts":"2024-01-15T14:32:10Z","level":"INFO","msg":"starting server","port":3000}
{"ts":"2024-01-15T14:32:15Z","level":"INFO","msg":"request started","method":"GET","path":"/","request_id":"550e8400-e29b-41d4"}
{"ts":"2024-01-15T14:32:15Z","level":"INFO","msg":"handling root request","request_id":"550e8400-e29b-41d4"}
{"ts":"2024-01-15T14:32:15Z","level":"INFO","msg":"request completed","request_id":"550e8400-e29b-41d4","status":200,"duration_ms":1}
```

## Benefits of This Approach

✅ **Zero dependencies** - go-logo doesn't depend on any web frameworks  
✅ **Full control** - Customize middleware to your exact needs  
✅ **Simple** - ~30 lines of code per framework  
✅ **Flexible** - Easy to add/remove logging points  
✅ **Testable** - Mock logger in tests without framework coupling

## See Also

- [Context Package Documentation](http.md)
- [Request ID Generators](http.md#request-id-generators)
- [CLI Configuration Examples](CLI_EXAMPLES.md)
e.Use(middleware.Logger())

// After
e.Use(echologo.Middleware(logger))
log := echologo.Logger(c)
```

## Best Practices

1. **Create one logger instance** and share across middleware
2. **Use context helpers** for request-scoped fields
3. **Configure request ID format** based on your needs (UUID for distributed systems, short for readability)
4. **Skip health checks** to reduce log noise
5. **Use structured fields** instead of string concatenation
6. **Leverage request ID** for distributed tracing correlation

## Performance

Custom middleware implementations shown above are highly efficient:
- Request ID generation: ~100ns per request (UUID)
- Context storage: zero allocations
- Minimal overhead: <1ms per request

The middleware pattern is simple (~30 lines) and has negligible performance impact.
