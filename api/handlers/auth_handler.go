package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/unknowncode44/unknown-go-server/api/presenter"
	"github.com/unknowncode44/unknown-go-server/pkg/auth"
	"github.com/unknowncode44/unknown-go-server/pkg/user"
)

// AuthHandler handles HTTP requests for authentication.
// It depends on a user.Service to verify credentials and on the auth
// package to sign JWTs.
type AuthHandler struct {
	userService user.Service
}

// NewAuthHandler creates a new AuthHandler with the provided user service.
func NewAuthHandler(userService user.Service) *AuthHandler {
	return &AuthHandler{userService: userService}
}

// Login handles POST /auth/login. Verifica email/password y devuelve un JWT
// junto con los datos públicos del usuario. Ante credenciales inválidas
// responde 401 con un error genérico (sin distinguir el motivo).
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req presenter.LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	u, err := h.userService.Authenticate(req.Email, req.Password)
	if err != nil {
		return respondError(c, fiber.StatusUnauthorized, err.Error())
	}

	token, err := auth.GenerateToken(u)
	if err != nil {
		return respondError(c, fiber.StatusInternalServerError, "no se pudo generar el token")
	}

	return c.JSON(presenter.UserSuccessResponse{
		Success: true,
		Data: presenter.LoginResponse{
			Token: token,
			User:  presenter.ToUserResponse(u),
		},
	})
}

// Me handles GET /auth/me. Devuelve los claims que RequireAuth dejó en
// fiber.Locals, sin volver a consultar la DB.
func (h *AuthHandler) Me(c *fiber.Ctx) error {
	return c.JSON(presenter.UserSuccessResponse{
		Success: true,
		Data: fiber.Map{
			"id":    c.Locals("user_id"),
			"email": c.Locals("user_email"),
			"role":  c.Locals("user_role"),
		},
	})
}
