package main

import (
	"net/http"
	"time"

	"github.com/aN0mad/go-logo/v2"
	"github.com/gin-gonic/gin"
)

// loggerMiddleware demonstrates how to integrate go-logo with Gin.
// Users can customize this middleware to fit their needs.
func loggerMiddleware(logger *logo.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Generate or extract request ID
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = logo.GenerateRequestID()
		}
		c.Header("X-Request-ID", requestID)

		// Add request ID to context for context-aware logging
		ctx := logo.WithRequestID(c.Request.Context(), requestID)
		c.Request = c.Request.WithContext(ctx)

		// Log request
		logger.With(
			logo.F("request_id", requestID),
			logo.F("method", c.Request.Method),
			logo.F("path", c.Request.URL.Path),
			logo.F("ip", c.ClientIP()),
		).Info("request started")

		// Process request
		c.Next()

		// Log response
		duration := time.Since(start)
		logger.With(
			logo.F("request_id", requestID),
			logo.F("status", c.Writer.Status()),
			logo.F("duration_ms", duration.Milliseconds()),
		).Info("request completed")
	}
}

func main() {
	// Create logger with CLI configuration
	cfg := logo.NewConfig()
	cfg.RegisterFlags(nil)

	logger, err := cfg.BuildLogger()
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	// Set Gin to release mode (or use gin.SetMode(gin.DebugMode) for development)
	gin.SetMode(gin.ReleaseMode)

	// Create Gin router
	router := gin.New()

	// Use custom logger middleware
	router.Use(loggerMiddleware(logger))

	// Recovery middleware
	router.Use(gin.Recovery())

	// Routes
	router.GET("/", func(c *gin.Context) {
		// Get request ID from response header
		requestID := c.Writer.Header().Get("X-Request-ID")

		// Use context-aware logging
		log := logo.FromContext(c.Request.Context())
		log.Info("handling root request")

		c.JSON(http.StatusOK, gin.H{
			"message":    "Hello from Gin!",
			"request_id": requestID,
		})
	})

	router.GET("/user/:id", func(c *gin.Context) {
		// Context-aware logging with additional fields
		log := logo.FromContext(c.Request.Context())
		log.With(
			logo.F("user_id", c.Param("id")),
			logo.F("action", "get_user"),
		).Info("fetching user data")

		c.JSON(http.StatusOK, gin.H{
			"user_id": c.Param("id"),
			"name":    "Alice Smith",
		})
	})

	router.POST("/orders", func(c *gin.Context) {
		log := logo.FromContext(c.Request.Context())
		log.With(
			logo.F("items", 5),
		).Info("creating order")

		c.JSON(http.StatusCreated, gin.H{
			"order_id": "ORD-99999",
			"status":   "created",
		})
	})

	logger.Info("starting Gin server", logo.F("port", 3001))
	if err := router.Run(":3001"); err != nil {
		logger.Fatal("server failed", logo.F("error", err.Error()))
	}
}
