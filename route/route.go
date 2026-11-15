package route

import (
	"findit-campus-backend/app/service"
	"findit-campus-backend/config"
	"findit-campus-backend/helper"
	"findit-campus-backend/middleware"

	"github.com/gofiber/fiber/v2"
)

type Services struct {
	Auth   *service.AuthService
	Report *service.ReportService
  Claim *service.ClaimService
	ReturnHistory *service.ReturnHistoryService
}

func Register(app *fiber.App, cfg *config.Config, svc Services) {
	app.Get("/health", func(c *fiber.Ctx) error {
		return helper.Success(c, "FindIt Campus API is running", fiber.Map{
			"status": "ok",
		})
	})

	api := app.Group("/api/v1")

	// Auth routes 
	auth := api.Group("/auth")
	auth.Post("/register", middleware.RequireJSON, svc.Auth.Register)
	auth.Post("/login", middleware.RequireJSON, svc.Auth.Login)
	auth.Get("/me", middleware.Protected(cfg.JWTSecret), svc.Auth.Me)

	// Categories, Locations, Reports 
	if svc.Report != nil {
		RegisterReportRoutes(api, cfg, svc)
	}
  // Claims
  RegisterClaimRoutes(api, cfg, svc.Claim)
	// Return & History routes
	registerReturnRoutes(api, cfg, svc.ReturnHistory)
}