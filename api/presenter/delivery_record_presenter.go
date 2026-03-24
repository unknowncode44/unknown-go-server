package presenter

import (
	"time"

	"github.com/unknowncode44/unknown-go-server/pkg/entities"
)

type CreateDeliveryRecordRequest struct {
	MaterialInventoryID string  `json:"material_inventory_id"`
	Type                string  `json:"type"`
	Quantity            float64 `json:"quantity"`
	DeliveryDate        string  `json:"delivery_date"`
	Notes               string  `json:"notes,omitempty"`
}

type DeliveryRecordResponse struct {
	ID                  string    `json:"id"`
	MaterialInventoryID string    `json:"material_inventory_id"`
	Type                string    `json:"type"`
	Quantity            float64   `json:"quantity"`
	DeliveryDate        time.Time `json:"delivery_date"`
	Notes               *string   `json:"notes,omitempty"`
	CreatedAt           time.Time `json:"created_at"`
}

type DeliveryRecordSuccessResponse struct {
	Success bool        `json:"ok"`
	Data    interface{} `json:"data"`
}

type DeliveryRecordErrorResponse struct {
	Success bool   `json:"ok"`
	Error   string `json:"error"`
}

type InventoryTotalResponse struct {
	MaterialInventoryID string  `json:"material_inventory_id"`
	InitialQuantity     float64 `json:"initial_quantity"`
	TotalInbound        float64 `json:"total_inbound"`
	TotalOutbound       float64 `json:"total_outbound"`
	CurrentTotal        float64 `json:"current_total"`
}

func ToDeliveryRecordResponse(dr *entities.DeliveryRecord) DeliveryRecordResponse {
	if dr == nil {
		return DeliveryRecordResponse{}
	}
	return DeliveryRecordResponse{
		ID:                  dr.ID.String(),
		MaterialInventoryID: dr.MaterialInventoryID.String(),
		Type:                string(dr.Type),
		Quantity:            dr.Quantity,
		DeliveryDate:        dr.DeliveryDate,
		Notes:               dr.Notes,
		CreatedAt:           dr.CreatedAt,
	}
}

func ToDeliveryRecordListResponse(list []entities.DeliveryRecord) []DeliveryRecordResponse {
	resp := make([]DeliveryRecordResponse, 0, len(list))
	for _, dr := range list {
		resp = append(resp, ToDeliveryRecordResponse(&dr))
	}
	return resp
}
