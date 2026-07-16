package entities

import (
	"time"

	"github.com/google/uuid"
)

// MaterialGroup clasifica el material según su grupo de negocio,
// independiente de cómo se trackea el stock (granel o serializado).
type MaterialGroup string

const (
	MaterialGroupBDC MaterialGroup = "BDC" // Bienes de Cambio: materiales para vender/instalar en clientes
	MaterialGroupBDU MaterialGroup = "BDU" // Bienes de Uso: activos de infraestructura propia
)

// Material represents a material resource used by the application.
//
// The struct includes persistence metadata for the ORM (GORM) and
// basic lifecycle fields. Field constraints (size, not-null, defaults)
// are defined via GORM struct tags.
type Material struct {
	ID     uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name   string    `gorm:"size:255;not null"`
	Sector string    `gorm:"size:100;not null"`
	// Columna explícita porque `group` es palabra reservada en SQL. El default
	// a nivel DB permite el ALTER TABLE sobre la tabla con filas existentes;
	// la validación real para altas nuevas vive en el servicio.
	Group         MaterialGroup `gorm:"column:material_group;type:varchar(10);not null;default:'BDC'"`
	UnitOfMeasure string        `gorm:"size:50;not null"`
	ERPCode       string        `gorm:"column:erp_code;index"`
	InternalCode  *string       `gorm:"column:internal_code;uniqueIndex"`
	Code          string        `gorm:"size:50;uniqueIndex;not null"`
	IsActive      bool          `gorm:"default:true;not null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
