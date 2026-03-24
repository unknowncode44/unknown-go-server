package entities

import (
	"time"

	"github.com/google/uuid"
)

type MaterialInventory struct {
	ID              uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	MaterialID      uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"`
	Material        Material  `gorm:"foreignKey:MaterialID"`
	InitialQuantity float64   `gorm:"type:numeric(15,4);not null;default:0"`
	StartDate       time.Time `gorm:"not null"`
	IsActive        bool      `gorm:"default:true;not null"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
