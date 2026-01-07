// Package handlers contains HTTP handlers for API resources.
// Handlers validate requests, map input DTOs to domain entities,
// invoke the service layer, and format responses.
package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/unknowncode44/unknown-go-server/api/presenter"
	"github.com/unknowncode44/unknown-go-server/pkg/entities"
	"github.com/unknowncode44/unknown-go-server/pkg/vendor"
)

// VendorHandler handles HTTP requests for vendor resources.
// It depends on a vendor.Service to perform business logic and
// persistence operations.
type VendorHandler struct {
	service vendor.Service
}

// NewVendorHandler creates a new VendorHandler with the provided service.
func NewVendorHandler(service vendor.Service) *VendorHandler {
	return &VendorHandler{service: service}
}

// Create handles POST /vendors and creates a new vendor.
// It expects a JSON body matching presenter.CreateVendorRequest.
// On success it returns HTTP 201 with the created vendor.
func (h *VendorHandler) Create(c *fiber.Ctx) error {
	var req presenter.CreateVendorRequest

	if err := c.BodyParser(&req); err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	vendorEntity := &entities.Vendor{
		Name:  req.Name,
		Code:  req.Code,
		TaxID: req.TaxID,
	}

	created, err := h.service.Create(vendorEntity)
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, err.Error())
	}

	return c.Status(fiber.StatusCreated).
		JSON(presenter.VendorSuccessResponse{
			Success: true,
			Data:    presenter.ToVendorResponse(created),
		})
}

// FindAll handles GET /vendors and returns all vendors.
func (h *VendorHandler) FindAll(c *fiber.Ctx) error {
	vendors, err := h.service.FindAll()
	if err != nil {
		return respondError(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(presenter.VendorSuccessResponse{
		Success: true,
		Data:    presenter.ToVendorListResponse(vendors),
	})
}

// FindByID handles GET /vendors/:id and returns a vendor by UUID.
// Returns 400 for invalid UUIDs and 404 if the vendor is not found.
func (h *VendorHandler) FindByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid vendor ID")
	}

	v, err := h.service.FindByID(id)
	if err != nil {
		return respondError(c, fiber.StatusNotFound, err.Error())
	}

	return c.JSON(presenter.VendorSuccessResponse{
		Success: true,
		Data:    presenter.ToVendorResponse(v),
	})
}

// Update handles PUT /vendors/:id and applies updates to an existing vendor.
// It requires at least one updatable field in the request body.
func (h *VendorHandler) Update(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid vendor ID")
	}

	var req presenter.UpdateVendorRequest
	if err := c.BodyParser(&req); err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// require name on create/update as per higher-level requirement
	if req.Name == "" {
		return respondError(c, fiber.StatusBadRequest, "name is required")
	}

	// require at least one field to update
	if req.Name == "" && req.Code == "" && req.TaxID == "" {
		return respondError(c, fiber.StatusBadRequest, "No fields provided for update")
	}

	vend, err := h.service.FindByID(id)
	if err != nil {
		return respondError(c, fiber.StatusNotFound, "Vendor with ID provided was not found")
	}

	if req.Name != "" {
		vend.Name = req.Name
	}
	if req.Code != "" {
		vend.Code = req.Code
	}
	if req.TaxID != "" {
		vend.TaxID = req.TaxID
	}

	updated, err := h.service.Update(vend)
	if err != nil {
		return respondError(c, fiber.StatusInternalServerError, "Vendor was not updated due to an internal server error")
	}

	return c.JSON(presenter.VendorSuccessResponse{
		Success: true,
		Data:    presenter.ToVendorResponse(updated),
	})
}

// Deactivate handles DELETE /vendors/:id and performs a logical delete.
// Returns 204 on success.
func (h *VendorHandler) Deactivate(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid vendor ID")
	}

	if err := h.service.Deactivate(id); err != nil {
		return respondError(c, fiber.StatusNotFound, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}
