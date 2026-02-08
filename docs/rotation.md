# Log Rotation

go-logo provides built-in log rotation for file handlers using the lumberjack library.

## Overview

Log rotation automatically manages log file size and retention, preventing disk space issues in long-running applications.

## Features

- **Size-Based Rotation** - Rotate when file reaches a specified size
- **Age-Based Retention** - Automatically delete old log files
- **Backup Limits** - Keep only a specified number of rotated files
- **Compression** - Optionally compress rotated files with gzip
- **Thread-Safe** - Safe for concurrent writes

## Configuration

### Default Configuration

```go
log := logo.New(logo.WithJSONFile("app.log"))
```

Default settings:
- **MaxSize**: 100 MB
- **MaxAge**: 30 days
- **MaxBackups**: 10 files
- **Compress**: true
- **LocalTime**: false (UTC timestamps)

### Custom Configuration

```go
rotationConfig := &logo.RotationConfig{
    MaxSize:    10,   // 10 MB
    MaxAge:     7,    // 7 days
    MaxBackups: 5,    // 5 files
    Compress:   true,
    LocalTime:  false,
}

formatter := logo.NewJSONFormatter()
fileHandler, err := logo.NewFileHandler("app.log", formatter,
    logo.WithRotation(rotationConfig),
    logo.WithFileLevel(logo.LevelInfo),
)

router := logo.NewRouter(fileHandler)
logger := logo.New(logo.WithRouter(router))
```

## Rotation Parameters

### MaxSize

Maximum size in megabytes before rotation occurs.

```go
rotationConfig := &logo.RotationConfig{
    MaxSize: 50,  // Rotate at 50 MB
}
```

- **Type**: `int`
- **Unit**: Megabytes
- **Default**: 100
- **Minimum**: 1

### MaxAge

Maximum number of days to retain old log files.

```go
rotationConfig := &logo.RotationConfig{
    MaxAge: 14,  // Keep logs for 14 days
}
```

- **Type**: `int`
- **Unit**: Days
- **Default**: 30
- **Behavior**: Files older than this are deleted during rotation

### MaxBackups

Maximum number of old log files to retain.

```go
rotationConfig := &logo.RotationConfig{
    MaxBackups: 3,  // Keep 3 old files
}
```

- **Type**: `int`
- **Unit**: Count
- **Default**: 10
- **Behavior**: Oldest files are deleted when limit is reached

### Compress

Whether to compress rotated files using gzip.

```go
rotationConfig := &logo.RotationConfig{
    Compress: true,  // Compress old files
}
```

- **Type**: `bool`
- **Default**: true
- **Behavior**: Rotated files are gzipped (`.log.gz`)
- **Savings**: Typically 80-90% size reduction

### LocalTime

Use local time instead of UTC for backup timestamps.

```go
rotationConfig := &logo.RotationConfig{
    LocalTime: true,  // Use local time
}
```

- **Type**: `bool`
- **Default**: false (UTC)
- **Behavior**: Affects timestamp in rotated file names

## File Naming

### Current Log File

```
app.log
```

The active log file always has the specified name.

### Rotated Files (UTC)

```
app-2026-02-07T15-04-05.123.log.gz
```

Format: `<name>-<timestamp>.log[.gz]`
- Timestamp format: `YYYY-MM-DDTHH-MM-SS.mmm`
- `.gz` extension added if compressed

### Rotated Files (Local Time)

```
app-2026-02-07T10-04-05.123.log.gz
```

Same format, but timestamp reflects local timezone.

## Rotation Behavior

### When Rotation Occurs

Rotation is triggered when:
1. Current file size exceeds `MaxSize`
2. A write operation is attempted

The rotation process:
1. Close current file
2. Rename to timestamped name
3. Compress if `Compress: true`
4. Open new file with original name
5. Delete old files exceeding `MaxAge` or `MaxBackups`

### Retention Policy

Files are deleted when they meet either condition:
- **Age**: Older than `MaxAge` days
- **Count**: More than `MaxBackups` files exist

Example with `MaxAge: 7` and `MaxBackups: 5`:
```
app.log                              # Current (5 MB)
app-2026-02-07T10-00-00.123.log.gz  # 1 day old
app-2026-02-06T10-00-00.123.log.gz  # 2 days old
app-2026-02-05T10-00-00.123.log.gz  # 3 days old
app-2026-02-04T10-00-00.123.log.gz  # 4 days old
app-2026-02-03T10-00-00.123.log.gz  # 5 days old
# Older files deleted (> MaxBackups or > MaxAge)
```

## Manual Rotation

Trigger rotation programmatically:

```go
fileHandler, _ := logo.NewFileHandler("app.log", formatter)

// Force rotation
err := fileHandler.Rotate()
if err != nil {
    // Handle error
}
```

Use cases:
- Logrotate integration
- Scheduled maintenance
- Testing

## CLI Configuration

Configure rotation via command-line flags:

```bash
./app -log.file.path=app.log \
      -log.rotation.maxsize=50 \
      -log.rotation.maxage=14 \
      -log.rotation.maxbackups=7 \
      -log.rotation.compress=true
```

## Examples

### High-Volume Application

```go
rotationConfig := &logo.RotationConfig{
    MaxSize:    10,   // Rotate every 10 MB
    MaxAge:     3,    // Keep 3 days
    MaxBackups: 100,  // Keep up to 100 files
    Compress:   true,
}
```

### Long-Term Retention

```go
rotationConfig := &logo.RotationConfig{
    MaxSize:    100,  // Large files
    MaxAge:     365,  // Keep 1 year
    MaxBackups: 50,   // Limit to 50 files
    Compress:   true,
}
```

### Development (No Rotation)

```go
rotationConfig := &logo.RotationConfig{
    MaxSize:    1000, // Very large (1 GB)
    MaxAge:     0,    // Never delete
    MaxBackups: 0,    // No limit
    Compress:   false,
}
```

## Integration with External Tools

### Logrotate

Use `copytruncate` if integrating with system logrotate:

```
/var/log/myapp/*.log {
    daily
    rotate 7
    compress
    copytruncate
    missingok
}
```

Or trigger manual rotation:

```go
// SIGUSR1 handler
signal.Notify(sigChan, syscall.SIGUSR1)
go func() {
    for range sigChan {
        fileHandler.Rotate()
    }
}()
```

### Container Logs

In containers, consider writing to stdout and letting the container runtime handle rotation:

```go
log := logo.New(logo.WithConsole())  // No file handler
```

Container runtime (Docker, Kubernetes) manages log rotation.

## Best Practices

1. **Set Appropriate MaxSize** - Balance between file size and rotation frequency
   - High-traffic: 10-50 MB
   - Low-traffic: 100-500 MB

2. **Enable Compression** - Saves disk space (typically 80-90% reduction)

3. **Use UTC Timestamps** - Easier to correlate logs across timezones

4. **Monitor Disk Space** - Set `MaxAge` and `MaxBackups` based on available storage

5. **Test Rotation** - Verify rotation works before production deployment

6. **Consider Volume** - For very high-volume apps, use external log aggregation

7. **Backup Logs** - Archive important logs before they're deleted

## Troubleshooting

**Rotation not occurring:**
- Check `MaxSize` is set (minimum 1 MB)
- Verify file handler is properly initialized
- Ensure writes are actually occurring

**Out of disk space:**
- Reduce `MaxAge` or `MaxBackups`
- Enable compression
- Increase rotation frequency (lower `MaxSize`)

**Permission errors:**
- Check file and directory permissions
- Ensure write access to log directory
- Verify user has permission to create files

**Old files not being deleted:**
- Verify `MaxAge` and `MaxBackups` are set
- Check file timestamps match expected format
- Ensure rotation is occurring (triggers cleanup)

## lumberjack Integration

go-logo uses [lumberjack](https://github.com/aN0mad/lumberjack/v2) for rotation:

```go
import "github.com/aN0mad/lumberjack/v2"

// FileHandler internally uses:
writer := &lumberjack.Logger{
    Filename:   path,
    MaxSize:    config.MaxSize,
    MaxAge:     config.MaxAge,
    MaxBackups: config.MaxBackups,
    Compress:   config.Compress,
    LocalTime:  config.LocalTime,
}
```

For advanced lumberjack features, create handlers manually with custom lumberjack configurations.

