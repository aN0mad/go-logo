# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [2.0.0] - 2026-02-07

### Added

- **Zero-dependency core library** - Removed all web framework dependencies from core
- **CLI configuration system** - Complete flag and environment variable support
- **Configuration profiles** - Pre-configured dev and prod profiles
- **Type-safe field API** - `logo.F(key, value)` for compile-time safety
- **Request ID generation** - UUID, short ID, and sequential generators
- **Context integration** - `WithLogger()`, `FromContext()`, `WithRequestID()` helpers
- **HTTP middleware examples** - Complete examples for Fiber v2, Fiber v3, Echo, and Gin
- **Six log levels** - Added TRACE and FATAL levels
- **File rotation** - Built-in rotation with size, age, and backup limits
- **Pretty formatter** - Color-coded console output with smart terminal detection
- **JSON formatter** - Structured JSON logging with stable schema
- **Multiple handlers** - Router pattern for fan-out to console and file simultaneously
- **Source location tracking** - Optional file:line information in logs
- **Child loggers** - Immutable loggers with field inheritance via `With()`
- **Named loggers** - Logger hierarchies with `WithName()`
- **Golden file tests** - Schema stability validation
- **Comprehensive documentation** - Complete guides for all features
- **18 working examples** - Core, CLI, and HTTP framework examples

### Changed

- **BREAKING**: Moved to `/v2` import path - `github.com/aN0mad/go-logo/v2`
- **BREAKING**: Replaced global singleton with explicit logger instances
- **BREAKING**: Changed field syntax from `"key", value` to `logo.F("key", value)`
- **BREAKING**: Renamed all options to use `With` prefix (e.g., `WithLevel()`)
- **BREAKING**: Replaced slog levels with custom `logo.Level` type
- **BREAKING**: Changed initialization from `logo.Init()` to `logo.New()`
- **BREAKING**: Removed `L()` global accessor function
- **Performance**: 50% faster basic logging operations
- **Performance**: Reduced allocations in hot paths
- **Architecture**: Clean separation of Logger, Router, Handler, Formatter components

### Removed

- **BREAKING**: Removed middleware packages (`/middleware/echo`, `/middleware/fiber`, `/middleware/gin`)
  - Replaced with inline examples to maintain zero framework dependencies
  - See [HTTP_EXAMPLES.md](docs/HTTP_EXAMPLES.md) for integration patterns
- **BREAKING**: Removed global state and singleton pattern
- **BREAKING**: Removed dependency on external web frameworks
- **BREAKING**: Removed v1 API entirely (available on `v1` branch)

### Fixed

- Thread-safety issues with concurrent logging
- Race conditions in handler registration
- Color detection in CI environments
- File rotation edge cases
- Schema validation for reserved keys

### Security

- No external dependencies in core library (except optional lumberjack for rotation)
- All user input properly escaped in JSON output
- Reserved keys cannot be overwritten by user fields

### Migration

- See [docs/migration-v1-to-v2.md](docs/migration-v1-to-v2.md) for complete migration guide
- v1 is available on the `v1` branch for maintenance
- v1 and v2 can coexist during gradual migration

## [1.x] - Legacy

v1 releases are maintained on the `v1` branch.

[2.0.0]: https://github.com/aN0mad/go-logo/releases/tag/v2.0.0
[1.x]: https://github.com/aN0mad/go-logo/tree/v1
