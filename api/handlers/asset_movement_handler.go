package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/unknowncode44/unknown-go-server/api/presenter"
	amsvc "github.com/unknowncode44/unknown-go-server/pkg/asset_movement"
	"github.com/unknowncode44/unknown-go-server/pkg/entities"
)

type AssetMovementHandler struct {
	service amsvc.Service
}

func NewAssetMovementHandler(s amsvc.Service) *AssetMovementHandler {
	return &AssetMovementHandler{service: s}
}

func (h *AssetMovementHandler) Create(c *fiber.Ctx) error {
	var req presenter.CreateAssetMovementRequest
	if err := c.BodyParser(&req); err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	aid, err := uuid.Parse(req.AssetID)
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid asset_id")
	}

	var from, to *uuid.UUID
	if req.FromLocationID != "" {
		fid, err := uuid.Parse(req.FromLocationID)
		if err != nil {
			return respondError(c, fiber.StatusBadRequest, "invalid from_location_id")
		}
		from = &fid
	}
	if req.ToLocationID != "" {
		tid, err := uuid.Parse(req.ToLocationID)
		if err != nil {
			return respondError(c, fiber.StatusBadRequest, "invalid to_location_id")
		}
		to = &tid
	}

	// parse movement date
	var mDate time.Time
	if req.MovementDate == "" {
		return respondError(c, fiber.StatusBadRequest, "movement_date is required")
	}
	mDate, err = time.Parse(time.RFC3339, req.MovementDate)
	if err != nil {
		mDate, err = time.Parse("2006-01-02", req.MovementDate)
		if err != nil {
			return respondError(c, fiber.StatusBadRequest, "invalid movement_date format")
		}
	}

	var notes *string
	if req.Notes != "" {
		notes = &req.Notes
	}

	am := &entities.AssetMovement{
		AssetID:        aid,
		Type:           entities.AssetMovementType(req.Type),
		FromLocationID: from,
		ToLocationID:   to,
		MovementDate:   mDate,
		Notes:          notes,
	}

	created, err := h.service.Create(am)
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(presenter.AssetMovementSuccessResponse{Success: true, Data: presenter.ToAssetMovementResponse(created)})
}

func (h *AssetMovementHandler) FindByAsset(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("asset_id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid asset id")
	}
	list, err := h.service.FindByAsset(id)
	if err != nil {
		return respondError(c, fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(presenter.AssetMovementSuccessResponse{Success: true, Data: presenter.ToAssetMovementListResponse(list)})
}

func (h *AssetMovementHandler) FindLastByAsset(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("asset_id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid asset id")
	}
	am, err := h.service.FindLastByAsset(id)
	if err != nil {
		return respondError(c, fiber.StatusNotFound, err.Error())
	}
	return c.JSON(presenter.AssetMovementSuccessResponse{Success: true, Data: presenter.ToAssetMovementResponse(am)})
}
