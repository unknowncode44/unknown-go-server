package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/unknowncode44/unknown-go-server/api/presenter"
	"github.com/unknowncode44/unknown-go-server/pkg/entities"
	"github.com/unknowncode44/unknown-go-server/pkg/material_cost"
)

type MaterialCostHandler struct {
	service material_cost.Service
}

func NewMaterialCostHandler(s material_cost.Service) *MaterialCostHandler {
	return &MaterialCostHandler{service: s}
}

func (h *MaterialCostHandler) Create(c *fiber.Ctx) error {
	var req presenter.CreateMaterialCostRequest
	if err := c.BodyParser(&req); err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	mid, err := uuid.Parse(req.MaterialID)
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid material_id")
	}
	vid, err := uuid.Parse(req.VendorID)
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid vendor_id")
	}
	cid, err := uuid.Parse(req.CurrencyID)
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid currency_id")
	}

	// Parse cost_date: accept RFC3339 or YYYY-MM-DD
	var costDate time.Time
	if req.CostDate == "" {
		return respondError(c, fiber.StatusBadRequest, "cost_date is required")
	}
	costDate, err = time.Parse(time.RFC3339, req.CostDate)
	if err != nil {
		costDate, err = time.Parse("2006-01-02", req.CostDate)
		if err != nil {
			return respondError(c, fiber.StatusBadRequest, "invalid cost_date format")
		}
	}

	mc := &entities.MaterialCost{
		MaterialID: mid,
		VendorID:   vid,
		CurrencyID: cid,
		Cost:       req.Cost,
		CostDate:   costDate,
	}

	created, err := h.service.Create(mc)
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(presenter.MaterialCostSuccessResponse{
		Success: true,
		Data:    presenter.ToMaterialCostResponse(created),
	})
}

func (h *MaterialCostHandler) FindAll(c *fiber.Ctx) error {
	list, err := h.service.FindAll()
	if err != nil {
		return respondError(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(presenter.MaterialCostSuccessResponse{
		Success: true,
		Data:    presenter.ToMaterialCostListResponse(list),
	})
}

func (h *MaterialCostHandler) FindByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid id")
	}

	mc, err := h.service.FindByID(id)
	if err != nil {
		return respondError(c, fiber.StatusNotFound, err.Error())
	}

	return c.JSON(presenter.MaterialCostSuccessResponse{
		Success: true,
		Data:    presenter.ToMaterialCostResponse(mc),
	})
}

func (h *MaterialCostHandler) FindByMaterial(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid material id")
	}

	list, err := h.service.FindByMaterial(id)
	if err != nil {
		return respondError(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(presenter.MaterialCostSuccessResponse{
		Success: true,
		Data:    presenter.ToMaterialCostListResponse(list),
	})
}
