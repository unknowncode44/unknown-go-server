package presenter

import (
	"time"

	"github.com/unknowncode44/unknown-go-server/pkg/entities"
)

type CreateMaterialCostRequest struct {
	MaterialID string  `json:"material_id"`
	VendorID   string  `json:"vendor_id"`
	CurrencyID string  `json:"currency_id"`
	Cost       float64 `json:"cost"`
	CostDate   string  `json:"cost_date"`
}

type MaterialCostResponse struct {
	ID         string    `json:"id"`
	MaterialID string    `json:"material_id"`
	VendorID   string    `json:"vendor_id"`
	CurrencyID string    `json:"currency_id"`
	Cost       float64   `json:"cost"`
	CostDate   time.Time `json:"cost_date"`
	CreatedAt  time.Time `json:"created_at"`
}

type MaterialCostSuccessResponse struct {
	Success bool        `json:"ok"`
	Data    interface{} `json:"data"`
}

type MaterialCostErrorResponse struct {
	Success bool   `json:"ok"`
	Error   string `json:"error"`
}

func ToMaterialCostResponse(mc *entities.MaterialCost) MaterialCostResponse {
	if mc == nil {
		return MaterialCostResponse{}
	}

	return MaterialCostResponse{
		ID:         mc.ID.String(),
		MaterialID: mc.MaterialID.String(),
		VendorID:   mc.VendorID.String(),
		CurrencyID: mc.CurrencyID.String(),
		Cost:       mc.Cost,
		CostDate:   mc.CostDate,
		CreatedAt:  mc.CreatedAt,
	}
}

func ToMaterialCostListResponse(list []entities.MaterialCost) []MaterialCostResponse {
	resp := make([]MaterialCostResponse, 0, len(list))
	for _, m := range list {
		resp = append(resp, ToMaterialCostResponse(&m))
	}
	return resp
}
