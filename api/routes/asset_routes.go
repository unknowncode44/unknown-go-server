package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/unknowncode44/unknown-go-server/api/handlers"
)

func AssetRoutes(router fiber.Router, handler *handlers.AssetHandler) {
	assets := router.Group("/assets")

	assets.Post("/", handler.Create)
	assets.Get("/", handler.FindAll)
	assets.Get("/:id", handler.FindByID)
	assets.Get("/serial/:serial", handler.FindBySerial)
	assets.Put("/:id", handler.Update)
	assets.Delete("/:id", handler.Deactivate)
}

// Public route for unauthenticated asset lookup
func PublicAssetRoute(router fiber.Router, handler *handlers.AssetHandler) {
	router.Get("/asset/:serial", handler.PublicBySerial)
}
