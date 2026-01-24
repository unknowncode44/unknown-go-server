package presenter

import (
	"time"

	"github.com/unknowncode44/unknown-go-server/pkg/entities"
)

type CreateAssetRequest struct {
	MaterialID         string `json:"material_id"`
	ManufacturerSerial string `json:"manufacturer_serial,omitempty"`
	ParentAssetID      string `json:"parent_asset_id,omitempty"`
	CurrentLocationID  string `json:"current_location_id,omitempty"`
}

type UpdateAssetRequest struct {
	ManufacturerSerial string `json:"manufacturer_serial,omitempty"`
	ParentAssetID      string `json:"parent_asset_id,omitempty"`
	CurrentLocationID  string `json:"current_location_id,omitempty"`
	Status             string `json:"status,omitempty"`
	IsActive           *bool  `json:"is_active,omitempty"`
}

type AssetResponse struct {
	ID                 string    `json:"id"`
	SerialVisible      string    `json:"serial_visible"`
	MaterialID         string    `json:"material_id"`
	ManufacturerSerial string    `json:"manufacturer_serial,omitempty"`
	Status             string    `json:"status"`
	ParentAssetID      string    `json:"parent_asset_id,omitempty"`
	CurrentLocationID  string    `json:"current_location_id,omitempty"`
	IsActive           bool      `json:"is_active"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type AssetSuccessResponse struct {
	Success bool        `json:"ok"`
	Data    interface{} `json:"data"`
}

type AssetErrorResponse struct {
	Success bool   `json:"ok"`
	Error   string `json:"error"`
}

// PublicAssetResponse limited view for public QR endpoint
type PublicAssetResponse struct {
	SerialVisible string `json:"serial_visible"`
	Material      struct {
		Name  string `json:"name"`
		Code  string `json:"code"`
		Model string `json:"model,omitempty"`
		Brand string `json:"brand,omitempty"`
	} `json:"material"`
	Status          string `json:"status"`
	CurrentLocation *struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"current_location,omitempty"`
	LastMovement *struct {
		Type string `json:"type"`
		Date string `json:"date"`
	} `json:"last_movement,omitempty"`
}

func ToAssetResponse(a *entities.Asset) AssetResponse {
	if a == nil {
		return AssetResponse{}
	}
	var parentID string
	if a.ParentAssetID != nil {
		parentID = a.ParentAssetID.String()
	}
	var locID string
	if a.CurrentLocationID != nil {
		locID = a.CurrentLocationID.String()
	}
	return AssetResponse{
		ID:                 a.ID.String(),
		SerialVisible:      a.SerialVisible,
		MaterialID:         a.MaterialID.String(),
		ManufacturerSerial: ptrToString(a.ManufacturerSerial),
		Status:             string(a.Status),
		ParentAssetID:      parentID,
		CurrentLocationID:  locID,
		IsActive:           a.IsActive,
		CreatedAt:          a.CreatedAt,
		UpdatedAt:          a.UpdatedAt,
	}
}

func ptrToString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
