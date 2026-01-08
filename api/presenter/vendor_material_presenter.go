package presenter

import (
	"time"

	"github.com/unknowncode44/unknown-go-server/pkg/entities"
)

type CreateVendorMaterialRequest struct {
	VendorID   string `json:"vendor_id"`
	MaterialID string `json:"material_id"`
	VendorCode string `json:"vendor_code"`
}

type UpdateVendorMaterialRequest struct {
	VendorCode string `json:"vendor_code,omitempty"`
}

type VendorMaterialResponse struct {
	ID         string    `json:"id"`
	VendorID   string    `json:"vendor_id"`
	MaterialID string    `json:"material_id"`
	VendorCode string    `json:"vendor_code"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type VendorMaterialSuccessResponse struct {
	Success bool        `json:"ok"`
	Data    interface{} `json:"data"`
}

type VendorMaterialErrorResponse struct {
	Success bool   `json:"ok"`
	Error   string `json:"error"`
}

func ToVendorMaterialResponse(vm *entities.VendorMaterial) VendorMaterialResponse {
	if vm == nil {
		return VendorMaterialResponse{}
	}
	return VendorMaterialResponse{
		ID:         vm.ID.String(),
		VendorID:   vm.VendorID.String(),
		MaterialID: vm.MaterialID.String(),
		VendorCode: vm.VendorCode,
		IsActive:   vm.IsActive,
		CreatedAt:  vm.CreatedAt,
		UpdatedAt:  vm.UpdatedAt,
	}
}

func ToVendorMaterialListResponse(vms []entities.VendorMaterial) []VendorMaterialResponse {
	resp := make([]VendorMaterialResponse, 0, len(vms))
	for _, vm := range vms {
		resp = append(resp, ToVendorMaterialResponse(&vm))
	}
	return resp
}
