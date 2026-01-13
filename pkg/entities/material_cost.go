package entities

import (
	"time"

	"github.com/google/uuid"
)

// MaterialCost represents a historical cost of a material
// offered by a specific supplier.
// It is an append-only entity: costs are not updated; instead, new records are created.
type MaterialCost struct {
	ID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`

	// Mandatory relationships
	MaterialID uuid.UUID `gorm:"type:uuid;not null;index"`
	VendorID   uuid.UUID `gorm:"type:uuid;not null;index"`
	CurrencyID uuid.UUID `gorm:"type:uuid;not null;index"`

	// Cost value
	Cost float64 `gorm:"type:numeric(15,4);not null"`

	// Date from which this cost is valid.
	CostDate time.Time `gorm:"not null;index"`

	CreatedAt time.Time
}
