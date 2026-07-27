package handlers

// Package handlers contains HTTP handlers for API resources.
// ShipmentHandler exposes the shipment (envíos) workflow: build a DRAFT
// shipment, process it to generate stock movements atomically, and reconcile
// it against the external ERP.

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/unknowncode44/unknown-go-server/api/presenter"
	"github.com/unknowncode44/unknown-go-server/pkg/entities"
	"github.com/unknowncode44/unknown-go-server/pkg/shipment"
)

// ShipmentHandler manages HTTP operations for shipments and their items.
type ShipmentHandler struct {
	service shipment.Service
}

// NewShipmentHandler returns a new ShipmentHandler using the given service.
func NewShipmentHandler(service shipment.Service) *ShipmentHandler {
	return &ShipmentHandler{service: service}
}

// Create handles POST /shipments.
func (h *ShipmentHandler) Create(c *fiber.Ctx) error {
	var req presenter.CreateShipmentRequest
	if err := c.BodyParser(&req); err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	destinationID, err := uuid.Parse(strings.TrimSpace(req.DestinationLocationID))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid destination_location_id")
	}

	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return respondError(c, fiber.StatusUnauthorized, "missing authenticated user")
	}

	created, err := h.service.Create(destinationID, userID, req.RequesterName, req.Notes)
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(presenter.ShipmentSuccessResponse{
		Success: true,
		Data:    presenter.ToShipmentResponse(created),
	})
}

// FindAll handles GET /shipments, with an optional ?status filter.
func (h *ShipmentHandler) FindAll(c *fiber.Ctx) error {
	var statusFilter *entities.ShipmentStatus
	if raw := strings.TrimSpace(c.Query("status")); raw != "" {
		status := entities.ShipmentStatus(strings.ToUpper(raw))
		switch status {
		case entities.ShipmentStatusDraft, entities.ShipmentStatusProcessed, entities.ShipmentStatusReconciled:
			statusFilter = &status
		default:
			return respondError(c, fiber.StatusBadRequest, "invalid status filter")
		}
	}

	shipments, err := h.service.FindAll(statusFilter)
	if err != nil {
		return respondError(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(presenter.ShipmentSuccessResponse{
		Success: true,
		Data:    presenter.ToShipmentListResponse(shipments),
	})
}

// FindByID handles GET /shipments/:id.
func (h *ShipmentHandler) FindByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid shipment ID")
	}

	s, err := h.service.FindByID(id)
	if err != nil {
		return respondError(c, fiber.StatusNotFound, "shipment not found")
	}

	return c.JSON(presenter.ShipmentSuccessResponse{
		Success: true,
		Data:    presenter.ToShipmentResponse(s),
	})
}

// AddItem handles POST /shipments/:id/items.
func (h *ShipmentHandler) AddItem(c *fiber.Ctx) error {
	shipmentID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid shipment ID")
	}

	var req presenter.AddShipmentItemRequest
	if err := c.BodyParser(&req); err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	itemType := entities.ShipmentItemType(strings.ToUpper(strings.TrimSpace(req.Type)))

	switch itemType {
	case entities.ShipmentItemTypeAsset:
		if strings.TrimSpace(req.Serial) == "" {
			return respondError(c, fiber.StatusBadRequest, "serial is required for type ASSET")
		}
		_, err = h.service.AddAssetItem(shipmentID, req.Serial)
	case entities.ShipmentItemTypeMaterial:
		if strings.TrimSpace(req.MaterialCode) == "" {
			return respondError(c, fiber.StatusBadRequest, "material_code is required for type MATERIAL")
		}
		if req.Quantity <= 0 {
			return respondError(c, fiber.StatusBadRequest, "quantity must be greater than zero for type MATERIAL")
		}
		_, err = h.service.AddMaterialItem(shipmentID, req.MaterialCode, req.Quantity)
	default:
		return respondError(c, fiber.StatusBadRequest, "type must be ASSET or MATERIAL")
	}

	if err != nil {
		return respondError(c, fiber.StatusBadRequest, err.Error())
	}

	// El frontend necesita el envío completo (con su array `items` actualizado)
	// para refrescar la lista, no el ítem suelto: se re-consulta el shipment,
	// que ya viene con Preload("Items").
	updated, err := h.service.FindByID(shipmentID)
	if err != nil {
		return respondError(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(presenter.ShipmentSuccessResponse{
		Success: true,
		Data:    presenter.ToShipmentResponse(updated),
	})
}

// RemoveItem handles DELETE /shipments/:id/items/:itemId.
func (h *ShipmentHandler) RemoveItem(c *fiber.Ctx) error {
	shipmentID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid shipment ID")
	}
	itemID, err := uuid.Parse(c.Params("itemId"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid item ID")
	}

	if err := h.service.RemoveItem(shipmentID, itemID); err != nil {
		return respondError(c, fiber.StatusBadRequest, err.Error())
	}

	// Igual que AddItem: se devuelve el envío completo ya sin el ítem, para que
	// el frontend actualice su lista con la misma fuente de verdad.
	updated, err := h.service.FindByID(shipmentID)
	if err != nil {
		return respondError(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(presenter.ShipmentSuccessResponse{
		Success: true,
		Data:    presenter.ToShipmentResponse(updated),
	})
}

// Process handles POST /shipments/:id/process.
func (h *ShipmentHandler) Process(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid shipment ID")
	}

	s, err := h.service.Process(id)
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(presenter.ShipmentSuccessResponse{
		Success: true,
		Data:    presenter.ToShipmentResponse(s),
	})
}

// LinkERP handles PUT /shipments/:id/erp-link.
func (h *ShipmentHandler) LinkERP(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid shipment ID")
	}

	var req presenter.LinkERPRequest
	if err := c.BodyParser(&req); err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	s, err := h.service.LinkERPMovement(id, req.ERPMovementNumber)
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(presenter.ShipmentSuccessResponse{
		Success: true,
		Data:    presenter.ToShipmentResponse(s),
	})
}

// Cancel handles DELETE /shipments/:id.
func (h *ShipmentHandler) Cancel(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid shipment ID")
	}

	if err := h.service.Cancel(id); err != nil {
		return respondError(c, fiber.StatusBadRequest, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}
