package entities

import (
	"time"

	"github.com/google/uuid"
)

type AssetStatus string

const (
	AssetStatusInStock   AssetStatus = "IN_STOCK"
	AssetStatusInstalled AssetStatus = "INSTALLED"
	AssetStatusInRepair  AssetStatus = "IN_REPAIR"
	AssetStatusRetired   AssetStatus = "RETIRED"
)

type Asset struct {
	ID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`

	// Identificador visible y público (QR)
	SerialVisible string `gorm:"column:serial_visible;type:varchar(64);uniqueIndex;not null"`

	// Relación con material
	MaterialID uuid.UUID `gorm:"type:uuid;not null"`
	Material   Material  `gorm:"foreignKey:MaterialID"`

	// Serial del fabricante (opcional)
	ManufacturerSerial *string `gorm:"column:manufacturer_serial;type:varchar(128)"`

	// Estado operativo actual
	Status AssetStatus `gorm:"type:varchar(32);not null"`

	// Kits / jerarquía de assets
	ParentAssetID *uuid.UUID `gorm:"type:uuid"`
	ParentAsset   *Asset     `gorm:"foreignKey:ParentAssetID"`
	ChildAssets   []Asset    `gorm:"foreignKey:ParentAssetID"`

	// Ubicación actual (estado derivado)
	CurrentLocationID *uuid.UUID `gorm:"type:uuid"`
	CurrentLocation   *Location  `gorm:"foreignKey:CurrentLocationID"`

	// Control
	IsActive bool `gorm:"default:true"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
