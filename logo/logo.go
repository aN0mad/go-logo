package logger

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	cblog "github.com/charmbracelet/log"
)

type LogLevel string

const (
	LevelDebug LogLevel = "DEBUG"
	LevelInfo  LogLevel = "INFO"
	LevelWarn  LogLevel = "WARN"
	LevelError LogLevel = "ERROR"
)

type LogMessage struct {
	Time    time.Time `json:"time"`
	Level   LogLevel  `json:"level"`
	Message string    `json:"message"`
}

type Logger struct {
	mu             sync.Mutex
	logChan        chan LogMessage
	consoleEnabled bool
	console        *cblog.Logger
}

// New creates a new Logger.
// logChan: if non-nil, logs will be sent to channel.
// consoleEnabled: if true, logs will also be printed using charmbracelet/log.
func New(logChan chan LogMessage, consoleEnabled bool) *Logger {
	output := io.Discard
	if consoleEnabled {
		output = os.Stderr
	}

	console := cblog.NewWithOptions(output, cblog.Options{
		ReportCaller: false,
		Prefix:       "lib",
		Level:        cblog.DebugLevel,
	})

	return &Logger{
		logChan:        logChan,
		consoleEnabled: consoleEnabled,
		console:        console,
	}
}

func (l *Logger) Debug(format string, args ...any) { l.log(LevelDebug, format, args...) }
func (l *Logger) Info(format string, args ...any)  { l.log(LevelInfo, format, args...) }
func (l *Logger) Warn(format string, args ...any)  { l.log(LevelWarn, format, args...) }
func (l *Logger) Error(format string, args ...any) { l.log(LevelError, format, args...) }

func (l *Logger) log(level LogLevel, format string, args ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()

	msg := fmt.Sprintf(format, args...)
	entry := LogMessage{
		Time:    time.Now(),
		Level:   level,
		Message: msg,
	}

	// Send to channel
	if l.logChan != nil {
		select {
		case l.logChan <- entry:
		default:
			// Optional: log drop warning
		}
	}

	// Send to console
	if l.consoleEnabled {
		switch level {
		case LevelDebug:
			l.console.Debug(msg)
		case LevelInfo:
			l.console.Info(msg)
		case LevelWarn:
			l.console.Warn(msg)
		case LevelError:
			l.console.Error(msg)
		default:
			l.console.Print(msg)
		}
	}
}
