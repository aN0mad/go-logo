package main

import (
	"net/http"
	"time"

	"github.com/aN0mad/go-logo/v2"
	"github.com/labstack/echo/v4"
)

// loggerMiddleware demonstrates how to integrate go-logo with Echo.
// Users can customize this middleware to fit their needs.
func loggerMiddleware(logger *logo.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			// Generate request ID
			requestID := c.Request().Header.Get("X-Request-ID")
			if requestID == "" {
				requestID = logo.GenerateRequestID()
			}
			c.Response().Header().Set("X-Request-ID", requestID)

			// Add request ID to context for context-aware logging
			req := c.Request()
			ctx := logo.WithRequestID(req.Context(), requestID)
			c.SetRequest(req.WithContext(ctx))

			// Log request
			logger.With(
				logo.F("request_id", requestID),
				logo.F("method", c.Request().Method),
				logo.F("path", c.Request().URL.Path),
				logo.F("ip", c.RealIP()),
			).Info("request started")

			// Process request
			err := next(c)

			// Log response
			duration := time.Since(start)
			logger.With(
				logo.F("request_id", requestID),
				logo.F("status", c.Response().Status),
				logo.F("duration_ms", duration.Milliseconds()),
			).Info("request completed")

			return err
		}
	}
}

func main() {
	// Create logger
	cfg := logo.NewConfig()
	cfg.RegisterFlags(nil)

	logger, err := cfg.BuildLogger()
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	// Create Echo instance
	e := echo.New()
	e.HideBanner = true

	// Use custom logger middleware
	e.Use(loggerMiddleware(logger))

	// Routes
	e.GET("/", func(c echo.Context) error {
		// Get request ID from response header
		requestID := c.Response().Header().Get("X-Request-ID")

		// Use context-aware logging
		log := logo.FromContext(c.Request().Context())
		log.Info("handling root request")

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message":    "Hello from Echo!",
			"request_id": requestID,
		})
	})

	e.GET("/user/:id", func(c echo.Context) error {
		// Context-aware logging with additional fields
		log := logo.FromContext(c.Request().Context())
		log.With(
			logo.F("user_id", c.Param("id")),
			logo.F("action", "get_user"),
		).Info("fetching user data")

		return c.JSON(http.StatusOK, map[string]interface{}{
			"user_id": c.Param("id"),
			"name":    "Bob Johnson",
		})
	})

	e.POST("/orders", func(c echo.Context) error {
		log := logo.FromContext(c.Request().Context())
		log.With(
			logo.F("items", 7),
		).Info("creating order")

		return c.JSON(http.StatusCreated, map[string]interface{}{
			"order_id": "ORD-11111",
			"status":   "created",
		})
	})

	logger.Info("starting Echo server", logo.F("port", 3002))
	if err := e.Start(":3002"); err != nil {
		logger.Fatal("server failed", logo.F("error", err.Error()))
	}
}
