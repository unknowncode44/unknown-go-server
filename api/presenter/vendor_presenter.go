package presenter

import (
	"github.com/unknowncode44/unknown-go-server/pkg/entities"
)

// CreateVendorRequest struct representa el cuerpo esperado para crear un vendor
// Las validaciones adicionales se aplican en la capa de servicio

type CreateVendorRequest struct {
	Name  string `json:"name"`
	Code  string `json:"code,omitempty"`
	TaxID string `json:"tax_id,omitempty"`
}

// UpdateVendorRequest define el cuerpo esperado para actualizar un proveedor.
// El ID viene por path, no en el body.
type UpdateVendorRequest struct {
	Name  string `json:"name,omitempty"`
	Code  string `json:"code,omitempty"`
	TaxID string `json:"tax_id,omitempty"`
}

// VendorResponse representa un proveedor en las respuestas HTTP.
type VendorResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Code     string `json:"code,omitempty"`
	TaxID    string `json:"tax_id,omitempty"`
	IsActive bool   `json:"is_active"`
}

// VendorSuccessResponse representa una respuesta exitosa con uno o más vendors.
type VendorSuccessResponse struct {
	Success bool        `json:"ok"`
	Data    interface{} `json:"data"`
}

// VendorErrorResponse representa una respuesta de error estándar.
type VendorErrorResponse struct {
	Success bool   `json:"ok"`
	Error   string `json:"error"`
}

// ToVendorResponse transforma una entidad Vendor en VendorResponse.
func ToVendorResponse(v *entities.Vendor) *VendorResponse {
	if v == nil {
		return nil
	}

	return &VendorResponse{
		ID:       v.ID.String(),
		Name:     v.Name,
		Code:     v.Code,
		TaxID:    v.TaxID,
		IsActive: v.IsActive,
	}
}

// ToVendorListResponse transforma una lista de entidades Vendor.
func ToVendorListResponse(vendors []entities.Vendor) []VendorResponse {
	response := make([]VendorResponse, 0, len(vendors))

	for _, v := range vendors {
		response = append(response, VendorResponse{
			ID:       v.ID.String(),
			Name:     v.Name,
			Code:     v.Code,
			TaxID:    v.TaxID,
			IsActive: v.IsActive,
		})
	}

	return response
}
