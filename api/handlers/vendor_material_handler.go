package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/unknowncode44/unknown-go-server/api/presenter"
	"github.com/unknowncode44/unknown-go-server/pkg/entities"
	vmservice "github.com/unknowncode44/unknown-go-server/pkg/vendor_material"
)

type VendorMaterialHandler struct {
	service vmservice.Service
}

func NewVendorMaterialHandler(s vmservice.Service) *VendorMaterialHandler {
	return &VendorMaterialHandler{service: s}
}

func (h *VendorMaterialHandler) Create(c *fiber.Ctx) error {
	var req presenter.CreateVendorMaterialRequest
	if err := c.BodyParser(&req); err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// basic handler-level validation
	if req.VendorID == "" || req.MaterialID == "" || req.VendorCode == "" {
		return respondError(c, fiber.StatusBadRequest, "vendor_id, material_id and vendor_code are required")
	}

	vendorID, err := uuid.Parse(req.VendorID)
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid vendor_id")
	}
	materialID, err := uuid.Parse(req.MaterialID)
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid material_id")
	}

	vm := &entities.VendorMaterial{
		VendorID:   vendorID,
		MaterialID: materialID,
		VendorCode: req.VendorCode,
	}

	created, err := h.service.Create(vm)
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(presenter.VendorMaterialSuccessResponse{
		Success: true,
		Data:    presenter.ToVendorMaterialResponse(created),
	})
}

func (h *VendorMaterialHandler) GetAll(c *fiber.Ctx) error {
	vms, err := h.service.FindAll()
	if err != nil {
		return respondError(c, fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(presenter.VendorMaterialSuccessResponse{
		Success: true,
		Data:    presenter.ToVendorMaterialListResponse(vms),
	})
}

func (h *VendorMaterialHandler) GetByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid id")
	}

	vm, err := h.service.FindByID(id)
	if err != nil {
		return respondError(c, fiber.StatusNotFound, err.Error())
	}

	return c.JSON(presenter.VendorMaterialSuccessResponse{
		Success: true,
		Data:    presenter.ToVendorMaterialResponse(vm),
	})
}

func (h *VendorMaterialHandler) Update(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid id")
	}

	var req presenter.UpdateVendorMaterialRequest
	if err := c.BodyParser(&req); err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// require at least one field to update
	if req.VendorCode == "" {
		return respondError(c, fiber.StatusBadRequest, "No fields provided for update")
	}

	existing, err := h.service.FindByID(id)
	if err != nil {
		return respondError(c, fiber.StatusNotFound, "VendorMaterial with id provided was not found")
	}

	if req.VendorCode != "" {
		existing.VendorCode = req.VendorCode
	}

	updated, err := h.service.Update(existing)
	if err != nil {
		return respondError(c, fiber.StatusInternalServerError, "VendorMaterial was not updated due to an internal server error")
	}

	return c.JSON(presenter.VendorMaterialSuccessResponse{
		Success: true,
		Data:    presenter.ToVendorMaterialResponse(updated),
	})
}

func (h *VendorMaterialHandler) Deactivate(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid id")
	}

	if err := h.service.Deactivate(id); err != nil {
		return respondError(c, fiber.StatusNotFound, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}
