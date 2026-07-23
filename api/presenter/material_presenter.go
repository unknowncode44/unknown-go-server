package presenter

import (
	"time"

	"github.com/unknowncode44/unknown-go-server/pkg/entities"
)

// CreateMaterialRequest representa el cuerpo esperado para crear un material.
// Los validadores adicionales se aplican en la capa de servicio.
type CreateMaterialRequest struct {
	Name          string `json:"name"`
	Sector        string `json:"sector"`
	Group         string `json:"group"`
	UnitOfMeasure string `json:"unitOfMeasure"`
	ERPCode       string `json:"erpCode"`
	InternalCode  string `json:"internalCode"`
	Code          string `json:"code"`
}

// UpdateMaterialRequest representa los campos permitidos para actualizar un material.
// Notar que `IsActive` permite reactivar/desactivar desde la API.
type UpdateMaterialRequest struct {
	Name          string `json:"name"`
	Sector        string `json:"sector"`
	Group         string `json:"group"`
	UnitOfMeasure string `json:"unitOfMeasure"`
	ERPCode       string `json:"erpCode"`
	InternalCode  string `json:"internalCode"`
	Code          string `json:"code"`
}

// MaterialResponse define la respuesta que devuelve la API para un material.
type MaterialResponse struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Sector        string    `json:"sector"`
	Group         string    `json:"group"`
	ERPCode       string    `json:"erpCode"`
	InternalCode  string    `json:"internalCode"`
	UnitOfMeasure string    `json:"unitOfMeasure"`
	Code          string    `json:"code"`
	IsActive      bool      `json:"isActive"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

// PublicMaterialResponse es la vista pública (sin autenticación) que se
// devuelve al escanear el QR de un material a granel (BDC).
type PublicMaterialResponse struct {
	Code          string   `json:"code"`
	Name          string   `json:"name"`
	Sector        string   `json:"sector"`
	UnitOfMeasure string   `json:"unit_of_measure"`
	Group         string   `json:"group"`
	CurrentStock  *float64 `json:"current_stock,omitempty"`
}

// MaterialSuccessResponse representa una respuesta exitosa con uno o más vendors.
type MaterialSuccessResponse struct {
	Success bool        `json:"ok"`
	Data    interface{} `json:"data"`
}

// MaterialErrorResponse representa una respuesta de error estándar.
type MaterialErrorResponse struct {
	Success bool   `json:"ok"`
	Error   string `json:"error"`
}

// ToMaterialResponse transforma una entidad Material en MaterialResponse.
func ToMaterialResponse(m *entities.Material) MaterialResponse {
	if m == nil {
		return MaterialResponse{}
	}

	return MaterialResponse{
		ID:            m.ID.String(),
		Name:          m.Name,
		Sector:        m.Sector,
		Group:         string(m.Group),
		ERPCode:       m.ERPCode,
		InternalCode:  ptrToString(m.InternalCode),
		UnitOfMeasure: m.UnitOfMeasure,
		Code:          m.Code,
		IsActive:      m.IsActive,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}

// ToMaterialListResponse transforma una lista de entidades Materials.
func ToMaterialListResponse(vendors []entities.Material) []MaterialResponse {
	response := make([]MaterialResponse, 0, len(vendors))

	for _, m := range vendors {
		response = append(response, MaterialResponse{
			ID:            m.ID.String(),
			Name:          m.Name,
			Sector:        m.Sector,
			Group:         string(m.Group),
			ERPCode:       m.ERPCode,
			InternalCode:  ptrToString(m.InternalCode),
			UnitOfMeasure: m.UnitOfMeasure,
			Code:          m.Code,
			IsActive:      m.IsActive,
			CreatedAt:     m.CreatedAt,
			UpdatedAt:     m.UpdatedAt,
		})
	}

	return response
}
