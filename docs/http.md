# HTTP Integration

go-logo integrates seamlessly with popular Go web frameworks through custom middleware patterns.

For complete API documentation, see the [API Reference](doc.md).

## Overview

HTTP middleware responsibilities:
1. Generate or extract request IDs
2. Attach a request-scoped logger to the context
3. Log request lifecycle events (start, completion)
4. Make the logger available to handlers via context

## Framework Examples

go-logo provides complete, ready-to-use examples for popular frameworks:

- **[Fiber v2](../examples/http_fiber_v2/)** - High-performance framework
- **[Fiber v3](../examples/http_fiber_v3/)** - Fiber v3 beta
- **[Echo](../examples/http_echo/)** - Minimalist web framework
- **[Gin](../examples/http_gin/)** - Popular HTTP framework

See [HTTP_EXAMPLES.md](HTTP_EXAMPLES.md) for detailed implementations.

## Common Pattern

All HTTP middleware examples follow this pattern:

```go
func loggerMiddleware(logger *logo.Logger) MiddlewareFunc {
    return func(c Context) error {
        start := time.Now()
        
        // 1. Generate request ID
        requestID := logo.GenerateRequestID()
        
        // 2. Add request ID to context
        ctx := logo.WithRequestID(c.Context(), requestID)
        c.SetContext(ctx)
        
        // 3. Log request start
        logger.With(
            logo.F("request_id", requestID),
            logo.F("method", c.Method()),
            logo.F("path", c.Path()),
        ).Info("request started")
        
        // 4. Process request
        err := c.Next()
        
        // 5. Log request completion
        duration := time.Since(start)
        logger.With(
            logo.F("request_id", requestID),
            logo.F("status", c.StatusCode()),
            logo.F("duration_ms", duration.Milliseconds()),
        ).Info("request completed")
        
        return err
    }
}
```

## Using Context-Aware Logging

In your handlers, retrieve the logger from context:

```go
func handleRequest(c *fiber.Ctx) error {
    // Get logger from context
    log := logo.FromContext(c.Context())
    
    // Log with automatic request ID
    log.Info("processing user request",
        logo.F("user_id", getUserID(c)),
        logo.F("action", "fetch_profile"),
    )
    
    // Business logic...
    
    return c.JSON(response)
}
```

## Request ID Generation

go-logo provides multiple request ID generators:

### UUID Generator (Default)

```go
requestID := logo.GenerateRequestID()
// Returns: "550e8400-e29b-41d4-a716-446655440000"
```

### Short ID Generator

```go
logo.SetRequestIDGenerator(&logo.ShortIDGenerator{})
requestID := logo.GenerateRequestID()
// Returns: "a1b2c3d4e5f6g7h8"
```

### Sequential Generator (Testing)

```go
logo.SetRequestIDGenerator(logo.NewSequentialGenerator("req"))
requestID := logo.GenerateRequestID()
// Returns: "req-1", "req-2", "req-3"...
```

## Common Logged Fields

Standard fields logged by HTTP middleware:

| Field | Description | Example |
|-------|-------------|---------|
| `request_id` | Unique request identifier | `"550e8400-e29b-41d4..."` |
| `method` | HTTP method | `"GET"`, `"POST"` |
| `path` | Request path | `"/api/users/123"` |
| `status` | HTTP status code | `200`, `404` |
| `duration_ms` | Request duration in milliseconds | `42` |
| `ip` | Client IP address | `"192.168.1.1"` |
| `user_agent` | Client user agent | `"Mozilla/5.0..."` |

## Fiber Example

```go
package main

import (
    "github.com/aN0mad/go-logo/v2"
    "github.com/gofiber/fiber/v2"
)

func main() {
    logger := logo.New(logo.WithConsole())
    defer logger.Close()
    
    app := fiber.New()
    
    // Use middleware
    app.Use(loggerMiddleware(logger))
    
    app.Get("/", func(c *fiber.Ctx) error {
        log := logo.FromContext(c.Context())
        log.Info("handling request")
        return c.SendString("Hello!")
    })
    
    app.Listen(":3000")
}
```

## Echo Example

```go
package main

import (
    "github.com/aN0mad/go-logo/v2"
    "github.com/labstack/echo/v4"
)

func main() {
    logger := logo.New(logo.WithConsole())
    defer logger.Close()
    
    e := echo.New()
    
    // Use middleware
    e.Use(loggerMiddleware(logger))
    
    e.GET("/", func(c echo.Context) error {
        log := logo.FromContext(c.Request().Context())
        log.Info("handling request")
        return c.String(200, "Hello!")
    })
    
    e.Start(":3000")
}
```

## Gin Example

```go
package main

import (
    "github.com/aN0mad/go-logo/v2"
    "github.com/gin-gonic/gin"
)

func main() {
    logger := logo.New(logo.WithConsole())
    defer logger.Close()
    
    r := gin.New()
    
    // Use middleware
    r.Use(loggerMiddleware(logger))
    
    r.GET("/", func(c *gin.Context) {
        log := logo.FromContext(c.Request.Context())
        log.Info("handling request")
        c.String(200, "Hello!")
    })
    
    r.Run(":3000")
}
```

## Advanced Patterns

### Request-Scoped Fields

Add persistent fields to the logger for the entire request:

```go
func authMiddleware(c Context) {
    log := logo.FromContext(c.Context())
    
    userID := authenticate(c)
    
    // Create logger with user context
    userLogger := log.With(
        logo.F("user_id", userID),
        logo.F("tenant_id", getTenant(userID)),
    )
    
    // Update context with new logger
    ctx := logo.WithLogger(c.Context(), userLogger)
    c.SetContext(ctx)
    
    // All downstream handlers get logger with user fields
    c.Next()
}
```

### Error Logging

Log errors with full context:

```go
func handleRequest(c Context) error {
    log := logo.FromContext(c.Context())
    
    result, err := doSomething()
    if err != nil {
        log.Error("operation failed",
            logo.F("error", err.Error()),
            logo.F("operation", "doSomething"),
        )
        return c.Status(500).JSON(errorResponse)
    }
    
    return c.JSON(result)
}
```

### Correlation IDs

Propagate correlation IDs across services:

```go
func middleware(logger *logo.Logger) HandlerFunc {
    return func(c Context) {
        // Extract or generate correlation ID
        correlationID := c.Get("X-Correlation-ID")
        if correlationID == "" {
            correlationID = logo.GenerateRequestID()
        }
        
        // Add to context
        ctx := logo.WithRequestID(c.Context(), correlationID)
        c.SetContext(ctx)
        
        // Log with correlation ID
        logger.With(
            logo.F("correlation_id", correlationID),
        ).Info("request received")
        
        c.Next()
    }
}
```

## Best Practices

1. **Generate Request IDs Early** - Create them in the first middleware
2. **Use Context Everywhere** - Always retrieve logger from context in handlers
3. **Log Request Start and End** - Provides complete request lifecycle visibility
4. **Include Duration** - Track request latency
5. **Structure Your Fields** - Use consistent field names across services
6. **Don't Log Sensitive Data** - Avoid logging passwords, tokens, or PII
7. **Use Appropriate Levels** - Info for normal requests, Error for failures
8. **Consider Sampling** - In high-traffic scenarios, sample verbose logs

## Integration with Observability

go-logo's structured logging integrates well with observability tools:

- **Log Aggregation**: Use JSON format for ELK, Splunk, Datadog
- **Distributed Tracing**: Include trace_id and span_id fields from OpenTelemetry
- **Metrics**: Log duration_ms for latency metrics
- **APM**: Include request_id for correlation with application performance monitoring

Example with OpenTelemetry:

```go
import "go.opentelemetry.io/otel/trace"

func middleware(logger *logo.Logger) HandlerFunc {
    return func(c Context) {
        span := trace.SpanFromContext(c.Context())
        
        log := logger.With(
            logo.F("trace_id", span.SpanContext().TraceID().String()),
            logo.F("span_id", span.SpanContext().SpanID().String()),
        )
        
        ctx := logo.WithLogger(c.Context(), log)
        c.SetContext(ctx)
        
        c.Next()
    }
}
```

## Troubleshooting

**Logger not in context:**
- Ensure middleware is registered before route handlers
- Check that you're using `c.Context()` not `context.Background()`

**Request ID missing:**
- Verify middleware generates and adds request ID to context
- Use `logo.WithRequestID()` to add ID to context

**Duplicate logs:**
- Check that middleware isn't registered multiple times
- Ensure you're not calling `logger.Info()` and also letting default middleware log

