package route

import (
	"findit-campus-backend/app/service"
	"findit-campus-backend/config"
	"findit-campus-backend/middleware"

	"github.com/gofiber/fiber/v2"
)

// RegisterClaimRoutes registers all endpoints related to claims
func RegisterClaimRoutes(api fiber.Router, cfg *config.Config, svc *service.ClaimService) {
	if svc == nil {
		return
	}

	protected := middleware.Protected(cfg.JWTSecret)
	adminOrSecurity := middleware.RequireRole("security", "admin")

	// Claim submission under report resource
	api.Post("/reports/:id/claims", protected, middleware.RequireJSON, svc.SubmitClaim)

	// Claim management resources
	claims := api.Group("/claims", protected)
	claims.Get("/me", svc.GetMyClaims)
	claims.Get("", adminOrSecurity, svc.ListClaims)
	claims.Get("/", adminOrSecurity, svc.ListClaims)
	claims.Patch("/:id/review", adminOrSecurity, middleware.RequireJSON, svc.ReviewClaim)
}
