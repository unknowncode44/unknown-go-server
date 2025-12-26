package presenter

import (
	"time"

	"github.com/google/uuid"
)

// CreateMaterialRequest representa el cuerpo esperado para crear un material.
// Los validadores adicionales se aplican en la capa de servicio.
type CreateMaterialRequest struct {
	Name          string `json:"name"`
	Sector        string `json:"sector"`
	UnitOfMeasure string `json:"unitOfMeasure"`
}

// UpdateMaterialRequest representa los campos permitidos para actualizar un material.
// Notar que `IsActive` permite reactivar/desactivar desde la API.
type UpdateMaterialRequest struct {
	Name          string `json:"name"`
	Sector        string `json:"sector"`
	UnitOfMeasure string `json:"unitOfMeasure"`
	IsActive      bool   `json:"isActive"`
}

// MaterialResponse define la respuesta que devuelve la API para un material.
type MaterialResponse struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	Sector        string    `json:"sector"`
	UnitOfMeasure string    `json:"unitOfMeasure"`
	IsActive      bool      `json:"isActive"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}
