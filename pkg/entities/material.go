package entities

import (
	"time"

	"github.com/google/uuid"
)

// Material representa el recurso de materiales en la aplicación.
// Contiene metadatos usados por GORM para persistencia.
type Material struct {
	ID            uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name          string    `gorm:"size:255;not null"`
	Sector        string    `gorm:"size:100;not null"`
	UnitOfMeasure string    `gorm:"size:50;not null"`
	IsActive      bool      `gorm:"default:true;not null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
