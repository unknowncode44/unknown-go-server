package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/unknowncode44/unknown-go-server/api/presenter"
	"github.com/unknowncode44/unknown-go-server/pkg/delivery_record"
	"github.com/unknowncode44/unknown-go-server/pkg/entities"
)

type DeliveryRecordHandler struct {
	service delivery_record.Service
}

func NewDeliveryRecordHandler(s delivery_record.Service) *DeliveryRecordHandler {
	return &DeliveryRecordHandler{service: s}
}

func (h *DeliveryRecordHandler) Create(c *fiber.Ctx) error {
	var req presenter.CreateDeliveryRecordRequest
	if err := c.BodyParser(&req); err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	invID, err := uuid.Parse(req.MaterialInventoryID)
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid material_inventory_id")
	}

	// validate type
	drType := entities.DeliveryType(req.Type)
	if drType != entities.DeliveryInbound && drType != entities.DeliveryOutbound {
		return respondError(c, fiber.StatusBadRequest, "type must be INBOUND or OUTBOUND")
	}

	// parse delivery_date
	var deliveryDate time.Time
	if req.DeliveryDate == "" {
		return respondError(c, fiber.StatusBadRequest, "delivery_date is required")
	}
	deliveryDate, err = time.Parse(time.RFC3339, req.DeliveryDate)
	if err != nil {
		deliveryDate, err = time.Parse("2006-01-02", req.DeliveryDate)
		if err != nil {
			return respondError(c, fiber.StatusBadRequest, "invalid delivery_date format")
		}
	}

	var notes *string
	if req.Notes != "" {
		notes = &req.Notes
	}

	dr := &entities.DeliveryRecord{
		MaterialInventoryID: invID,
		Type:                drType,
		Quantity:            req.Quantity,
		DeliveryDate:        deliveryDate,
		Notes:               notes,
	}

	created, err := h.service.Create(dr)
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(presenter.DeliveryRecordSuccessResponse{
		Success: true,
		Data:    presenter.ToDeliveryRecordResponse(created),
	})
}

func (h *DeliveryRecordHandler) GetAll(c *fiber.Ctx) error {
	list, err := h.service.FindAll()
	if err != nil {
		return respondError(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(presenter.DeliveryRecordSuccessResponse{
		Success: true,
		Data:    presenter.ToDeliveryRecordListResponse(list),
	})
}

func (h *DeliveryRecordHandler) FindByInventoryID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid inventory ID")
	}

	list, err := h.service.FindByInventoryID(id)
	if err != nil {
		return respondError(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(presenter.DeliveryRecordSuccessResponse{
		Success: true,
		Data:    presenter.ToDeliveryRecordListResponse(list),
	})
}
