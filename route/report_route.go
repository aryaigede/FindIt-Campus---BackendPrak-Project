package route

import (
	"findit-campus-backend/config"
	"findit-campus-backend/middleware"

	"github.com/gofiber/fiber/v2"
)

// RegisterReportRoutes registers all routes handled by PIC 2 (categories, locations, reports, images).
func RegisterReportRoutes(router fiber.Router, cfg *config.Config, svc Services) {
	adminAuth := []fiber.Handler{
		middleware.Protected(cfg.JWTSecret),
		middleware.RequireRole("admin"),
	}

	// Categories
	categories := router.Group("/categories")
	categories.Get("/", svc.Report.ListCategories)
	categories.Get("", svc.Report.ListCategories)
	categories.Get("/:id", svc.Report.GetCategory)
	categories.Post("/", append(adminAuth, middleware.RequireJSON, svc.Report.CreateCategory)...)
	categories.Post("", append(adminAuth, middleware.RequireJSON, svc.Report.CreateCategory)...)
	categories.Put("/:id", append(adminAuth, middleware.RequireJSON, svc.Report.UpdateCategory)...)
	categories.Delete("/:id", append(adminAuth, svc.Report.DeleteCategory)...)

	// Locations
	locations := router.Group("/locations")
	locations.Get("/", svc.Report.ListLocations)
	locations.Get("", svc.Report.ListLocations)
	locations.Get("/:id", svc.Report.GetLocation)
	locations.Post("/", append(adminAuth, middleware.RequireJSON, svc.Report.CreateLocation)...)
	locations.Post("", append(adminAuth, middleware.RequireJSON, svc.Report.CreateLocation)...)
	locations.Put("/:id", append(adminAuth, middleware.RequireJSON, svc.Report.UpdateLocation)...)
	locations.Delete("/:id", append(adminAuth, svc.Report.DeleteLocation)...)
}
