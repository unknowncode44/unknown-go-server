package handlers

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/unknowncode44/unknown-go-server/api/presenter"
	"github.com/unknowncode44/unknown-go-server/pkg/entities"
	"github.com/unknowncode44/unknown-go-server/pkg/user"
)

// UserHandler handles HTTP requests for user resources.
// It depends on a user.Service to perform business logic and
// persistence operations.
type UserHandler struct {
	service user.Service
}

// NewUserHandler creates a new UserHandler with the provided service.
func NewUserHandler(service user.Service) *UserHandler {
	return &UserHandler{service: service}
}

// Create handles POST /users and creates a new user.
// It expects a JSON body matching presenter.CreateUserRequest.
// On success it returns HTTP 201 with the created user (sin el hash).
func (h *UserHandler) Create(c *fiber.Ctx) error {
	var req presenter.CreateUserRequest

	if err := c.BodyParser(&req); err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	userEntity := &entities.User{
		Name:  req.Name,
		Email: req.Email,
		Role:  entities.UserRole(req.Role),
	}

	created, err := h.service.Create(userEntity, req.Password)
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, err.Error())
	}

	return c.Status(fiber.StatusCreated).
		JSON(presenter.UserSuccessResponse{
			Success: true,
			Data:    presenter.ToUserResponse(created),
		})
}

// FindAll handles GET /users and returns all users.
func (h *UserHandler) FindAll(c *fiber.Ctx) error {
	users, err := h.service.FindAll()
	if err != nil {
		return respondError(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(presenter.UserSuccessResponse{
		Success: true,
		Data:    presenter.ToUserListResponse(users),
	})
}

// FindByID handles GET /users/:id and returns a user by UUID.
// Returns 400 for invalid UUIDs and 404 if the user is not found.
func (h *UserHandler) FindByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid user ID")
	}

	u, err := h.service.FindByID(id)
	if err != nil {
		return respondError(c, fiber.StatusNotFound, err.Error())
	}

	return c.JSON(presenter.UserSuccessResponse{
		Success: true,
		Data:    presenter.ToUserResponse(u),
	})
}

// Update handles PUT /users/:id and applies updates to an existing user.
// Solo admite nombre y rol; el password se cambia por /users/:id/password
// y el email no se edita en este alcance.
func (h *UserHandler) Update(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid user ID")
	}

	var req presenter.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// require at least one field to update
	if req.Name == "" && req.Role == "" {
		return respondError(c, fiber.StatusBadRequest, "No fields provided for update")
	}

	u, err := h.service.FindByID(id)
	if err != nil {
		return respondError(c, fiber.StatusNotFound, "User with ID provided was not found")
	}

	if req.Name != "" {
		u.Name = req.Name
	}
	if req.Role != "" {
		u.Role = entities.UserRole(req.Role)
	}

	updated, err := h.service.Update(u)
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(presenter.UserSuccessResponse{
		Success: true,
		Data:    presenter.ToUserResponse(updated),
	})
}

// SetPassword handles PUT /users/:id/password and resets a user's password
// (operación de Admin, sin flujo de email). Returns 204 on success.
func (h *UserHandler) SetPassword(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid user ID")
	}

	var req presenter.SetPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if err := h.service.SetPassword(id, req.Password); err != nil {
		return respondError(c, fiber.StatusBadRequest, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// Deactivate handles DELETE /users/:id and performs a logical delete.
// Returns 204 on success, 409 si se intenta desactivar al último admin activo.
func (h *UserHandler) Deactivate(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid user ID")
	}

	if err := h.service.Deactivate(id); err != nil {
		if errors.Is(err, user.ErrLastActiveAdmin) {
			return respondError(c, fiber.StatusConflict, err.Error())
		}
		return respondError(c, fiber.StatusNotFound, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}
