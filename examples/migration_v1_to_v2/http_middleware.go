package main

// This example demonstrates v2's HTTP middleware integration
// which has no equivalent in v1.

import (
	logo "github.com/aN0mad/go-logo/v2"
	fiberlogo "github.com/aN0mad/go-logo/v2/middleware/fiber"
	"github.com/gofiber/fiber/v2"
)

func main() {
	// Create logger
	logger := logo.New(
		logo.WithConsole(),
		logo.WithLevel(logo.LevelInfo),
	)
	defer logger.Close()

	// Create Fiber app
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	// v2 NEW: HTTP middleware with request-scoped logging
	app.Use(fiberlogo.Middleware(logger))

	// Routes
	app.Get("/users/:id", getUser)
	app.Post("/orders", createOrder)
	app.Get("/health", healthCheck)

	logger.Info("Starting HTTP server", logo.F("port", 3000))
	app.Listen(":3000")
}

func getUser(c *fiber.Ctx) error {
	// v2 NEW: Get request-scoped logger
	log := fiberlogo.Logger(c)
	requestID := fiberlogo.RequestID(c)

	userID := c.Params("id")

	// All logs include request ID automatically
	log.Info("Fetching user", logo.F("user_id", userID))

	// Simulate work
	user := map[string]interface{}{
		"id":   userID,
		"name": "John Doe",
	}

	log.Info("User fetched successfully", logo.F("user_id", userID))

	return c.JSON(fiber.Map{
		"request_id": requestID,
		"user":       user,
	})
}

func createOrder(c *fiber.Ctx) error {
	log := fiberlogo.Logger(c)

	// v2 NEW: Add fields to request context
	log = logo.ContextWithFields(c.Context(),
		logo.F("action", "create_order"),
		logo.F("user_id", 123),
	)

	log.Info("Creating order")

	// Simulate validation
	log.Debug("Validating order data")

	// Simulate saving
	log.Info("Order created", logo.F("order_id", "ORD-12345"))

	return c.Status(201).JSON(fiber.Map{
		"order_id": "ORD-12345",
		"status":   "created",
	})
}

func healthCheck(c *fiber.Ctx) error {
	// Health checks typically don't need verbose logging
	return c.SendString("OK")
}

/*
Run this example and test:

  go run http_middleware.go

Then in another terminal:

  curl http://localhost:3000/users/123
  curl -X POST http://localhost:3000/orders
  curl http://localhost:3000/health

Watch the server logs show:
- Automatic request logging (method, path, status, duration)
- Request IDs for correlation
- Request-scoped context fields
- All logs tied to specific requests
*/
