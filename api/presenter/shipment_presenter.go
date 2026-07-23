package presenter

import (
	"time"

	"github.com/unknowncode44/unknown-go-server/pkg/entities"
)

// CreateShipmentRequest is the body expected to create a shipment. Status,
// timestamps and the creator are managed server-side.
type CreateShipmentRequest struct {
	DestinationLocationID string `json:"destination_location_id"`
	RequesterName         string `json:"requester_name"`
	Notes                 string `json:"notes"`
}

// AddShipmentItemRequest is the body expected to add an item to a shipment.
// serial is required when type=ASSET; material_code and quantity are required
// when type=MATERIAL.
type AddShipmentItemRequest struct {
	Type         string  `json:"type"`
	Serial       string  `json:"serial"`
	MaterialCode string  `json:"material_code"`
	Quantity     float64 `json:"quantity"`
}

// LinkERPRequest is the body expected to reconcile a processed shipment against
// the external ERP movement.
type LinkERPRequest struct {
	ERPMovementNumber string `json:"erp_movement_number"`
}

// ShipmentItemResponse is the API representation of a shipment line item.
type ShipmentItemResponse struct {
	ID                        string   `json:"id"`
	Type                      string   `json:"type"`
	AssetID                   *string  `json:"asset_id,omitempty"`
	AssetSerial               *string  `json:"asset_serial,omitempty"`
	MaterialID                *string  `json:"material_id,omitempty"`
	MaterialCode              *string  `json:"material_code,omitempty"`
	MaterialName              *string  `json:"material_name,omitempty"`
	Quantity                  *float64 `json:"quantity,omitempty"`
	ResultingAssetMovementID  *string  `json:"resulting_asset_movement_id,omitempty"`
	ResultingDeliveryRecordID *string  `json:"resulting_delivery_record_id,omitempty"`
}

// ShipmentResponse is the API representation of a shipment.
type ShipmentResponse struct {
	ID                      string                 `json:"id"`
	DestinationLocationID   string                 `json:"destination_location_id"`
	DestinationLocationName string                 `json:"destination_location_name"`
	Status                  string                 `json:"status"`
	CreatedByUserID         string                 `json:"created_by_user_id"`
	RequesterName           string                 `json:"requester_name"`
	Notes                   string                 `json:"notes"`
	ErpMovementNumber       *string                `json:"erp_movement_number,omitempty"`
	ProcessedAt             *string                `json:"processed_at,omitempty"`
	ReconciledAt            *string                `json:"reconciled_at,omitempty"`
	Items                   []ShipmentItemResponse `json:"items"`
	CreatedAt               string                 `json:"created_at"`
}

// ShipmentSuccessResponse is the uniform success envelope for the domain.
type ShipmentSuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

func formatTimePtr(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := t.Format(time.RFC3339)
	return &s
}

// ToShipmentItemResponse maps a ShipmentItem entity to its API view.
func ToShipmentItemResponse(item *entities.ShipmentItem) ShipmentItemResponse {
	if item == nil {
		return ShipmentItemResponse{}
	}

	resp := ShipmentItemResponse{
		ID:   item.ID.String(),
		Type: string(item.Type),
	}

	if item.AssetID != nil {
		s := item.AssetID.String()
		resp.AssetID = &s
	}
	if item.Asset != nil && item.Asset.SerialVisible != "" {
		s := item.Asset.SerialVisible
		resp.AssetSerial = &s
	}
	if item.MaterialID != nil {
		s := item.MaterialID.String()
		resp.MaterialID = &s
	}
	if item.Material != nil {
		code := item.Material.Code
		name := item.Material.Name
		resp.MaterialCode = &code
		resp.MaterialName = &name
	}
	if item.Quantity != nil {
		q := *item.Quantity
		resp.Quantity = &q
	}
	if item.ResultingAssetMovementID != nil {
		s := item.ResultingAssetMovementID.String()
		resp.ResultingAssetMovementID = &s
	}
	if item.ResultingDeliveryRecordID != nil {
		s := item.ResultingDeliveryRecordID.String()
		resp.ResultingDeliveryRecordID = &s
	}

	return resp
}

// ToShipmentItemListResponse maps a slice of ShipmentItem entities.
func ToShipmentItemListResponse(list []entities.ShipmentItem) []ShipmentItemResponse {
	items := make([]ShipmentItemResponse, 0, len(list))
	for i := range list {
		items = append(items, ToShipmentItemResponse(&list[i]))
	}
	return items
}

// ToShipmentResponse maps a Shipment entity to its API view.
func ToShipmentResponse(s *entities.Shipment) ShipmentResponse {
	if s == nil {
		return ShipmentResponse{}
	}

	return ShipmentResponse{
		ID:                      s.ID.String(),
		DestinationLocationID:   s.DestinationLocationID.String(),
		DestinationLocationName: s.DestinationLocation.Name,
		Status:                  string(s.Status),
		CreatedByUserID:         s.CreatedByUserID.String(),
		RequesterName:           s.RequesterName,
		Notes:                   s.Notes,
		ErpMovementNumber:       s.ERPMovementNumber,
		ProcessedAt:             formatTimePtr(s.ProcessedAt),
		ReconciledAt:            formatTimePtr(s.ReconciledAt),
		Items:                   ToShipmentItemListResponse(s.Items),
		CreatedAt:               s.CreatedAt.Format(time.RFC3339),
	}
}

// ToShipmentListResponse maps a slice of Shipment entities.
func ToShipmentListResponse(list []entities.Shipment) []ShipmentResponse {
	items := make([]ShipmentResponse, 0, len(list))
	for i := range list {
		items = append(items, ToShipmentResponse(&list[i]))
	}
	return items
}
