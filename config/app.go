package config

import (
	"errors"
	"log/slog"

	"findit-campus-backend/helper"
	"findit-campus-backend/middleware"

	"github.com/gofiber/fiber/v2"
)

func NewApp(cfg *Config, logger *slog.Logger) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName: cfg.AppName,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			msg := "internal server error"

			var fiberErr *fiber.Error
			if errors.As(err, &fiberErr) {
				code = fiberErr.Code
				msg = fiberErr.Message
			}

			return helper.Fail(c, code, msg)
		},
	})

	middleware.Register(app, logger)
	return app
}