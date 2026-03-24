package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/unknowncode44/unknown-go-server/api/handlers"
)

func PublicSyncPurchaseOrdersRoute(router fiber.Router, handler *handlers.SyncPurchaseOrderHandler) {
	router.Post("/sync/purchase-orders", handler.SyncPurchaseOrders)
	router.Get("/sync/purchase-orders", handler.GetPurchaseOrders)
	router.Patch("/sync/purchase-orders/:id/:material_id/:vendor_tax_id", handler.UpdatePORecords)
}
