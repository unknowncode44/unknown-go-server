package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/unknowncode44/unknown-go-server/api/handlers"
)

func MaterialCostRoutes(router fiber.Router, handler *handlers.MaterialCostHandler) {
	mcs := router.Group("/material-costs")

	mcs.Post("/", handler.Create)
	mcs.Get("/", handler.FindAll)
	mcs.Get("/:id", handler.FindByID)

	// route to get costs for a given material
	materials := router.Group("/materials")
	materials.Get("/:id/costs", handler.FindByMaterial)
	materials.Get("/:materialId/vendors/:vendorId/current-cost", handler.GetCurrentCost)
	materials.Get("/:materialId/vendor-comparison", handler.CompareVendorsByMaterial)
}
