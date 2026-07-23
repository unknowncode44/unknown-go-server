package entities

import (
	"time"

	"github.com/google/uuid"
)

type ShipmentStatus string

const (
	ShipmentStatusDraft      ShipmentStatus = "DRAFT"
	ShipmentStatusProcessed  ShipmentStatus = "PROCESSED"
	ShipmentStatusReconciled ShipmentStatus = "RECONCILED"
)

type ShipmentItemType string

const (
	ShipmentItemTypeAsset    ShipmentItemType = "ASSET"
	ShipmentItemTypeMaterial ShipmentItemType = "MATERIAL"
)

type Shipment struct {
	ID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`

	DestinationLocationID uuid.UUID `gorm:"type:uuid;not null"`
	DestinationLocation   Location  `gorm:"foreignKey:DestinationLocationID"`

	Status ShipmentStatus `gorm:"type:varchar(32);not null;default:'DRAFT'"`

	CreatedByUserID uuid.UUID `gorm:"type:uuid;not null"`
	CreatedByUser   User      `gorm:"foreignKey:CreatedByUserID"`

	// Quién solicita la entrega — campo de texto libre y obligatorio (mismo
	// criterio que ya usa entities.Need.RequesterName: Quinar todavía no
	// tiene un dominio de "solicitantes" propio, así que no es una FK).
	RequesterName string `gorm:"type:varchar(255);not null"`

	// Descripción libre de la entrega — opcional.
	Notes string `gorm:"type:text"`

	// Se completa recién en la reconciliación de escritorio.
	ERPMovementNumber *string `gorm:"type:varchar(64)"`

	ProcessedAt  *time.Time
	ReconciledAt *time.Time

	Items []ShipmentItem `gorm:"foreignKey:ShipmentID"`

	IsActive  bool `gorm:"default:true;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ShipmentItem struct {
	ID         uuid.UUID        `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	ShipmentID uuid.UUID        `gorm:"type:uuid;not null"`
	Type       ShipmentItemType `gorm:"type:varchar(16);not null"`

	// Uno de los dos según Type — nunca ambos.
	AssetID *uuid.UUID `gorm:"type:uuid"`
	Asset   *Asset     `gorm:"foreignKey:AssetID"`

	MaterialID *uuid.UUID `gorm:"type:uuid"`
	Material   *Material  `gorm:"foreignKey:MaterialID"`
	Quantity   *float64   `gorm:"type:numeric(15,4)"` // solo para MATERIAL

	// Se completan al procesar el envío, para trazabilidad.
	ResultingAssetMovementID  *uuid.UUID `gorm:"type:uuid"`
	ResultingDeliveryRecordID *uuid.UUID `gorm:"type:uuid"`

	CreatedAt time.Time
}
