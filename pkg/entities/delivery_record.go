package entities

import (
	"time"

	"github.com/google/uuid"
)

type DeliveryType string

const (
	DeliveryInbound  DeliveryType = "INBOUND"
	DeliveryOutbound DeliveryType = "OUTBOUND"
)

type DeliveryRecord struct {
	ID                  uuid.UUID         `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	MaterialInventoryID uuid.UUID         `gorm:"type:uuid;not null;index"`
	MaterialInventory   MaterialInventory `gorm:"foreignKey:MaterialInventoryID"`
	Type                DeliveryType      `gorm:"type:varchar(16);not null"`
	Quantity            float64           `gorm:"type:numeric(15,4);not null"`
	DeliveryDate        time.Time         `gorm:"not null;index"`
	Notes               *string           `gorm:"type:text"`
	CreatedAt           time.Time
}
