package handlers

import (
	"github.com/gofiber/fiber/v2"
)

// respondError centralizes error responses with a uniform JSON shape.
func respondError(c *fiber.Ctx, status int, msg string) error {
	return c.Status(status).JSON(fiber.Map{"success": false, "error": msg})
}
