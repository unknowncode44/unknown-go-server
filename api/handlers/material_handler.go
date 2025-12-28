package handlers

// Package handlers contiene los controladores HTTP para los recursos de la API.
// Cada handler se encarga de validar la petición, convertir DTOs a entidades,
// invocar la capa de servicio y formatear la respuesta.

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/unknowncode44/unknown-go-server/api/presenter"
	"github.com/unknowncode44/unknown-go-server/pkg/entities"
	"github.com/unknowncode44/unknown-go-server/pkg/material"
)

// MaterialHandler gestiona las operaciones HTTP relacionadas con materiales.
// Contiene la dependencia a la capa de servicio para delegar la lógica de negocio.
type MaterialHandler struct {
	service material.Service
}

// NewMaterialHandler crea una nueva instancia de MaterialHandler con la dependencia inyectada.
func NewMaterialHandler(service material.Service) *MaterialHandler {
	return &MaterialHandler{service: service}
}

// Create procesa la solicitud POST para crear un nuevo material.
// - Valida el body de la petición y los campos obligatorios.
// - Convierte el DTO a entidad y delega la creación al servicio.
// - Devuelve 201 con el material creado o 400/500 según corresponda.
func (h *MaterialHandler) Create(c *fiber.Ctx) error {
	var req presenter.CreateMaterialRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(presenter.MaterialErrorResponse{Success: false, Error: "Invalid request body"})
	}

	// Validación básica a nivel de handler (el servicio puede implementar validaciones adicionales)
	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).
			JSON(presenter.MaterialErrorResponse{Success: false, Error: "Invalid request body"})
	}

	materialEntity := &entities.Material{
		Name:          req.Name,
		Sector:        req.Sector,
		UnitOfMeasure: req.UnitOfMeasure,
	}

	created, err := h.service.Create(materialEntity)
	if err != nil {
		// Devolvemos el error tal cual; la capa de servicio debe responsabilizarse de errores de negocio.
		return c.Status(fiber.StatusBadRequest).
			JSON(presenter.MaterialErrorResponse{Success: false, Error: err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(presenter.MaterialSuccessResponse{
		Success: true,
		Data:    presenter.ToMaterialResponse(created),
	})
}

// GetAll devuelve la lista completa de materiales (GET).
// Prealoca el slice de respuesta para mejorar rendimiento en colecciones grandes.
func (h *MaterialHandler) GetAll(c *fiber.Ctx) error {
	materials, err := h.service.FindAll()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(presenter.MaterialErrorResponse{Success: false, Error: err.Error()})
	}

	return c.JSON(presenter.MaterialSuccessResponse{
		Success: true,
		Data:    presenter.ToMaterialListResponse(materials),
	})
}

// GetById devuelve un material por su ID (GET /:id).
// Valida que el ID tenga formato UUID y delega la búsqueda al servicio.
func (h *MaterialHandler) GetById(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(presenter.MaterialErrorResponse{Success: false, Error: "Invalid material ID"})
	}

	materialEntity, err := h.service.FindByID(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(presenter.MaterialErrorResponse{Success: false, Error: err.Error()})
	}

	return c.JSON(presenter.MaterialSuccessResponse{
		Success: true,
		Data:    presenter.ToMaterialResponse(materialEntity),
	})
}

// Update aplica cambios sobre un material existente (PUT /:id).
// - Valida ID y body, busca la entidad y delega la actualización al servicio.
func (h *MaterialHandler) Update(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(presenter.MaterialErrorResponse{Success: false, Error: "Invalid material ID"})
	}

	var req presenter.UpdateMaterialRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(presenter.MaterialErrorResponse{Success: false, Error: "Invalid request body"})
	}

	materialEntity, err := h.service.FindByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).
			JSON(presenter.MaterialErrorResponse{Success: false, Error: "Material No Encontrado"})
	}

	// Aplicar cambios permitidos desde el DTO
	materialEntity.Name = req.Name
	materialEntity.Sector = req.Sector
	materialEntity.UnitOfMeasure = req.UnitOfMeasure
	materialEntity.IsActive = req.IsActive

	updated, err := h.service.Update(materialEntity)
	if err != nil {
		return c.Status(fiber.StatusNotFound).
			JSON(presenter.MaterialErrorResponse{Success: false, Error: err.Error()})
	}

	return c.JSON(presenter.MaterialSuccessResponse{
		Success: true,
		Data:    presenter.ToMaterialResponse(updated),
	})
}

// Deactivate realiza el borrado lógico de un material (DELETE /:id).
// Devuelve 204 cuando la operación es exitosa.
func (h *MaterialHandler) Deactivate(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusNotFound).
			JSON(presenter.MaterialErrorResponse{Success: false, Error: "ID de Material invalido"})
	}

	if err := h.service.Deactivate(id); err != nil {
		return c.Status(fiber.StatusNotFound).
			JSON(presenter.MaterialErrorResponse{Success: false, Error: "Material No Encontrado"})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
