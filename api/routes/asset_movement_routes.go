package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/unknowncode44/unknown-go-server/api/handlers"
)

func AssetMovementRoutes(router fiber.Router, handler *handlers.AssetMovementHandler) {
	ams := router.Group("/asset-movements")

	ams.Post("/", handler.Create)
	// optionally could add GET /asset-movements/:id in future
}

// Asset-scoped movement endpoints
func AssetMovementByAssetRoutes(router fiber.Router, handler *handlers.AssetMovementHandler) {
	assets := router.Group("/assets")
	assets.Get("/:asset_id/movements", handler.FindByAsset)
	assets.Get("/:asset_id/movements/last", handler.FindLastByAsset)
}
