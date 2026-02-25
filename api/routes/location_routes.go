package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/unknowncode44/unknown-go-server/api/handlers"
)

func LocationRoutes(router fiber.Router, handler *handlers.LocationHandler) {
	locations := router.Group("/locations")

	locations.Post("/", handler.Create)
	locations.Get("/", handler.GetAll)
	locations.Get("/:id", handler.GetByID)
	locations.Put("/:id", handler.Update)
	locations.Delete("/:id", handler.Deactivate)
}
