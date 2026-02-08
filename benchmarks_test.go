package logo

import (
	"testing"
)

// Benchmark logger creation
func BenchmarkNew(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		cfg := NewConfig()
		cfg.Console = false
		logger, _ := cfg.BuildLogger()
		logger.Close()
	}
}

// Benchmark basic logging without fields
func BenchmarkInfo(b *testing.B) {
	cfg := NewConfig()
	cfg.Console = false
	logger, _ := cfg.BuildLogger()
	defer logger.Close()

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		logger.Info("benchmark message")
	}
}

// Benchmark logging with fields
func BenchmarkInfoWithFields(b *testing.B) {
	cfg := NewConfig()
	cfg.Console = false
	logger, _ := cfg.BuildLogger()
	defer logger.Close()

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		logger.Info("benchmark message",
			F("key1", "value1"),
			F("key2", 42),
			F("key3", true),
		)
	}
}

// Benchmark logging with many fields
func BenchmarkInfoWith10Fields(b *testing.B) {
	cfg := NewConfig()
	cfg.Console = false
	logger, _ := cfg.BuildLogger()
	defer logger.Close()

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		logger.Info("benchmark message",
			F("f1", "v1"), F("f2", "v2"), F("f3", "v3"), F("f4", "v4"), F("f5", "v5"),
			F("f6", "v6"), F("f7", "v7"), F("f8", "v8"), F("f9", "v9"), F("f10", "v10"),
		)
	}
}

// Benchmark JSON formatting
func BenchmarkJSONOutput(b *testing.B) {
	cfg := NewConfig()
	cfg.Console = false
	cfg.FileFormat = "json"
	cfg.FilePath = "/dev/null"
	logger, _ := cfg.BuildLogger()
	defer logger.Close()

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		logger.Info("benchmark message",
			F("user_id", 12345),
			F("action", "login"),
			F("ip", "192.168.1.1"),
		)
	}
}

// Benchmark console output
func BenchmarkConsoleOutput(b *testing.B) {
	cfg := NewConfig()
	cfg.Console = true
	cfg.ConsoleFormat = "text"
	logger, _ := cfg.BuildLogger()
	defer logger.Close()

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		logger.Info("benchmark message", F("key", "value"))
	}
}

// Benchmark level filtering
func BenchmarkLevelFiltering(b *testing.B) {
	cfg := NewConfig()
	cfg.Level = "warn"
	cfg.Console = false
	logger, _ := cfg.BuildLogger()
	defer logger.Close()

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		// These should be filtered out
		logger.Debug("debug message")
		logger.Info("info message")
		// This should pass
		logger.Warn("warn message")
	}
}

// Benchmark With() (child logger creation)
func BenchmarkWith(b *testing.B) {
	cfg := NewConfig()
	cfg.Console = false
	logger, _ := cfg.BuildLogger()
	defer logger.Close()

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		child := logger.With(F("request_id", "abc123"))
		child.Info("message")
	}
}

// Benchmark WithName() (named logger)
func BenchmarkWithName(b *testing.B) {
	cfg := NewConfig()
	cfg.Console = false
	logger, _ := cfg.BuildLogger()
	defer logger.Close()

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		named := logger.WithName("service")
		named.Info("message")
	}
}

// Benchmark concurrent logging
func BenchmarkConcurrent(b *testing.B) {
	cfg := NewConfig()
	cfg.Console = false
	logger, _ := cfg.BuildLogger()
	defer logger.Close()

	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("concurrent message",
				F("worker", "test"),
				F("count", 1),
			)
		}
	})
}

// Benchmark different log levels
func BenchmarkTrace(b *testing.B) {
	cfg := NewConfig()
	cfg.Level = "trace"
	cfg.Console = false
	logger, _ := cfg.BuildLogger()
	defer logger.Close()

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		logger.Trace("trace message")
	}
}

func BenchmarkDebug(b *testing.B) {
	cfg := NewConfig()
	cfg.Level = "debug"
	cfg.Console = false
	logger, _ := cfg.BuildLogger()
	defer logger.Close()

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		logger.Debug("debug message")
	}
}

func BenchmarkWarn(b *testing.B) {
	cfg := NewConfig()
	cfg.Level = "warn"
	cfg.Console = false
	logger, _ := cfg.BuildLogger()
	defer logger.Close()

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		logger.Warn("warn message")
	}
}

func BenchmarkError(b *testing.B) {
	cfg := NewConfig()
	cfg.Level = "error"
	cfg.Console = false
	logger, _ := cfg.BuildLogger()
	defer logger.Close()

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		logger.Error("error message")
	}
}

// Benchmark field types
func BenchmarkFieldString(b *testing.B) {
	cfg := NewConfig()
	cfg.Console = false
	logger, _ := cfg.BuildLogger()
	defer logger.Close()

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		logger.Info("msg", F("key", "string value"))
	}
}

func BenchmarkFieldInt(b *testing.B) {
	cfg := NewConfig()
	cfg.Console = false
	logger, _ := cfg.BuildLogger()
	defer logger.Close()

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		logger.Info("msg", F("key", 42))
	}
}

func BenchmarkFieldBool(b *testing.B) {
	cfg := NewConfig()
	cfg.Console = false
	logger, _ := cfg.BuildLogger()
	defer logger.Close()

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		logger.Info("msg", F("key", true))
	}
}

func BenchmarkFieldComplex(b *testing.B) {
	cfg := NewConfig()
	cfg.Console = false
	logger, _ := cfg.BuildLogger()
	defer logger.Close()

	data := map[string]interface{}{
		"user_id": 123,
		"name":    "test",
		"active":  true,
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		logger.Info("msg", F("data", data))
	}
}

// Request ID generator benchmarks are in requestid_test.go
