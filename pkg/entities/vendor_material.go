package entities

import (
	"time"

	"github.com/google/uuid"
)

type VendorMaterial struct {
	ID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`

	VendorID   uuid.UUID `gorm:"type:uuid;not null"`
	MaterialID uuid.UUID `gorm:"type:uuid;not null"`

	// Vendor id code (SKU, part number, etc.)
	VendorCode string `gorm:"not null"`

	IsActive bool `gorm:"default:true"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
