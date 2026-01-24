package entities

import (
	"time"

	"github.com/google/uuid"
)

// Material represents a material resource used by the application.
//
// The struct includes persistence metadata for the ORM (GORM) and
// basic lifecycle fields. Field constraints (size, not-null, defaults)
// are defined via GORM struct tags.
type Material struct {
	ID            uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name          string    `gorm:"size:255;not null"`
	Sector        string    `gorm:"size:100;not null"`
	UnitOfMeasure string    `gorm:"size:50;not null"`
	Code          string    `gorm:"size:50;uniqueIndex;not null"`
	IsActive      bool      `gorm:"default:true;not null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
