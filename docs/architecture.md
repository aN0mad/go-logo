# Architecture

go-logo follows a clean, modular architecture that separates concerns and allows for flexible configuration.

## Component Overview

```
Application Code
      ↓
   Logger (immutable, thread-safe)
      ↓
   Router (fan-out dispatcher)
      ↓
   Handlers (level filtering, formatting)
      ↓
   Writers (console, file, custom)
```

## Core Components

### Logger

The `Logger` is the main entry point for application code. Key characteristics:

- **Immutable** - Creating child loggers or adding fields returns new logger instances
- **Thread-Safe** - Safe for concurrent use across goroutines
- **Field Inheritance** - Child loggers inherit fields from parents
- **Named Loggers** - Support for logger hierarchies via names

```go
logger := logo.New(logo.WithRouter(router))
childLogger := logger.With(logo.F("service", "api"))
```

### Router

The `Router` implements a fan-out pattern, dispatching log entries to multiple handlers simultaneously.

- **Multiple Handlers** - Send logs to console and file at the same time
- **Thread-Safe** - Concurrent routing without locks on the hot path
- **Dynamic** - Add handlers at runtime if needed

```go
router := logo.NewRouter(consoleHandler, fileHandler)
```

### Handler

Handlers define three key aspects of log processing:

1. **Level Filtering** - Minimum log level to process
2. **Formatting** - How entries are formatted (pretty, text, JSON)
3. **Destination** - Where logs are written (stdout, file, etc.)

Built-in handlers:
- `ConsoleHandler` - Pretty console output with colors
- `FileHandler` - File output with rotation
- `JSONHandler` - JSON-formatted output
- `WriterHandler` - Generic writer-based handler

```go
handler := logo.NewConsoleHandler(
    logo.WithConsoleLevel(logo.LevelInfo),
)
```

### Formatter

Formatters convert log entries into byte streams:

- `PrettyFormatter` - Human-readable colored output
- `TextFormatter` - Plain text output
- `JSONFormatter` - Structured JSON output

### Entry

The `Entry` represents a single log event:

```go
type Entry struct {
    Timestamp time.Time
    Level     Level
    Message   string
    Logger    string
    Fields    map[string]interface{}
    Caller    *Caller
}
```

Entries are immutable once created and cloneable for handler isolation.

## Data Flow

1. Application calls `logger.Info("message", logo.F("key", "value"))`
2. Logger creates an `Entry` with timestamp, level, message, and fields
3. Logger optionally adds caller information (file, line, function)
4. Logger sends entry to Router
5. Router fans out to all registered handlers
6. Each handler:
   - Checks if entry level meets its minimum level
   - Formats the entry using its formatter
   - Writes formatted bytes to its writer
7. Entry processing completes

## Concurrency Model

- **Loggers** - Read-only after creation, safe to share across goroutines
- **Router** - Uses RWMutex for handler list access
- **Handlers** - Use mutexes to protect writer access
- **No Channels** - Synchronous logging for predictable ordering

## Extension Points

For detailed API documentation of all types and interfaces, see the [API Reference](doc.md).

### Custom Handlers

Implement the `Handler` interface:

```go
type Handler interface {
    Handle(entry *Entry) bool
    Level() Level
    Enabled(level Level) bool
    Close() error
}
```

### Custom Formatters

Implement the `Formatter` interface:

```go
type Formatter interface {
    Format(entry *Entry) ([]byte, error)
}
```

### Custom Writers

Any `io.Writer` can be used with `WriterHandler`:

```go
handler := logo.NewWriterHandler(myWriter, formatter, level)
```

## Design Principles

1. **Simplicity** - Clear, straightforward APIs
2. **Immutability** - Loggers don't change after creation
3. **Composability** - Mix and match handlers, formatters, and writers
4. **Zero Allocations** - Minimal allocations in hot paths (where possible)
5. **No Reflection** - Type-safe, compile-time verified code
6. **Stable Schema** - JSON output schema is versioned and stable

## Performance Considerations

- **Synchronous Logging** - Logs are written synchronously for ordering guarantees
- **Level Filtering** - Early filtering at handler level avoids formatting costs
- **Field Reuse** - Field maps are pre-allocated and reused where possible
- **String Building** - Efficient string formatting with minimal allocations

## Testing Strategy

- **Golden Files** - Schema stability validated via golden file tests
- **Unit Tests** - Comprehensive coverage of all components
- **Integration Tests** - End-to-end logging scenarios
- **Benchmark Tests** - Performance regression detection

