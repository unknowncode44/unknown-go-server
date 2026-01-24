package entities

import (
	"time"

	"github.com/google/uuid"
)

type LocationType string

const (
	LocationTypeInternal LocationType = "INTERNAL"
	LocationTypeClient   LocationType = "CLIENT"
)

type Location struct {
	ID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`

	// Nombre visible de la ubicación
	Name string `gorm:"type:varchar(128);not null"`

	// Tipo de ubicación (interna / cliente)
	Type LocationType `gorm:"type:varchar(32);not null"`

	// Jerarquía de ubicaciones
	ParentLocationID *uuid.UUID `gorm:"type:uuid"`
	ParentLocation   *Location  `gorm:"foreignKey:ParentLocationID"`
	ChildLocations   []Location `gorm:"foreignKey:ParentLocationID"`

	// Relación opcional con cliente (futuro)
	// ClientID *uuid.UUID

	IsActive bool `gorm:"default:true"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
