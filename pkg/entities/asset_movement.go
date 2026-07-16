package entities

import (
	"time"

	"github.com/google/uuid"
)

type AssetMovementType string

const (
	AssetMovementInbound      AssetMovementType = "INBOUND"      // ingreso a stock
	AssetMovementTransfer     AssetMovementType = "TRANSFER"     // traslado interno
	AssetMovementInstall      AssetMovementType = "INSTALL"      // instalación
	AssetMovementUninstall    AssetMovementType = "UNINSTALL"    // retiro
	AssetMovementRepair       AssetMovementType = "REPAIR"       // envío a reparación
	AssetMovementReturn       AssetMovementType = "RETURN"       // regreso de reparación
	AssetMovementDecommission AssetMovementType = "DECOMMISSION" // baja definitiva
	AssetMovementScrap        AssetMovementType = "SCRAP"        // baja definitiva por rotura o obsolescencia
	AssetMovementSold         AssetMovementType = "SOLD"         // venta de BDC serializado a cliente final
)

type AssetMovement struct {
	ID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`

	// Asset afectado
	AssetID uuid.UUID `gorm:"type:uuid;not null"`
	Asset   Asset     `gorm:"foreignKey:AssetID"`

	// Tipo de movimiento
	Type AssetMovementType `gorm:"type:varchar(32);not null"`

	// Ubicación origen y destino
	FromLocationID *uuid.UUID `gorm:"type:uuid"`
	FromLocation   *Location  `gorm:"foreignKey:FromLocationID"`

	ToLocationID *uuid.UUID `gorm:"type:uuid"`
	ToLocation   *Location  `gorm:"foreignKey:ToLocationID"`

	// Fecha efectiva del movimiento
	MovementDate time.Time `gorm:"not null"`

	// Observaciones / motivo
	Notes *string `gorm:"type:text"`

	CreatedAt time.Time
}
