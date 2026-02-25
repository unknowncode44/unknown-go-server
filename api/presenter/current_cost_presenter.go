package presenter

import "time"

type CurrentCostData struct {
	MaterialID string    `json:"material_id"`
	VendorID   string    `json:"vendor_id"`
	Currency   string    `json:"currency"`
	Cost       float64   `json:"cost"`
	CostDate   time.Time `json:"cost_date"`
}

type CurrentCostSuccessResponse struct {
	Success bool        `json:"ok"`
	Data    interface{} `json:"data"`
}

type CurrentCostErrorResponse struct {
	Success bool   `json:"ok"`
	Error   string `json:"error"`
}

type VendorComparisonItem struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Currency string    `json:"currency"`
	Cost     float64   `json:"cost"`
	CostDate time.Time `json:"cost_date"`
}

type VendorComparisonResponse struct {
	MaterialID string                 `json:"material_id"`
	At         string                 `json:"at"`
	Vendors    []VendorComparisonItem `json:"vendors"`
}

func ToCurrentCostData(materialID, vendorID string, currency string, cost float64, costDate time.Time) CurrentCostData {
	return CurrentCostData{
		MaterialID: materialID,
		VendorID:   vendorID,
		Currency:   currency,
		Cost:       cost,
		CostDate:   costDate,
	}
}

func ToVendorComparisonItem(id, name, currency string, cost float64, costDate time.Time) VendorComparisonItem {
	return VendorComparisonItem{
		ID:       id,
		Name:     name,
		Currency: currency,
		Cost:     cost,
		CostDate: costDate,
	}
}
