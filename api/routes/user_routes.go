package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/unknowncode44/unknown-go-server/api/handlers"
	"github.com/unknowncode44/unknown-go-server/pkg/auth"
	"github.com/unknowncode44/unknown-go-server/pkg/entities"
)

// UserRoutes registra las rutas de gestión de usuarios.
// Todas requieren estar logueado con rol ADMIN.
func UserRoutes(router fiber.Router, handler *handlers.UserHandler) {
	users := router.Group("/users", auth.RequireAuth(), auth.RequireRole(entities.UserRoleAdmin))

	users.Post("/", handler.Create)
	users.Get("/", handler.FindAll)
	users.Get("/:id", handler.FindByID)
	users.Put("/:id", handler.Update)
	users.Put("/:id/password", handler.SetPassword)
	users.Delete("/:id", handler.Deactivate)
}
