package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/unknowncode44/unknown-go-server/api/handlers"
)

func NeedRoutes(router fiber.Router, handler *handlers.NeedHandler) {
	n := router.Group("/needs")

	n.Post("/", handler.Create)
	n.Get("/", handler.GetAll)
	n.Get("/:id", handler.GetByID)
	n.Put("/:id", handler.Update)

	// State transition actions
	n.Post("/:id/promote", handler.Promote)
	n.Post("/:id/mark-converted", handler.MarkConverted)
	n.Post("/:id/discard", handler.Discard)

	// Items
	n.Get("/:id/items", handler.ListItems)
	n.Post("/:id/items", handler.AddItem)
	n.Put("/items/:itemId", handler.UpdateItem)
	n.Delete("/items/:itemId", handler.RemoveItem)
}
