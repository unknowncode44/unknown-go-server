package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/unknowncode44/unknown-go-server/api/handlers"
)

// CurrencyRoutes registers currency routes under the provided router.
func CurrencyRoutes(router fiber.Router, handler *handlers.CurrencyHandler) {
	currencies := router.Group("/currencies")

	currencies.Post("/", handler.Create)
	currencies.Get("/", handler.GetAll)
	currencies.Get(":id", handler.GetById)
	currencies.Patch(":id", handler.Update)
	currencies.Post(":id/deactivate", handler.Deactivate)
}
