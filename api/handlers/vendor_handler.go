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
		return c.Status(fiber.StatusBadRequest).
			JSON(presenter.VendorErrorResponse{Error: "Invalid request body"})
	}

	vendorEntity := &entities.Vendor{
		Name:  req.Name,
		Code:  req.Code,
		TaxID: req.TaxID,
	}

	created, err := h.service.Create(vendorEntity)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(presenter.VendorErrorResponse{Error: err.Error()})
	}

	return c.Status(fiber.StatusCreated).
		JSON(presenter.VendorSuccessResponse{
			Data: presenter.ToVendorResponse(created),
		})
}

// GET /vendors
func (h *VendorHandler) FindAll(c *fiber.Ctx) error {
	vendors, err := h.service.FindAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(presenter.VendorErrorResponse{Error: err.Error()})
	}

	return c.JSON(presenter.VendorSuccessResponse{
		Data: presenter.ToVendorListResponse(vendors),
	})
}

// GET /vendors/:id
func (h *VendorHandler) FindByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(presenter.VendorErrorResponse{Error: "Invalid vendor ID"})
	}

	v, err := h.service.FindByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).
			JSON(presenter.VendorErrorResponse{Error: err.Error()})
	}

	return c.JSON(presenter.VendorSuccessResponse{
		Data: presenter.ToVendorResponse(v),
	})
}

// PUT /vendors/:id
func (h *VendorHandler) Update(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(presenter.VendorErrorResponse{Error: "Invalid vendor ID"})
	}

	var v entities.Vendor
	if err := c.BodyParser(&v); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(presenter.VendorErrorResponse{Error: "Invalid request body"})
	}

	v.ID = id

	updated, err := h.service.Update(&v)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(presenter.VendorErrorResponse{Error: err.Error()})
	}

	return c.JSON(presenter.VendorSuccessResponse{
		Data: presenter.ToVendorResponse(updated),
	})
}

// DELETE /vendors/:id
func (h *VendorHandler) Deactivate(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(presenter.VendorErrorResponse{Error: "Invalid vendor ID"})
	}

	if err := h.service.Deactivate(id); err != nil {
		return c.Status(fiber.StatusNotFound).
			JSON(presenter.VendorErrorResponse{Error: err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
