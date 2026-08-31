package sync_purchase_order

import (
	"github.com/google/uuid"
	"github.com/unknowncode44/unknown-go-server/pkg/entities"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	BatchCreate(orders []entities.PurchaseOrderSync) error
	FindByID(id uint) (*entities.PurchaseOrderSync, error)
	FindAll() ([]entities.PurchaseOrderSync, error)
	UpdateLinks(po *entities.PurchaseOrderSync, materialID *uuid.UUID, vendorID *uuid.UUID) (*entities.PurchaseOrderSync, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) Repository {
	return &repository{db: db}
}

// BatchCreate upserts by the (oc_bejerman, articulo) natural key — the VBA
// macro resends whole date ranges on every run (no local "already sent"
// tracking), so the same line is retransmitted whenever ranges overlap.
// On conflict it refreshes only the fields the Excel actually carries;
// material_id/vendor_id are deliberately left out so a resync never wipes
// out links assigned by hand from the Purchase Orders view.
func (r *repository) BatchCreate(orders []entities.PurchaseOrderSync) error {
	// GORM devuelve "empty slice found" si se le pasa un batch vacio.
	if len(orders) == 0 {
		return nil
	}

	return r.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "oc_bejerman"}, {Name: "articulo"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"cco", "nota_pedido", "fecha_np", "fecha_oc",
			"descripcion", "cantidad", "proveedor", "moneda",
			"importe_unitario", "sincronizado_en",
		}),
	}).Create(&orders).Error
}

// FindAll retrieves all synced purchase orders
func (r *repository) FindAll() ([]entities.PurchaseOrderSync, error) {
	var orders []entities.PurchaseOrderSync
	// Return ordered by latest synchronization
	err := r.db.Order("fecha_oc DESC").Find(&orders).Error
	return orders, err
}

func (r *repository) FindByID(id uint) (*entities.PurchaseOrderSync, error) {
	var po entities.PurchaseOrderSync
	if err := r.db.First(&po, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &po, nil
}

func (r *repository) UpdateLinks(po *entities.PurchaseOrderSync, materialID *uuid.UUID, vendorID *uuid.UUID) (*entities.PurchaseOrderSync, error) {
	if err := r.db.
		Model(&entities.PurchaseOrderSync{}).
		Where("id = ?", po.ID).
		Updates(map[string]interface{}{
			"material_id": materialID,
			"vendor_id":   vendorID,
		}).Error; err != nil {
		return nil, err
	}

	var mPo entities.PurchaseOrderSync

	if err := r.db.First(&mPo, "id = ?", po.ID).Error; err != nil {
		return nil, err
	}

	return &mPo, nil
}
