package handlers

// Package handlers contains HTTP handlers for API resources.
// Handlers validate requests, map input DTOs to domain entities,
// invoke the service layer, and format responses.

import (
	"net/url"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/unknowncode44/unknown-go-server/api/presenter"
	"github.com/unknowncode44/unknown-go-server/pkg/delivery_record"
	"github.com/unknowncode44/unknown-go-server/pkg/entities"
	"github.com/unknowncode44/unknown-go-server/pkg/material"
	"github.com/unknowncode44/unknown-go-server/pkg/material_inventory"
)

// MaterialHandler manages HTTP operations related to materials.
// It delegates business logic to the provided material.Service.
type MaterialHandler struct {
	service   material.Service
	miService material_inventory.Service
	drService delivery_record.Service
}

// NewMaterialHandler returns a new MaterialHandler using the given services.
// miService/drService are used by the public endpoint to attach current stock.
func NewMaterialHandler(service material.Service, miService material_inventory.Service, drService delivery_record.Service) *MaterialHandler {
	return &MaterialHandler{service: service, miService: miService, drService: drService}
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

	// internal_code is optional and uniquely indexed: store empty as NULL so
	// multiple materials without an internal code don't collide on the index.
	var internalCode *string
	if req.InternalCode != "" {
		internalCode = &req.InternalCode
	}

	materialEntity := &entities.Material{
		Name:          req.Name,
		Sector:        req.Sector,
		Group:         entities.MaterialGroup(req.Group),
		UnitOfMeasure: req.UnitOfMeasure,
		InternalCode:  internalCode,
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

// GetMaterialsByERP handles GET /materials and returns the list of materials that matchs erpCode provided.
func (h *MaterialHandler) GetMaterialsByERP(c *fiber.Ctx) error {
	// Fiber's default config does not unescape path params, so decode it here:
	// real Bejerman "bag" ERP codes contain spaces (e.g. "0 MAT GOP21") that
	// the client must percent-encode (".../by-erp/0%20MAT%20GOP21").
	erpCode, err := url.PathUnescape(c.Params("erp_code"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid erp_code")
	}
	if erpCode == "" {
		return respondError(c, fiber.StatusBadRequest, "Invalid request body")
	}
	materials, err := h.service.FindByERPCode(erpCode)
	if err != nil {
		return respondError(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(presenter.MaterialSuccessResponse{
		Success: true,
		Data:    presenter.ToMaterialListResponse(materials),
	})
}

// PublicByCode returns a limited, unauthenticated view for QR scans of
// bulk (BDC) materials.
func (h *MaterialHandler) PublicByCode(c *fiber.Ctx) error {
	code := c.Params("code")
	mat, err := h.service.FindByCode(code)
	if err != nil {
		return respondError(c, fiber.StatusNotFound, "Material not found")
	}

	resp := presenter.PublicMaterialResponse{
		Code:          mat.Code,
		Name:          mat.Name,
		Sector:        mat.Sector,
		UnitOfMeasure: mat.UnitOfMeasure,
		Group:         string(mat.Group),
	}

	// Best-effort: si no hay inventario cargado para este material, la
	// respuesta igual es válida, solo sin current_stock.
	if inv, err := h.miService.FindByMaterialID(mat.ID); err == nil {
		if total, err := h.drService.GetTotalQuantity(inv.ID); err == nil {
			resp.CurrentStock = &total
		}
	}

	return c.JSON(presenter.MaterialSuccessResponse{Success: true, Data: resp})
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
	if req.Name == "" && req.Sector == "" && req.Group == "" && req.UnitOfMeasure == "" &&
		req.ERPCode == "" && req.InternalCode == "" && req.Code == "" {
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

	if req.Group != "" {
		mat.Group = entities.MaterialGroup(req.Group)
	}

	if req.UnitOfMeasure != "" {
		mat.UnitOfMeasure = req.UnitOfMeasure
	}

	if req.ERPCode != "" {
		mat.ERPCode = req.ERPCode
	}

	if req.InternalCode != "" {
		mat.InternalCode = &req.InternalCode
	}

	if req.Code != "" {
		mat.Code = req.Code
	}

	updated, err := h.service.Update(mat)
	if err != nil {
		// Devolvemos el error tal cual, igual que en Create: los errores de
		// validación del servicio (ej. group inválido) deben llegar como 400.
		return respondError(c, fiber.StatusBadRequest, err.Error())
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
