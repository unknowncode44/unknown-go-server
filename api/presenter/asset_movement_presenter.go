package presenter

import (
	"time"

	"github.com/unknowncode44/unknown-go-server/pkg/entities"
)

type CreateAssetMovementRequest struct {
	AssetID        string `json:"asset_id"`
	Type           string `json:"type"`
	FromLocationID string `json:"from_location_id,omitempty"`
	ToLocationID   string `json:"to_location_id,omitempty"`
	MovementDate   string `json:"movement_date"` // RFC3339 or YYYY-MM-DD
	Notes          string `json:"notes,omitempty"`
}

type AssetMovementResponse struct {
	ID             string    `json:"id"`
	AssetID        string    `json:"asset_id"`
	Type           string    `json:"type"`
	FromLocationID string    `json:"from_location_id,omitempty"`
	ToLocationID   string    `json:"to_location_id,omitempty"`
	MovementDate   time.Time `json:"movement_date"`
	Notes          *string   `json:"notes,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}

type AssetMovementSuccessResponse struct {
	Success bool        `json:"ok"`
	Data    interface{} `json:"data"`
}

type AssetMovementErrorResponse struct {
	Success bool   `json:"ok"`
	Error   string `json:"error"`
}

func ToAssetMovementResponse(am *entities.AssetMovement) AssetMovementResponse {
	if am == nil {
		return AssetMovementResponse{}
	}
	var fromID, toID string
	if am.FromLocationID != nil {
		fromID = am.FromLocationID.String()
	}
	if am.ToLocationID != nil {
		toID = am.ToLocationID.String()
	}
	return AssetMovementResponse{
		ID:             am.ID.String(),
		AssetID:        am.AssetID.String(),
		Type:           string(am.Type),
		FromLocationID: fromID,
		ToLocationID:   toID,
		MovementDate:   am.MovementDate,
		Notes:          am.Notes,
		CreatedAt:      am.CreatedAt,
	}
}

func ToAssetMovementListResponse(list []entities.AssetMovement) []AssetMovementResponse {
	resp := make([]AssetMovementResponse, 0, len(list))
	for _, m := range list {
		resp = append(resp, ToAssetMovementResponse(&m))
	}
	return resp
}
