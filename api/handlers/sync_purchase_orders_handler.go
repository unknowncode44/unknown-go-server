package handlers

import (
	"github.com/gofiber/fiber/v2"
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
