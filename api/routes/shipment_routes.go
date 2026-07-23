package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/unknowncode44/unknown-go-server/api/handlers"
	"github.com/unknowncode44/unknown-go-server/pkg/auth"
)

// ShipmentRoutes registers the authenticated shipment endpoints.
func ShipmentRoutes(router fiber.Router, handler *handlers.ShipmentHandler) {
	shipments := router.Group("/shipments", auth.RequireAuth())

	shipments.Post("/", handler.Create)
	shipments.Get("/", handler.FindAll)
	shipments.Get("/:id", handler.FindByID)
	shipments.Post("/:id/items", handler.AddItem)
	shipments.Delete("/:id/items/:itemId", handler.RemoveItem)
	shipments.Post("/:id/process", handler.Process)
	shipments.Put("/:id/erp-link", handler.LinkERP)
	shipments.Delete("/:id", handler.Cancel)
}
