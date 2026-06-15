package entities

import (
	"time"

	"github.com/google/uuid"
)

// NeedCounter stores the last correlative used per year to generate the
// Number field of Need (e.g. NEED-2026-0001) safely under transaction locking.
type NeedCounter struct {
	ID   uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Year int       `gorm:"not null;uniqueIndex"`
	Last int       `gorm:"not null;default:0"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
