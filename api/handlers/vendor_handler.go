package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/unknowncode44/unknown-go-server/api/presenter"
	"github.com/unknowncode44/unknown-go-server/pkg/entities"
	"github.com/unknowncode44/unknown-go-server/pkg/vendor"
)

type VendorHandler struct {
	service vendor.Service
}

func NewVendorHandler(service vendor.Service) *VendorHandler {
	return &VendorHandler{service: service}
}

// POST /vendors
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

// GET /vendors
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

// GET /vendors/:id
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

// PUT /vendors/:id
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

// DELETE /vendors/:id
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
