package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/unknowncode44/unknown-go-server/api/presenter"
	assetsvc "github.com/unknowncode44/unknown-go-server/pkg/asset"
	amsvc "github.com/unknowncode44/unknown-go-server/pkg/asset_movement"
	"github.com/unknowncode44/unknown-go-server/pkg/entities"
)

type AssetHandler struct {
	service    assetsvc.Service
	movService amsvc.Service
}

func NewAssetHandler(s assetsvc.Service, ms amsvc.Service) *AssetHandler {
	return &AssetHandler{service: s, movService: ms}
}

func (h *AssetHandler) Create(c *fiber.Ctx) error {
	var req presenter.CreateAssetRequest
	if err := c.BodyParser(&req); err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	mid, err := uuid.Parse(req.MaterialID)
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid material_id")
	}

	var parent *uuid.UUID
	if req.ParentAssetID != "" {
		pid, err := uuid.Parse(req.ParentAssetID)
		if err != nil {
			return respondError(c, fiber.StatusBadRequest, "invalid parent_asset_id")
		}
		parent = &pid
	}

	var loc *uuid.UUID
	if req.CurrentLocationID != "" {
		lid, err := uuid.Parse(req.CurrentLocationID)
		if err != nil {
			return respondError(c, fiber.StatusBadRequest, "invalid current_location_id")
		}
		loc = &lid
	}

	asset := &entities.Asset{
		MaterialID:         mid,
		ManufacturerSerial: &req.ManufacturerSerial,
		ParentAssetID:      parent,
		CurrentLocationID:  loc,
	}

	created, err := h.service.Create(asset)
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(presenter.AssetSuccessResponse{Success: true, Data: presenter.ToAssetResponse(created)})
}

func (h *AssetHandler) FindAll(c *fiber.Ctx) error {
	list, err := h.service.FindAll()
	if err != nil {
		return respondError(c, fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(presenter.AssetSuccessResponse{Success: true, Data: list})
}

func (h *AssetHandler) FindByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid id")
	}
	asset, err := h.service.FindByID(id)
	if err != nil {
		return respondError(c, fiber.StatusNotFound, err.Error())
	}
	return c.JSON(presenter.AssetSuccessResponse{Success: true, Data: presenter.ToAssetResponse(asset)})
}

func (h *AssetHandler) FindBySerial(c *fiber.Ctx) error {
	serial := c.Params("serial")
	asset, err := h.service.FindBySerial(serial)
	if err != nil {
		return respondError(c, fiber.StatusNotFound, err.Error())
	}
	return c.JSON(presenter.AssetSuccessResponse{Success: true, Data: presenter.ToAssetResponse(asset)})
}

// PublicBySerial returns a limited, unauthenticated view for QR scans.
func (h *AssetHandler) PublicBySerial(c *fiber.Ctx) error {
	serial := c.Params("serial")
	asset, err := h.service.FindBySerial(serial)
	if err != nil {
		return respondError(c, fiber.StatusNotFound, err.Error())
	}

	var resp presenter.PublicAssetResponse
	resp.SerialVisible = asset.SerialVisible
	resp.Material.Name = asset.Material.Name
	resp.Material.Code = asset.Material.Code
	resp.Status = string(asset.Status)
	if asset.CurrentLocationID != nil {
		resp.CurrentLocation = &struct {
			ID   string "json:\"id\""
			Name string "json:\"name\""
		}{
			ID:   asset.CurrentLocationID.String(),
			Name: asset.CurrentLocation.Name,
		}
	}

	// try to get last movement (best-effort)
	if h.movService != nil {
		last, err := h.movService.FindLastByAsset(asset.ID)
		if err == nil && last != nil {
			resp.LastMovement = &struct {
				Type string "json:\"type\""
				Date string "json:\"date\""
			}{
				Type: string(last.Type),
				Date: last.MovementDate.Format(time.RFC3339),
			}
		}
	}

	return c.JSON(resp)
}

func (h *AssetHandler) Update(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid id")
	}

	var req presenter.UpdateAssetRequest
	if err := c.BodyParser(&req); err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	asset, err := h.service.FindByID(id)
	if err != nil {
		return respondError(c, fiber.StatusNotFound, "Asset not found")
	}

	if req.ManufacturerSerial != "" {
		asset.ManufacturerSerial = &req.ManufacturerSerial
	}
	if req.ParentAssetID != "" {
		pid, err := uuid.Parse(req.ParentAssetID)
		if err != nil {
			return respondError(c, fiber.StatusBadRequest, "invalid parent_asset_id")
		}
		asset.ParentAssetID = &pid
	}
	if req.CurrentLocationID != "" {
		lid, err := uuid.Parse(req.CurrentLocationID)
		if err != nil {
			return respondError(c, fiber.StatusBadRequest, "invalid current_location_id")
		}
		asset.CurrentLocationID = &lid
	}
	if req.Status != "" {
		asset.Status = entities.AssetStatus(req.Status)
	}
	if req.IsActive != nil {
		asset.IsActive = *req.IsActive
	}

	updated, err := h.service.Update(asset)
	if err != nil {
		return respondError(c, fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(presenter.AssetSuccessResponse{Success: true, Data: presenter.ToAssetResponse(updated)})
}

func (h *AssetHandler) Deactivate(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid id")
	}
	if err := h.service.Deactivate(id); err != nil {
		return respondError(c, fiber.StatusNotFound, err.Error())
	}
	return c.SendStatus(fiber.StatusNoContent)
}
