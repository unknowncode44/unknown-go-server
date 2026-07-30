package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/unknowncode44/unknown-go-server/api/presenter"
	"github.com/unknowncode44/unknown-go-server/pkg/asset"
	"github.com/unknowncode44/unknown-go-server/pkg/entities"
	"github.com/unknowncode44/unknown-go-server/pkg/location"
)

type LocationHandler struct {
	service      location.Service
	assetService asset.Service
}

func NewLocationHandler(s location.Service, assetService asset.Service) *LocationHandler {
	return &LocationHandler{service: s, assetService: assetService}
}

func (h *LocationHandler) Create(c *fiber.Ctx) error {
	var req presenter.CreateLocationRequest
	if err := c.BodyParser(&req); err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	var parent *uuid.UUID
	if req.ParentLocationID != "" {
		pid, err := uuid.Parse(req.ParentLocationID)
		if err != nil {
			return respondError(c, fiber.StatusBadRequest, "invalid parent_location_id")
		}
		parent = &pid
	}

	loc := &entities.Location{
		Name:             req.Name,
		Type:             entities.LocationType(req.Type),
		ParentLocationID: parent,
	}

	created, err := h.service.Create(loc)
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, err.Error())
	}
	return c.Status(fiber.StatusCreated).JSON(presenter.LocationSuccessResponse{Success: true, Data: presenter.ToLocationResponse(created)})
}

func (h *LocationHandler) GetAll(c *fiber.Ctx) error {
	list, err := h.service.FindAll()
	if err != nil {
		return respondError(c, fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(presenter.LocationSuccessResponse{Success: true, Data: presenter.ToLocationListResponse(list)})
}

func (h *LocationHandler) GetByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid id")
	}
	l, err := h.service.FindByID(id)
	if err != nil {
		return respondError(c, fiber.StatusNotFound, err.Error())
	}
	return c.JSON(presenter.LocationSuccessResponse{Success: true, Data: presenter.ToLocationResponse(l)})
}

// PublicByID returns the location's name and the serialized assets
// currently located there — unauthenticated view for QR scans of
// shelves/estanterías. Bulk (BDC) materials aren't tracked per-location
// yet, so this only lists Assets.
func (h *LocationHandler) PublicByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid location ID")
	}

	loc, err := h.service.FindByID(id)
	if err != nil {
		return respondError(c, fiber.StatusNotFound, "Location not found")
	}

	assets, err := h.assetService.FindByLocationID(id)
	if err != nil {
		return respondError(c, fiber.StatusInternalServerError, err.Error())
	}

	resp := presenter.PublicLocationResponse{
		ID:   loc.ID.String(),
		Name: loc.Name,
		Type: string(loc.Type),
	}
	for _, a := range assets {
		var manufacturerSerial string
		if a.ManufacturerSerial != nil {
			manufacturerSerial = *a.ManufacturerSerial
		}
		resp.Assets = append(resp.Assets, presenter.PublicLocationAssetItem{
			SerialVisible:      a.SerialVisible,
			ManufacturerSerial: manufacturerSerial,
			MaterialName:       a.Material.Name,
			MaterialCode:       a.Material.Code,
			Status:             string(a.Status),
		})
	}

	return c.JSON(presenter.LocationSuccessResponse{Success: true, Data: resp})
}

func (h *LocationHandler) Update(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid id")
	}
	var req presenter.UpdateLocationRequest
	if err := c.BodyParser(&req); err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid request body")
	}
	loc, err := h.service.FindByID(id)
	if err != nil {
		return respondError(c, fiber.StatusNotFound, "Location not found")
	}
	if req.Name != "" {
		loc.Name = req.Name
	}
	if req.Type != "" {
		loc.Type = entities.LocationType(req.Type)
	}
	if req.ParentLocationID != "" {
		pid, err := uuid.Parse(req.ParentLocationID)
		if err != nil {
			return respondError(c, fiber.StatusBadRequest, "invalid parent_location_id")
		}
		loc.ParentLocationID = &pid
	}
	if req.IsActive != nil {
		loc.IsActive = *req.IsActive
	}
	updated, err := h.service.Update(loc)
	if err != nil {
		return respondError(c, fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(presenter.LocationSuccessResponse{Success: true, Data: presenter.ToLocationResponse(updated)})
}

func (h *LocationHandler) Deactivate(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid id")
	}
	if err := h.service.Deactivate(id); err != nil {
		return respondError(c, fiber.StatusNotFound, err.Error())
	}
	return c.SendStatus(fiber.StatusNoContent)
}
