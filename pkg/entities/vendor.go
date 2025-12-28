package entities

import (
	"time"

	"github.com/google/uuid"
)

type Vendor struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name      string    `gorm:"size:255;not null"`
	Code      string    `gorm:"size:100;uniqueIndex;not null"`
	IsActive  bool      `gorm:"default:true;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
