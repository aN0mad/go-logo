package main

import (
"os"
"time"

"github.com/aN0mad/go-logo/v2"
)

func main() {
	jsonHandler := logo.NewJSONHandler(os.Stdout, logo.LevelDebug)
	router := logo.NewRouter(jsonHandler)
	log := logo.New(logo.WithRouter(router))

	log.Info("application started", logo.F("version", "2.0.0"), logo.F("env", "production"))
	log.Debug("debug information", logo.F("detail", "low-level"))
	log.Warn("resource usage high", logo.F("cpu_percent", 85.5))
	log.Error("operation failed", logo.F("error", "connection timeout"), logo.F("retry_count", 3))

	log.Info("request processed",
logo.F("method", "POST"),
logo.F("path", "/api/users"),
logo.F("status", 201),
logo.F("duration", 45*time.Millisecond),
)

	requestLog := log.With(logo.F("request_id", "req-abc-123"), logo.F("user_id", 12345))
	requestLog.Info("processing user request")
	requestLog.Info("database query complete", logo.F("rows", 42))

	log.Info("application shutdown")
}
