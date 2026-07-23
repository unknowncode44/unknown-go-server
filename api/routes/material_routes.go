package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/unknowncode44/unknown-go-server/api/handlers"
)

func MaterialRoutes(router fiber.Router, handler *handlers.MaterialHandler) {
	materials := router.Group("/materials")

	materials.Post("/", handler.Create)
	materials.Get("/", handler.GetAll)
	materials.Get("/:id", handler.GetById)
	materials.Get("/by-erp/:erp_code", handler.GetMaterialsByERP)
	materials.Put("/:id", handler.Update)
	materials.Delete("/:id", handler.Deactivate)
}

// Public route for unauthenticated material lookup (bulk/BDC materials)
func PublicMaterialRoute(router fiber.Router, handler *handlers.MaterialHandler) {
	router.Get("/material/:code", handler.PublicByCode)
}
