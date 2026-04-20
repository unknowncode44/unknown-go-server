package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/unknowncode44/unknown-go-server/api/handlers"
)

func MaterialInventoryRoutes(router fiber.Router, handler *handlers.MaterialInventoryHandler) {
	inv := router.Group("/material-inventories")
	inv.Post("/", handler.Create)
	inv.Get("/", handler.GetAll)
	inv.Get("/:id", handler.GetByID)
	inv.Get("/:id/total", handler.GetTotal)
	inv.Put("/:id", handler.Update)
	inv.Delete("/:id", handler.Deactivate)
}

func MaterialInventoryByMaterialRoutes(router fiber.Router, handler *handlers.MaterialInventoryHandler) {
	materials := router.Group("/materials")
	materials.Get("/:id/inventory", handler.GetByMaterialID)
}
