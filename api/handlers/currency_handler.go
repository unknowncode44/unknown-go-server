package handlers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/unknowncode44/unknown-go-server/api/presenter"
	"github.com/unknowncode44/unknown-go-server/pkg/currency"
	"github.com/unknowncode44/unknown-go-server/pkg/entities"
)

// CurrencyHandler manages HTTP operations related to currencies.
type CurrencyHandler struct {
	service currency.Service
}

// NewCurrencyHandler returns a new CurrencyHandler using the given service.
func NewCurrencyHandler(s currency.Service) *CurrencyHandler {
	return &CurrencyHandler{service: s}
}

// Create handles POST /currencies
func (h *CurrencyHandler) Create(c *fiber.Ctx) error {
	var req presenter.CreateCurrencyRequest
	if err := c.BodyParser(&req); err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if strings.TrimSpace(req.Code) == "" || strings.TrimSpace(req.Name) == "" {
		return respondError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	cur := &entities.Currency{
		Code: strings.ToUpper(strings.TrimSpace(req.Code)),
		Name: strings.TrimSpace(req.Name),
	}

	created, err := h.service.Create(cur)
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(presenter.CurrencySuccessResponse{
		Success: true,
		Data:    presenter.ToCurrencyResponse(created),
	})
}

// GetAll handles GET /currencies
func (h *CurrencyHandler) GetAll(c *fiber.Ctx) error {
	items, err := h.service.FindAll()
	if err != nil {
		return respondError(c, fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(presenter.CurrencySuccessResponse{
		Success: true,
		Data:    presenter.ToCurrencyListResponse(items),
	})
}

// GetById handles GET /currencies/:id
func (h *CurrencyHandler) GetById(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid currency ID")
	}

	cur, err := h.service.FindByID(id)
	if err != nil {
		return respondError(c, fiber.StatusNotFound, err.Error())
	}

	return c.JSON(presenter.CurrencySuccessResponse{
		Success: true,
		Data:    presenter.ToCurrencyResponse(cur),
	})
}

// Update handles PATCH /currencies/:id
func (h *CurrencyHandler) Update(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid currency ID")
	}

	var req presenter.UpdateCurrencyRequest
	if err := c.BodyParser(&req); err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if strings.TrimSpace(req.Name) == "" {
		return respondError(c, fiber.StatusBadRequest, "No fields provided for update")
	}

	cur, err := h.service.FindByID(id)
	if err != nil {
		return respondError(c, fiber.StatusNotFound, "Currency with id provided was not found")
	}

	cur.Name = strings.TrimSpace(req.Name)

	updated, err := h.service.Update(cur)
	if err != nil {
		return respondError(c, fiber.StatusInternalServerError, "Currency was not updated due to an internal server error")
	}

	return c.JSON(presenter.CurrencySuccessResponse{
		Success: true,
		Data:    presenter.ToCurrencyResponse(updated),
	})
}

// Deactivate handles POST /currencies/:id/deactivate performing a logical delete.
func (h *CurrencyHandler) Deactivate(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid Currency ID")
	}

	if err := h.service.Deactivate(id); err != nil {
		return respondError(c, fiber.StatusNotFound, "Currency not found")
	}

	return c.SendStatus(fiber.StatusNoContent)
}
