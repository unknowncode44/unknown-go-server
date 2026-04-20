package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/unknowncode44/unknown-go-server/api/handlers"
)

func DeliveryRecordRoutes(router fiber.Router, handler *handlers.DeliveryRecordHandler) {
	dr := router.Group("/delivery-records")
	dr.Post("/", handler.Create)
	dr.Get("/", handler.GetAll)
}

func DeliveryRecordByInventoryRoutes(router fiber.Router, handler *handlers.MaterialInventoryHandler) {
	inv := router.Group("/material-inventories")
	inv.Get("/:id/deliveries", handler.FindDeliveriesByInventoryID)
}
