package middleware

import (
	"log/slog"
	"strings"
	"time"

	"findit-campus-backend/helper"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func Register(app *fiber.App, logger *slog.Logger) {
	app.Use(requestid.New())
	app.Use(recover.New())
	app.Use(helmet.New())
	app.Use(cors.New())
	app.Use(RequestLogger(logger))
}

func RequestLogger(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		duration := time.Since(start)

		logger.Info("http_request",
			"request_id", c.Locals("requestid"),
			"method", c.Method(),
			"path", c.Path(),
			"status", c.Response().StatusCode(),
			"latency_ms", duration.Milliseconds(),
		)
		return err
	}
}

func RequireJSON(c *fiber.Ctx) error {
	contentType := c.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		return helper.Fail(c, fiber.StatusUnsupportedMediaType, "Content-Type must be application/json")
	}
	return c.Next()
}