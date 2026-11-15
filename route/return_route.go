package route

import (
	"findit-campus-backend/app/service"
	"findit-campus-backend/config"
	"findit-campus-backend/middleware"

	"github.com/gofiber/fiber/v2"
)

// registerReturnRoutes memetakan endpoint serah terima barang & riwayat status laporan (PIC 4).
func registerReturnRoutes(api fiber.Router, cfg *config.Config, svc *service.ReturnHistoryService) {
	protected := middleware.Protected(cfg.JWTSecret)

	returns := api.Group("/returns")
	returns.Post("/", protected, middleware.RequireRole("security", "admin"), middleware.RequireJSON, svc.CreateReturn)
	returns.Get("/:report_id", protected, svc.GetReturnByReport)

	api.Get("/reports/:id/histories", protected, svc.ListHistories)
}
