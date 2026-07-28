package entities

import (
	"time"

	"github.com/google/uuid"
)

type PurchaseOrderSync struct {
	ID         uint      `gorm:"primaryKey"`
	CCO        string    `gorm:"size:50"`
	NotaPedido string    `gorm:"size:100"`
	FechaNP    string    `gorm:"size:20"`
	FechaOC    time.Time `gorm:"type:date"`
	// OCBejerman + Articulo identifican una línea de OC: el número de OC se
	// repite una vez por cada artículo que incluye (una OC puede tener varios
	// materiales), así que ninguno de los dos solo alcanza como clave natural.
	OCBejerman      string     `gorm:"size:100;uniqueIndex:idx_po_natural_key,priority:1"`
	Articulo        string     `gorm:"size:100;uniqueIndex:idx_po_natural_key,priority:2"`
	Descripcion     string     `gorm:"type:text"`
	Cantidad        float64    `gorm:"type:decimal(19,4)"`
	Proveedor       string     `gorm:"size:255"`
	Moneda          string     `gorm:"size:10"`
	ImporteUnitario float64    `gorm:"type:decimal(19,4)"`
	SincronizadoEn  time.Time  `gorm:"autoCreateTime"`
	MaterialID      *uuid.UUID `gorm:"type:uuid;index"`
	VendorID        *uuid.UUID `gorm:"type:uuid;index"`
}
