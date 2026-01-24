package entities

import (
	"time"

	"github.com/google/uuid"
)

// AssetSerialCounter stores the last used correlative per material to
// generate serial_visible values safely under transaction / locking.
type AssetSerialCounter struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	MaterialID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"`
	Last       int       `gorm:"not null;default:0"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
