package handlers

// Package handlers contains HTTP handlers for API resources.
// The NeedHandler exposes the pre-funnel purchase "need" workflow.

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/unknowncode44/unknown-go-server/api/presenter"
	"github.com/unknowncode44/unknown-go-server/pkg/entities"
	"github.com/unknowncode44/unknown-go-server/pkg/need"
)

// NeedHandler manages HTTP operations related to needs and their items.
type NeedHandler struct {
	service need.Service
}

// NewNeedHandler returns a new NeedHandler using the given service.
func NewNeedHandler(service need.Service) *NeedHandler {
	return &NeedHandler{service: service}
}

// parseFlexibleDate accepts an empty string (nil), RFC3339 or YYYY-MM-DD.
func parseFlexibleDate(value string) (*time.Time, error) {
	if value == "" {
		return nil, nil
	}
	if t, err := time.Parse(time.RFC3339, value); err == nil {
		return &t, nil
	}
	t, err := time.Parse("2006-01-02", value)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// Create handles POST /needs.
func (h *NeedHandler) Create(c *fiber.Ctx) error {
	var req presenter.CreateNeedRequest
	if err := c.BodyParser(&req); err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	requiredDate, err := parseFlexibleDate(req.RequiredDate)
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid required_date format")
	}

	needEntity := &entities.Need{
		RequesterName: req.RequesterName,
		BuyerName:     req.BuyerName,
		CostCenter:    req.CostCenter,
		Justification: req.Justification,
		RequiredDate:  requiredDate,
	}

	created, err := h.service.Create(needEntity)
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(presenter.NeedSuccessResponse{
		Success: true,
		Data:    presenter.ToNeedResponse(created),
	})
}

// GetAll handles GET /needs.
func (h *NeedHandler) GetAll(c *fiber.Ctx) error {
	needs, err := h.service.FindAll()
	if err != nil {
		return respondError(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(presenter.NeedSuccessResponse{
		Success: true,
		Data:    presenter.ToNeedListResponse(needs),
	})
}

// GetByID handles GET /needs/:id.
func (h *NeedHandler) GetByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid need ID")
	}

	n, err := h.service.FindByID(id)
	if err != nil {
		return respondError(c, fiber.StatusNotFound, err.Error())
	}

	return c.JSON(presenter.NeedSuccessResponse{
		Success: true,
		Data:    presenter.ToNeedResponse(n),
	})
}

// Update handles PUT /needs/:id.
func (h *NeedHandler) Update(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid need ID")
	}

	var req presenter.UpdateNeedRequest
	if err := c.BodyParser(&req); err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	n, err := h.service.FindByID(id)
	if err != nil {
		return respondError(c, fiber.StatusNotFound, "Need with id provided was not found")
	}

	requiredDate, err := parseFlexibleDate(req.RequiredDate)
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid required_date format")
	}

	n.RequesterName = req.RequesterName
	n.BuyerName = req.BuyerName
	n.CostCenter = req.CostCenter
	n.Justification = req.Justification
	if requiredDate != nil {
		n.RequiredDate = requiredDate
	}

	updated, err := h.service.Update(n)
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(presenter.NeedSuccessResponse{
		Success: true,
		Data:    presenter.ToNeedResponse(updated),
	})
}

// Promote handles POST /needs/:id/promote.
func (h *NeedHandler) Promote(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid need ID")
	}

	n, err := h.service.Promote(id)
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(presenter.NeedSuccessResponse{
		Success: true,
		Data:    presenter.ToNeedResponse(n),
	})
}

// MarkConverted handles POST /needs/:id/mark-converted.
func (h *NeedHandler) MarkConverted(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid need ID")
	}

	n, err := h.service.MarkConverted(id)
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(presenter.NeedSuccessResponse{
		Success: true,
		Data:    presenter.ToNeedResponse(n),
	})
}

// Discard handles POST /needs/:id/discard.
func (h *NeedHandler) Discard(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid need ID")
	}

	n, err := h.service.Discard(id)
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(presenter.NeedSuccessResponse{
		Success: true,
		Data:    presenter.ToNeedResponse(n),
	})
}

// AddItem handles POST /needs/:id/items.
func (h *NeedHandler) AddItem(c *fiber.Ctx) error {
	needID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid need ID")
	}

	var req presenter.CreateNeedItemRequest
	if err := c.BodyParser(&req); err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	materialID, err := uuid.Parse(req.MaterialID)
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid material_id")
	}

	item := &entities.NeedItem{
		MaterialID: materialID,
		Quantity:   req.Quantity,
		Unit:       req.Unit,
		Notes:      req.Notes,
	}

	if req.SelectedCostID != "" {
		costID, err := uuid.Parse(req.SelectedCostID)
		if err != nil {
			return respondError(c, fiber.StatusBadRequest, "invalid selected_cost_id")
		}
		item.SelectedCostID = &costID
	}

	created, err := h.service.AddItem(needID, item)
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(presenter.NeedSuccessResponse{
		Success: true,
		Data:    presenter.ToNeedItemResponse(created),
	})
}

// ListItems handles GET /needs/:id/items.
func (h *NeedHandler) ListItems(c *fiber.Ctx) error {
	needID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid need ID")
	}

	items, err := h.service.ListItems(needID)
	if err != nil {
		return respondError(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(presenter.NeedSuccessResponse{
		Success: true,
		Data:    presenter.ToNeedItemListResponse(items),
	})
}

// UpdateItem handles PUT /needs/items/:itemId.
func (h *NeedHandler) UpdateItem(c *fiber.Ctx) error {
	itemID, err := uuid.Parse(c.Params("itemId"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid item ID")
	}

	var req presenter.UpdateNeedItemRequest
	if err := c.BodyParser(&req); err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	item := &entities.NeedItem{
		ID:       itemID,
		Quantity: req.Quantity,
		Unit:     req.Unit,
		Notes:    req.Notes,
	}

	if req.MaterialID != "" {
		materialID, err := uuid.Parse(req.MaterialID)
		if err != nil {
			return respondError(c, fiber.StatusBadRequest, "invalid material_id")
		}
		item.MaterialID = materialID
	}

	if req.SelectedCostID != "" {
		costID, err := uuid.Parse(req.SelectedCostID)
		if err != nil {
			return respondError(c, fiber.StatusBadRequest, "invalid selected_cost_id")
		}
		item.SelectedCostID = &costID
	}

	updated, err := h.service.UpdateItem(item)
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(presenter.NeedSuccessResponse{
		Success: true,
		Data:    presenter.ToNeedItemResponse(updated),
	})
}

// RemoveItem handles DELETE /needs/items/:itemId.
func (h *NeedHandler) RemoveItem(c *fiber.Ctx) error {
	itemID, err := uuid.Parse(c.Params("itemId"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid item ID")
	}

	if err := h.service.RemoveItem(itemID); err != nil {
		return respondError(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}
