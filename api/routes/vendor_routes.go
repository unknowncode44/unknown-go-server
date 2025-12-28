package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/unknowncode44/unknown-go-server/api/handlers"
)

func VendorRoutes(router fiber.Router, handler *handlers.VendorHandler) {
	vendors := router.Group("/vendors")

	vendors.Post("/", handler.Create)
	vendors.Get("/", handler.FindAll)
	vendors.Get("/:id", handler.FindByID)
	vendors.Put("/:id", handler.Update)
	vendors.Delete("/:id", handler.Deactivate)
}
