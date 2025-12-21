package presenter

import (
	"time"

	"github.com/google/uuid"
)

type CreateMaterialRequest struct {
	Name          string `json:"name"`
	Sector        string `json:"sector"`
	UnitOfMeasure string `json:"unitOfMeasure"`
}

type UpdateMaterialRequest struct {
	Name          string `json:"name"`
	Sector        string `json:"sector"`
	UnitOfMeasure string `json:"unitOfMeasure"`
	IsActive      bool   `json:"isActive"`
}

type MaterialResponse struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	Sector        string    `json:"sector"`
	UnitOfMeasure string    `json:"unitOfMeasure"`
	IsActive      bool      `json:"isActive"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}
