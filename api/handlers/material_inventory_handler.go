package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/unknowncode44/unknown-go-server/api/presenter"
	"github.com/unknowncode44/unknown-go-server/pkg/delivery_record"
	"github.com/unknowncode44/unknown-go-server/pkg/entities"
	"github.com/unknowncode44/unknown-go-server/pkg/material_inventory"
)

type MaterialInventoryHandler struct {
	service   material_inventory.Service
	drService delivery_record.Service
}

func NewMaterialInventoryHandler(s material_inventory.Service, drService delivery_record.Service) *MaterialInventoryHandler {
	return &MaterialInventoryHandler{service: s, drService: drService}
}

func (h *MaterialInventoryHandler) Create(c *fiber.Ctx) error {
	var req presenter.CreateMaterialInventoryRequest
	if err := c.BodyParser(&req); err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	mid, err := uuid.Parse(req.MaterialID)
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid material_id")
	}

	// parse start_date
	var startDate time.Time
	if req.StartDate == "" {
		return respondError(c, fiber.StatusBadRequest, "start_date is required")
	}
	startDate, err = time.Parse(time.RFC3339, req.StartDate)
	if err != nil {
		startDate, err = time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			return respondError(c, fiber.StatusBadRequest, "invalid start_date format")
		}
	}

	inv := &entities.MaterialInventory{
		MaterialID:      mid,
		InitialQuantity: req.InitialQuantity,
		StartDate:       startDate,
	}

	created, err := h.service.Create(inv)
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(presenter.MaterialInventorySuccessResponse{
		Success: true,
		Data:    presenter.ToMaterialInventoryResponse(created),
	})
}

func (h *MaterialInventoryHandler) GetAll(c *fiber.Ctx) error {
	list, err := h.service.FindAll()
	if err != nil {
		return respondError(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(presenter.MaterialInventorySuccessResponse{
		Success: true,
		Data:    presenter.ToMaterialInventoryListResponse(list),
	})
}

func (h *MaterialInventoryHandler) GetByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid inventory ID")
	}

	inv, err := h.service.FindByID(id)
	if err != nil {
		return respondError(c, fiber.StatusNotFound, err.Error())
	}

	return c.JSON(presenter.MaterialInventorySuccessResponse{
		Success: true,
		Data:    presenter.ToMaterialInventoryResponse(inv),
	})
}

func (h *MaterialInventoryHandler) GetByMaterialID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid material ID")
	}

	inv, err := h.service.FindByMaterialID(id)
	if err != nil {
		return respondError(c, fiber.StatusNotFound, err.Error())
	}

	return c.JSON(presenter.MaterialInventorySuccessResponse{
		Success: true,
		Data:    presenter.ToMaterialInventoryResponse(inv),
	})
}

func (h *MaterialInventoryHandler) Update(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid inventory ID")
	}

	var req presenter.UpdateMaterialInventoryRequest
	if err := c.BodyParser(&req); err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	inv, err := h.service.FindByID(id)
	if err != nil {
		return respondError(c, fiber.StatusNotFound, "Material inventory not found")
	}

	if req.InitialQuantity != nil {
		inv.InitialQuantity = *req.InitialQuantity
	}

	if req.StartDate != "" {
		startDate, err := time.Parse(time.RFC3339, req.StartDate)
		if err != nil {
			startDate, err = time.Parse("2006-01-02", req.StartDate)
			if err != nil {
				return respondError(c, fiber.StatusBadRequest, "invalid start_date format")
			}
		}
		inv.StartDate = startDate
	}

	updated, err := h.service.Update(inv)
	if err != nil {
		return respondError(c, fiber.StatusInternalServerError, "Material inventory was not updated due to an internal server error")
	}

	return c.JSON(presenter.MaterialInventorySuccessResponse{
		Success: true,
		Data:    presenter.ToMaterialInventoryResponse(updated),
	})
}

func (h *MaterialInventoryHandler) Deactivate(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid inventory ID")
	}

	if err := h.service.Deactivate(id); err != nil {
		return respondError(c, fiber.StatusNotFound, "Material inventory not found")
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *MaterialInventoryHandler) GetTotal(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid inventory ID")
	}

	inv, err := h.service.FindByID(id)
	if err != nil {
		return respondError(c, fiber.StatusNotFound, "Material inventory not found")
	}

	records, err := h.drService.FindByInventoryID(id)
	if err != nil {
		return respondError(c, fiber.StatusInternalServerError, err.Error())
	}

	var inbound, outbound float64
	for _, r := range records {
		switch r.Type {
		case entities.DeliveryInbound:
			inbound += r.Quantity
		case entities.DeliveryOutbound:
			outbound += r.Quantity
		}
	}

	resp := presenter.InventoryTotalResponse{
		MaterialInventoryID: id.String(),
		InitialQuantity:     inv.InitialQuantity,
		TotalInbound:        inbound,
		TotalOutbound:       outbound,
		CurrentTotal:        inv.InitialQuantity + inbound - outbound,
	}

	return c.JSON(presenter.DeliveryRecordSuccessResponse{
		Success: true,
		Data:    resp,
	})
}

func (h *MaterialInventoryHandler) FindDeliveriesByInventoryID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid inventory ID")
	}

	list, err := h.drService.FindByInventoryID(id)
	if err != nil {
		return respondError(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(presenter.DeliveryRecordSuccessResponse{
		Success: true,
		Data:    presenter.ToDeliveryRecordListResponse(list),
	})
}
