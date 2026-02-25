package handlers

// Package handlers contains HTTP handlers for API resources.
// Handlers validate requests, map input DTOs to domain entities,
// invoke the service layer, and format responses.

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/unknowncode44/unknown-go-server/api/presenter"
	"github.com/unknowncode44/unknown-go-server/pkg/entities"
	"github.com/unknowncode44/unknown-go-server/pkg/material"
)

// MaterialHandler manages HTTP operations related to materials.
// It delegates business logic to the provided material.Service.
type MaterialHandler struct {
	service material.Service
}

// NewMaterialHandler returns a new MaterialHandler using the given service.
func NewMaterialHandler(service material.Service) *MaterialHandler {
	return &MaterialHandler{service: service}
}

// Create handles POST /materials. It validates the request body and
// required fields, converts the DTO to an entity, and delegates creation
// to the service. Returns 201 with the created material on success.
func (h *MaterialHandler) Create(c *fiber.Ctx) error {
	var req presenter.CreateMaterialRequest
	if err := c.BodyParser(&req); err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// Validación básica a nivel de handler (el servicio puede implementar validaciones adicionales)
	if req.Name == "" {
		return respondError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	materialEntity := &entities.Material{
		Name:          req.Name,
		Sector:        req.Sector,
		UnitOfMeasure: req.UnitOfMeasure,
		InternalCode:  req.InternalCode,
		ERPCode:       req.ERPCode,
		Code:          req.Code,
	}

	created, err := h.service.Create(materialEntity)
	if err != nil {
		// Devolvemos el error tal cual; la capa de servicio debe responsabilizarse de errores de negocio.
		return respondError(c, fiber.StatusBadRequest, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(presenter.MaterialSuccessResponse{
		Success: true,
		Data:    presenter.ToMaterialResponse(created),
	})
}

// GetAll handles GET /materials and returns the list of materials.
func (h *MaterialHandler) GetAll(c *fiber.Ctx) error {
	materials, err := h.service.FindAll()
	if err != nil {
		return respondError(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(presenter.MaterialSuccessResponse{
		Success: true,
		Data:    presenter.ToMaterialListResponse(materials),
	})
}

// GetById handles GET /materials/:id. It validates the UUID and
// returns the material or an appropriate error status.
func (h *MaterialHandler) GetById(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid material ID")
	}

	m, err := h.service.FindByID(id)
	if err != nil {
		return respondError(c, fiber.StatusNotFound, err.Error())
	}

	return c.JSON(presenter.MaterialSuccessResponse{
		Success: true,
		Data:    presenter.ToMaterialResponse(m),
	})
}

// Update handles PUT /materials/:id. It requires at least one field to
// update, applies changes to the entity and delegates persistence to the service.
func (h *MaterialHandler) Update(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid material ID")
	}

	var req presenter.UpdateMaterialRequest
	if err := c.BodyParser(&req); err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// require at least one field to update
	if req.Name == "" && req.Sector == "" && req.UnitOfMeasure == "" {
		return respondError(c, fiber.StatusBadRequest, "No fields provided for update")
	}

	mat, err := h.service.FindByID(id)
	if err != nil {
		return respondError(c, fiber.StatusNotFound, "Material with id provided was not found")
	}

	if req.Name != "" {
		mat.Name = req.Name
	}

	if req.Sector != "" {
		mat.Sector = req.Sector
	}

	if req.UnitOfMeasure != "" {
		mat.UnitOfMeasure = req.UnitOfMeasure
	}

	updated, err := h.service.Update(mat)
	if err != nil {
		return respondError(c, fiber.StatusInternalServerError, "Material was not updated due to an internal server error")
	}

	return c.JSON(presenter.MaterialSuccessResponse{
		Success: true,
		Data:    presenter.ToMaterialResponse(updated),
	})
}

// Deactivate handles DELETE /materials/:id performing a logical delete.
// Returns HTTP 204 on success.
func (h *MaterialHandler) Deactivate(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid Material ID")
	}

	if err := h.service.Deactivate(id); err != nil {
		return respondError(c, fiber.StatusNotFound, "Material not found")
	}

	return c.SendStatus(fiber.StatusNoContent)
}
