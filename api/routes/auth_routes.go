package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/unknowncode44/unknown-go-server/api/handlers"
	"github.com/unknowncode44/unknown-go-server/pkg/auth"
)

// AuthRoutes registra las rutas de autenticación. /login es pública (es la
// puerta de entrada); /me exige un token válido.
func AuthRoutes(router fiber.Router, handler *handlers.AuthHandler) {
	authGroup := router.Group("/auth")

	authGroup.Post("/login", handler.Login)              // público, sin RequireAuth
	authGroup.Get("/me", auth.RequireAuth(), handler.Me) // requiere login
}
