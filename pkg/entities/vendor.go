package entities

import (
	"time"

	"github.com/google/uuid"
)

type Vendor struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name      string    `gorm:"type:varchar(255);not null"`
	Code      string    `gorm:"type:varchar(100)"`
	TaxID     string    `gorm:"type:varchar(50)"` // opcional
	IsActive  bool      `gorm:"not null;default:true"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
