package presenter

import (
	"time"

	"github.com/unknowncode44/unknown-go-server/pkg/entities"
)

// CreateNeedRequest is the body expected to create a need. Number and Status
// are managed server-side and therefore not accepted here.
type CreateNeedRequest struct {
	RequesterName string `json:"requester_name"`
	BuyerName     string `json:"buyer_name"`
	CostCenter    string `json:"cost_center"`
	Justification string `json:"justification"`
	RequiredDate  string `json:"required_date"` // RFC3339 or YYYY-MM-DD, optional
}

// UpdateNeedRequest holds the editable fields of a need (everything except
// Number and Status, which are controlled by the server / transitions).
type UpdateNeedRequest struct {
	RequesterName string `json:"requester_name"`
	BuyerName     string `json:"buyer_name"`
	CostCenter    string `json:"cost_center"`
	Justification string `json:"justification"`
	RequiredDate  string `json:"required_date"` // RFC3339 or YYYY-MM-DD, optional
}

// NeedResponse is the API representation of a need.
type NeedResponse struct {
	ID            string             `json:"id"`
	Number        string             `json:"number"`
	RequesterName string             `json:"requester_name"`
	BuyerName     string             `json:"buyer_name"`
	CostCenter    string             `json:"cost_center"`
	Justification string             `json:"justification"`
	RequiredDate  *time.Time         `json:"required_date,omitempty"`
	Status        string             `json:"status"`
	PromotedAt    *time.Time         `json:"promoted_at,omitempty"`
	ConvertedAt   *time.Time         `json:"converted_at,omitempty"`
	IsActive      bool               `json:"is_active"`
	Items         []NeedItemResponse `json:"items,omitempty"`
	CreatedAt     time.Time          `json:"created_at"`
	UpdatedAt     time.Time          `json:"updated_at"`
}

// CreateNeedItemRequest is the body expected to add an item to a need.
type CreateNeedItemRequest struct {
	MaterialID     string  `json:"material_id"`
	Quantity       float64 `json:"quantity"`
	Unit           string  `json:"unit"`
	SelectedCostID string  `json:"selected_cost_id"` // optional MaterialCost UUID
	Notes          string  `json:"notes"`
}

// UpdateNeedItemRequest holds the editable fields of an item line.
type UpdateNeedItemRequest struct {
	MaterialID     string  `json:"material_id"`
	Quantity       float64 `json:"quantity"`
	Unit           string  `json:"unit"`
	SelectedCostID string  `json:"selected_cost_id"`
	Notes          string  `json:"notes"`
}

// NeedItemResponse is the API representation of a need line item. It embeds the
// material name and ERP code so the client can render without a second round-trip.
type NeedItemResponse struct {
	ID              string    `json:"id"`
	NeedID          string    `json:"need_id"`
	MaterialID      string    `json:"material_id"`
	MaterialName    string    `json:"material_name"`
	MaterialERPCode string    `json:"material_erp_code"`
	Quantity        float64   `json:"quantity"`
	Unit            string    `json:"unit"`
	SelectedCostID  *string   `json:"selected_cost_id,omitempty"`
	Notes           string    `json:"notes"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// NeedSuccessResponse represents a successful response with one or more needs.
type NeedSuccessResponse struct {
	Success bool        `json:"ok"`
	Data    interface{} `json:"data"`
}

// NeedErrorResponse represents a standard error response for the need domain.
type NeedErrorResponse struct {
	Success bool   `json:"ok"`
	Error   string `json:"error"`
}

// ToNeedResponse transforms a Need entity into a NeedResponse, embedding its
// items when they have been loaded.
func ToNeedResponse(n *entities.Need) NeedResponse {
	if n == nil {
		return NeedResponse{}
	}

	return NeedResponse{
		ID:            n.ID.String(),
		Number:        n.Number,
		RequesterName: n.RequesterName,
		BuyerName:     n.BuyerName,
		CostCenter:    n.CostCenter,
		Justification: n.Justification,
		RequiredDate:  n.RequiredDate,
		Status:        string(n.Status),
		PromotedAt:    n.PromotedAt,
		ConvertedAt:   n.ConvertedAt,
		IsActive:      n.IsActive,
		Items:         ToNeedItemListResponse(n.Items),
		CreatedAt:     n.CreatedAt,
		UpdatedAt:     n.UpdatedAt,
	}
}

// ToNeedListResponse transforms a list of Need entities.
func ToNeedListResponse(list []entities.Need) []NeedResponse {
	response := make([]NeedResponse, 0, len(list))
	for _, n := range list {
		item := n
		response = append(response, ToNeedResponse(&item))
	}
	return response
}

// ToNeedItemResponse transforms a NeedItem entity into a NeedItemResponse.
func ToNeedItemResponse(item *entities.NeedItem) NeedItemResponse {
	if item == nil {
		return NeedItemResponse{}
	}

	var selectedCostID *string
	if item.SelectedCostID != nil {
		s := item.SelectedCostID.String()
		selectedCostID = &s
	}

	return NeedItemResponse{
		ID:              item.ID.String(),
		NeedID:          item.NeedID.String(),
		MaterialID:      item.MaterialID.String(),
		MaterialName:    item.Material.Name,
		MaterialERPCode: item.Material.ERPCode,
		Quantity:        item.Quantity,
		Unit:            item.Unit,
		SelectedCostID:  selectedCostID,
		Notes:           item.Notes,
		CreatedAt:       item.CreatedAt,
		UpdatedAt:       item.UpdatedAt,
	}
}

// ToNeedItemListResponse transforms a list of NeedItem entities.
func ToNeedItemListResponse(list []entities.NeedItem) []NeedItemResponse {
	response := make([]NeedItemResponse, 0, len(list))
	for _, item := range list {
		i := item
		response = append(response, ToNeedItemResponse(&i))
	}
	return response
}
