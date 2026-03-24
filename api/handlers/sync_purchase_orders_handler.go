package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/unknowncode44/unknown-go-server/api/presenter"
	"github.com/unknowncode44/unknown-go-server/pkg/sync_purchase_order"
)

type SyncPurchaseOrderHandler struct {
	service sync_purchase_order.Service
}

func NewSyncOrderHandler(service sync_purchase_order.Service) *SyncPurchaseOrderHandler {
	return &SyncPurchaseOrderHandler{service: service}
}

func (h *SyncPurchaseOrderHandler) SyncPurchaseOrders(c *fiber.Ctx) error {
	var request []presenter.SyncPurchaseOrderRequest

	// Parse JSON array from VBA
	if err := c.BodyParser(&request); err != nil {
		return c.Status(400).JSON(presenter.SyncResponse{Ok: false, Error: err.Error()})
	}

	count, err := h.service.SyncOrders(request)
	if err != nil {
		return c.Status(500).JSON(presenter.SyncResponse{Ok: false, Error: err.Error()})
	}

	return c.Status(201).JSON(presenter.SyncResponse{
		Ok:      true,
		Message: "Sync completed",
		Data:    fiber.Map{"records_synced": count},
	})
}

func (h *SyncPurchaseOrderHandler) GetPurchaseOrders(c *fiber.Ctx) error {
	orders, err := h.service.GetAllOrders()
	if err != nil {
		return c.Status(500).JSON(presenter.SyncResponse{
			Ok:    false,
			Error: "Could not retrieve orders: " + err.Error(),
		})
	}
	// Map entities to presenter if needed, or return directly
	return c.Status(200).JSON(presenter.SyncResponse{
		Ok:   true,
		Data: orders,
	})

}

// example of request http://localhost:3000/sync/purchase-orders/125/550e8400-e29b-41d4-a716-446655440000/TAX12345678
func (h *SyncPurchaseOrderHandler) UpdatePORecords(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 0)
	if err != nil {
		return c.Status(500).JSON(presenter.SyncResponse{
			Ok:    false,
			Error: "Parse error during parsing the record: " + err.Error(),
		})
	}
	mId, err := uuid.Parse(c.Params("material_id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid material ID")
	}
	vId := c.Params("vendor_tax_id")
	if vId == "" {
		return respondError(c, fiber.StatusBadRequest, "No vendor tax id in the params")
	}

	order, err := h.service.AssignLinks(uint(id), mId, vId)
	if err != nil {
		return respondError(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(200).JSON(presenter.SyncResponse{
		Ok:   true,
		Data: order,
	})

}
