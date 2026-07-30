package presenter

import (
	"time"

	"github.com/unknowncode44/unknown-go-server/pkg/entities"
)

type CreateLocationRequest struct {
	Name             string `json:"name"`
	Type             string `json:"type"`
	ParentLocationID string `json:"parent_location_id,omitempty"`
}

type UpdateLocationRequest struct {
	Name             string `json:"name,omitempty"`
	Type             string `json:"type,omitempty"`
	ParentLocationID string `json:"parent_location_id,omitempty"`
	IsActive         *bool  `json:"is_active,omitempty"`
}

type LocationResponse struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	Type             string    `json:"type"`
	ParentLocationID string    `json:"parent_location_id,omitempty"`
	IsActive         bool      `json:"is_active"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// PublicLocationAssetItem es el detalle mínimo de un asset ubicado en una
// estantería, para la vista pública (QR).
type PublicLocationAssetItem struct {
	SerialVisible      string `json:"serial_visible"`
	ManufacturerSerial string `json:"manufacturer_serial,omitempty"`
	MaterialName       string `json:"material_name"`
	MaterialCode       string `json:"material_code"`
	Status             string `json:"status"`
}

// PublicLocationResponse es la vista pública (sin autenticación) de una
// ubicación y los assets serializados que contiene directamente.
type PublicLocationResponse struct {
	ID     string                    `json:"id"`
	Name   string                    `json:"name"`
	Type   string                    `json:"type"`
	Assets []PublicLocationAssetItem `json:"assets"`
}

type LocationSuccessResponse struct {
	Success bool        `json:"ok"`
	Data    interface{} `json:"data"`
}

type LocationErrorResponse struct {
	Success bool   `json:"ok"`
	Error   string `json:"error"`
}

func ToLocationResponse(l *entities.Location) LocationResponse {
	if l == nil {
		return LocationResponse{}
	}
	var parentID string
	if l.ParentLocationID != nil {
		parentID = l.ParentLocationID.String()
	}
	return LocationResponse{
		ID:               l.ID.String(),
		Name:             l.Name,
		Type:             string(l.Type),
		ParentLocationID: parentID,
		IsActive:         l.IsActive,
		CreatedAt:        l.CreatedAt,
		UpdatedAt:        l.UpdatedAt,
	}
}

func ToLocationListResponse(list []entities.Location) []LocationResponse {
	resp := make([]LocationResponse, 0, len(list))
	for _, l := range list {
		resp = append(resp, ToLocationResponse(&l))
	}
	return resp
}
