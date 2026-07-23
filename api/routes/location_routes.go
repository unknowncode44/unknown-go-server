package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/unknowncode44/unknown-go-server/api/handlers"
	"github.com/unknowncode44/unknown-go-server/pkg/auth"
	"github.com/unknowncode44/unknown-go-server/pkg/entities"
)

// LocationRoutes registra las rutas de ubicaciones. La gestión de
// depósitos/sectores/estanterías es exclusiva del rol ADMIN.
func LocationRoutes(router fiber.Router, handler *handlers.LocationHandler) {
	locations := router.Group("/locations", auth.RequireAuth(), auth.RequireRole(entities.UserRoleAdmin))

	locations.Post("/", handler.Create)
	locations.Get("/", handler.GetAll)
	locations.Get("/:id", handler.GetByID)
	locations.Put("/:id", handler.Update)
	locations.Delete("/:id", handler.Deactivate)
}

// Public route for unauthenticated location lookup (shelf/estantería QR)
func PublicLocationRoute(router fiber.Router, handler *handlers.LocationHandler) {
	router.Get("/location/:id", handler.PublicByID)
}
