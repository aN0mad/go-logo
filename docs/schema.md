# JSON Schema Reference

go-logo v2 uses a stable, well-defined JSON schema for structured logging.

For complete API documentation, see the [API Reference](doc.md).

## Overview

The JSON schema defines required fields, field types, and reserved keys to ensure consistent, parseable log output across all services.

## Schema Version

- **Version**: 2.0
- **Stability**: Stable - no breaking changes to existing fields
- **Format**: JSON (RFC 8259)
- **Encoding**: UTF-8

## Required Fields

Every log entry includes these fields:

| Field | Type | Description | Format/Example |
|-------|------|-------------|----------------|
| `ts` | string | Timestamp | RFC3339Nano UTC: `"2026-02-07T15:04:05.123456789Z"` |
| `level` | string | Log level | `"TRACE"`, `"DEBUG"`, `"INFO"`, `"WARN"`, `"ERROR"`, `"FATAL"` |
| `msg` | string | Log message | `"user login successful"` |
| `logger` | string | Logger name | `"http"`, `"database"`, `"default"` |

### Field Specifications

#### ts (Timestamp)

- **Format**: RFC3339Nano
- **Timezone**: Always UTC (Z suffix)
- **Precision**: Nanoseconds
- **Example**: `"2026-02-07T15:04:05.123456789Z"`

```json
{
  "ts": "2026-02-07T15:04:05.123456789Z"
}
```

#### level

- **Type**: String (uppercase)
- **Values**: `TRACE`, `DEBUG`, `INFO`, `WARN`, `ERROR`, `FATAL`
- **Example**: `"INFO"`

```json
{
  "level": "INFO"
}
```

#### msg

- **Type**: String
- **Content**: Human-readable log message
- **Example**: `"request processed successfully"`

```json
{
  "msg": "request processed successfully"
}
```

#### logger

- **Type**: String
- **Content**: Logger name/identifier
- **Default**: `"default"`
- **Example**: `"http"`, `"database"`, `"api"`

```json
{
  "logger": "http"
}
```

## Optional Fields

### source

Source location information (when enabled with `WithSource(true)`):

- **Format**: `"file:line"`
- **Example**: `"main.go:42"`

```json
{
  "source": "handler.go:127"
}
```

### User-Defined Fields

Any non-reserved field name can be used for custom structured data:

```json
{
  "ts": "2026-02-07T15:04:05.123Z",
  "level": "INFO",
  "msg": "user login",
  "logger": "auth",
  "user_id": 12345,
  "ip": "192.168.1.1",
  "success": true,
  "duration_ms": 42
}
```

## Reserved Keys

These field names are reserved and cannot be overwritten by user fields:

- `ts` - Timestamp
- `level` - Log level
- `msg` - Message
- `logger` - Logger name
- `source` - Source location
- `caller` - Alternative name for source

Attempting to use reserved keys will result in those fields being filtered out:

```go
// Reserved key "ts" is filtered out
log.Info("message", 
    logo.F("ts", "custom-value"),  // Filtered
    logo.F("user_id", 123),        // Included
)
```

Output:
```json
{
  "ts": "2026-02-07T15:04:05.123Z",
  "level": "INFO",
  "msg": "message",
  "logger": "default",
  "user_id": 123
}
```

## Field Types

User-defined fields support all JSON-serializable types:

### Strings

```go
log.Info("message", logo.F("name", "John Doe"))
```

```json
{
  "name": "John Doe"
}
```

### Numbers

```go
log.Info("message", 
    logo.F("count", 42),
    logo.F("price", 19.99),
)
```

```json
{
  "count": 42,
  "price": 19.99
}
```

### Booleans

```go
log.Info("message", logo.F("success", true))
```

```json
{
  "success": true
}
```

### Arrays/Slices

```go
log.Info("message", logo.F("tags", []string{"api", "v2", "auth"}))
```

```json
{
  "tags": ["api", "v2", "auth"]
}
```

### Objects/Maps

```go
log.Info("message", logo.F("metadata", map[string]interface{}{
    "version": "2.0",
    "region": "us-east-1",
}))
```

```json
{
  "metadata": {
    "version": "2.0",
    "region": "us-east-1"
  }
}
```

### Null Values

```go
log.Info("message", logo.F("optional", nil))
```

```json
{
  "optional": null
}
```

## Complete Example

```json
{
  "ts": "2026-02-07T15:04:05.123456789Z",
  "level": "INFO",
  "msg": "HTTP request processed",
  "logger": "http",
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "method": "POST",
  "path": "/api/v2/users",
  "status": 201,
  "duration_ms": 45,
  "ip": "192.168.1.100",
  "user_agent": "Mozilla/5.0",
  "user_id": 12345,
  "tenant_id": "acme-corp"
}
```

## Recommended Fields

For consistency across services, use these recommended field names:

### HTTP/API

- `request_id` - Unique request identifier
- `method` - HTTP method (GET, POST, etc.)
- `path` - Request path
- `status` - HTTP status code
- `duration_ms` - Request duration in milliseconds
- `ip` - Client IP address
- `user_agent` - Client user agent

### Authentication

- `user_id` - User identifier
- `tenant_id` - Tenant/organization identifier
- `session_id` - Session identifier
- `role` - User role

### Errors

- `error` - Error message
- `error_type` - Error type/class
- `stack_trace` - Stack trace (use sparingly)

### Distributed Tracing

- `trace_id` - Distributed trace ID
- `span_id` - Span ID
- `parent_span_id` - Parent span ID

### Application

- `service` - Service name
- `version` - Application version
- `env` - Environment (dev, staging, prod)
- `host` - Hostname
- `region` - Cloud region

## Schema Validation

go-logo validates fields at runtime:

```go
// Validates reserved keys
logo.ValidateFields(fields)

// Filters reserved keys
filtered := logo.FilterReservedKeys(fields)

// Checks if key is reserved
if logo.IsReservedKey("ts") {
    // true
}
```

## Golden File Testing

Schema stability is validated through golden file tests:

```bash
# Run golden file tests
go test -v -run TestGoldenJSON

# Update golden files (when making intentional changes)
UPDATE_GOLDEN=1 go test -run TestGoldenJSON
```

Golden files in `testdata/golden/` capture expected output for regression detection.

## Schema Guarantees

go-logo provides these guarantees:

1. **Required Fields Always Present** - `ts`, `level`, `msg`, `logger` in every entry
2. **Consistent Field Types** - Field types don't change between versions
3. **RFC3339Nano Timestamps** - Standard, parseable timestamp format
4. **Reserved Keys Protected** - User fields cannot override schema fields
5. **Valid JSON** - All output is valid JSON (RFC 8259)
6. **UTF-8 Encoding** - Proper handling of Unicode characters
7. **No Breaking Changes** - Schema is append-only (new fields may be added)

## Parsing Examples

### Go

```go
type LogEntry struct {
    Timestamp string                 `json:"ts"`
    Level     string                 `json:"level"`
    Message   string                 `json:"msg"`
    Logger    string                 `json:"logger"`
    Fields    map[string]interface{} `json:",inline"`
}

var entry LogEntry
json.Unmarshal(logLine, &entry)
```

### Python

```python
import json

log_entry = json.loads(log_line)
timestamp = log_entry['ts']
level = log_entry['level']
message = log_entry['msg']
user_id = log_entry.get('user_id')
```

### jq

```bash
# Filter by level
cat app.log | jq 'select(.level == "ERROR")'

# Extract specific fields
cat app.log | jq '{time: .ts, message: .msg, user: .user_id}'

# Calculate average duration
cat app.log | jq -s 'map(.duration_ms) | add / length'
```

## Best Practices

1. **Use Consistent Field Names** - Standardize across services
2. **Avoid Reserved Keys** - Don't use `ts`, `level`, `msg`, `logger`, `source`
3. **Use Descriptive Names** - `user_id` not `uid`, `duration_ms` not `dur`
4. **Include Units** - `duration_ms`, `size_bytes`, `timeout_sec`
5. **Use Underscores** - `user_id` not `userId` or `user-id`
6. **Avoid Nesting** - Flat structure is easier to query
7. **Be Consistent** - Same event should have same fields
8. **Don't Log Secrets** - Never log passwords, tokens, or sensitive data

## Version Compatibility

### Adding Fields

New optional fields may be added in minor versions. Your parsers should:
- Ignore unknown fields
- Not fail on additional fields
- Use defaults for missing optional fields

### Deprecating Fields

Fields are never removed, only deprecated with warnings in documentation.

### Schema Evolution

- **Major version**: Breaking changes to existing fields (rare)
- **Minor version**: New optional fields added
- **Patch version**: No schema changes

Current schema version is locked for v2.x releases.

