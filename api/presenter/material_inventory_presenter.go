package presenter

import (
	"time"

	"github.com/unknowncode44/unknown-go-server/pkg/entities"
)

type CreateMaterialInventoryRequest struct {
	MaterialID      string  `json:"material_id"`
	InitialQuantity float64 `json:"initial_quantity"`
	StartDate       string  `json:"start_date"`
}

type UpdateMaterialInventoryRequest struct {
	InitialQuantity *float64 `json:"initial_quantity"`
	StartDate       string   `json:"start_date"`
}

type MaterialInventoryResponse struct {
	ID              string    `json:"id"`
	MaterialID      string    `json:"material_id"`
	MaterialName    string    `json:"material_name"`
	MaterialCode    string    `json:"material_code"`
	InitialQuantity float64   `json:"initial_quantity"`
	StartDate       time.Time `json:"start_date"`
	IsActive        bool      `json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type MaterialInventorySuccessResponse struct {
	Success bool        `json:"ok"`
	Data    interface{} `json:"data"`
}

type MaterialInventoryErrorResponse struct {
	Success bool   `json:"ok"`
	Error   string `json:"error"`
}

func ToMaterialInventoryResponse(inv *entities.MaterialInventory) MaterialInventoryResponse {
	if inv == nil {
		return MaterialInventoryResponse{}
	}
	return MaterialInventoryResponse{
		ID:              inv.ID.String(),
		MaterialID:      inv.MaterialID.String(),
		MaterialName:    inv.Material.Name,
		MaterialCode:    inv.Material.Code,
		InitialQuantity: inv.InitialQuantity,
		StartDate:       inv.StartDate,
		IsActive:        inv.IsActive,
		CreatedAt:       inv.CreatedAt,
		UpdatedAt:       inv.UpdatedAt,
	}
}

func ToMaterialInventoryListResponse(list []entities.MaterialInventory) []MaterialInventoryResponse {
	resp := make([]MaterialInventoryResponse, 0, len(list))
	for _, inv := range list {
		resp = append(resp, ToMaterialInventoryResponse(&inv))
	}
	return resp
}
