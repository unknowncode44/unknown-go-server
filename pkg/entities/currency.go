package entities

import (
	"time"

	"github.com/google/uuid"
)

// Currency represents a currency supported by the system
// Is a reference table (master data)
type Currency struct {
	ID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`

	Code string `gorm:"type:varchar(3);not null;uniqueIndex"` // ISO 4217 (ARS, USD, EUR)
	Name string `gorm:"type:varchar(100);not null"`

	IsActive bool `gorm:"default:true"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
