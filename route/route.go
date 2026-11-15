package route

import (
	"findit-campus-backend/app/service"
	"findit-campus-backend/config"
	"findit-campus-backend/helper"
	"findit-campus-backend/middleware"

	"github.com/gofiber/fiber/v2"
)

type Services struct {
	Auth  *service.AuthService
	Claim *service.ClaimService
	// Anggota tim (PIC 2, 3, 4) menambahkan field service di sini saat PR fitur masing-masing
}

func Register(app *fiber.App, cfg *config.Config, svc Services) {
	app.Get("/health", func(c *fiber.Ctx) error {
		return helper.Success(c, "FindIt Campus API is running", fiber.Map{
			"status": "ok",
		})
	})

	api := app.Group("/api/v1")

	// Auth routes (Utuh - PIC 1)
	auth := api.Group("/auth")
	auth.Post("/register", middleware.RequireJSON, svc.Auth.Register)
	auth.Post("/login", middleware.RequireJSON, svc.Auth.Login)
	auth.Get("/me", middleware.Protected(cfg.JWTSecret), svc.Auth.Me)

	// TODO (PIC 2): /categories, /locations, /reports
	RegisterClaimRoutes(api, cfg, svc.Claim)
	// TODO (PIC 4): /returns, /reports/:id/histories
}