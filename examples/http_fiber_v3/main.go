package main

import (
	"time"

	"github.com/aN0mad/go-logo/v2"
	"github.com/gofiber/fiber/v3"
)

// loggerMiddleware demonstrates how to integrate go-logo with Fiber v3.
// Users can customize this middleware to fit their needs.
func loggerMiddleware(logger *logo.Logger) fiber.Handler {
	return func(c fiber.Ctx) error {
		start := time.Now()

		// Generate request ID
		requestID := logo.GenerateRequestID()
		c.Locals("request_id", requestID)

		// Add request ID to context for context-aware logging
		ctx := logo.WithRequestID(c.Context(), requestID)
		c.SetUserContext(ctx)

		// Log request
		logger.With(
			logo.F("request_id", requestID),
			logo.F("method", c.Method()),
			logo.F("path", c.Path()),
			logo.F("ip", c.IP()),
		).Info("request started")

		// Process request
		err := c.Next()

		// Log response
		duration := time.Since(start)
		logger.With(
			logo.F("request_id", requestID),
			logo.F("status", c.Response().StatusCode()),
			logo.F("duration_ms", duration.Milliseconds()),
		).Info("request completed")

		return err
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

	// Create Fiber v3 app
	app := fiber.New()

	// Use custom logger middleware
	app.Use(loggerMiddleware(logger))

	// Routes
	app.Get("/", func(c fiber.Ctx) error {
		// Get request ID from context
		requestID := c.Locals("request_id").(string)

		// Use context-aware logging
		log := logo.FromContext(c.Context())
		log.Info("handling root request")

		return c.JSON(fiber.Map{
			"message":    "Hello from Fiber v3!",
			"request_id": requestID,
		})
	})

	app.Get("/user/:id", func(c fiber.Ctx) error {
		// Context-aware logging with additional fields
		log := logo.FromContext(c.Context())
		log.With(
			logo.F("user_id", c.Params("id")),
			logo.F("action", "get_user"),
		).Info("fetching user data")

		return c.JSON(fiber.Map{
			"user_id": c.Params("id"),
			"name":    "Jane Doe",
		})
	})

	app.Post("/orders", func(c fiber.Ctx) error {
		log := logo.FromContext(c.Context())
		log.With(
			logo.F("items", 3),
		).Info("creating order")

		return c.Status(201).JSON(fiber.Map{
			"order_id": "ORD-54321",
			"status":   "created",
		})
	})

	logger.Info("starting Fiber v3 server", logo.F("port", 3003))
	if err := app.Listen(":3003"); err != nil {
		logger.Fatal("server failed", logo.F("error", err.Error()))
	}
}
