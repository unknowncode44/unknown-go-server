package entities

import (
	"time"

	"github.com/google/uuid"
)

// NeedItem is a line of a Need: one material, a quantity, and optionally a
// MaterialCost reference that fixes the chosen price for that line (among
// the available vendor quotations).
type NeedItem struct {
	ID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`

	NeedID uuid.UUID `gorm:"type:uuid;not null;index"`
	Need   Need      `gorm:"foreignKey:NeedID"`

	MaterialID uuid.UUID `gorm:"type:uuid;not null;index"`
	Material   Material  `gorm:"foreignKey:MaterialID"`

	Quantity float64 `gorm:"type:numeric(15,4);not null"`
	Unit     string  `gorm:"type:varchar(50);not null"`

	// Reference price chosen by the buyer. May be nil while quotations are
	// still being collected. Points to an already-registered MaterialCost.
	SelectedCostID *uuid.UUID    `gorm:"type:uuid;index"`
	SelectedCost   *MaterialCost `gorm:"foreignKey:SelectedCostID"`

	Notes string `gorm:"type:text"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
