package sync_purchase_order

import (
	"github.com/unknowncode44/unknown-go-server/pkg/entities"
	"gorm.io/gorm"
)

type Repository interface {
	BatchCreate(orders []entities.PurchaseOrderSync) error
	FindAll() ([]entities.PurchaseOrderSync, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) BatchCreate(orders []entities.PurchaseOrderSync) error {
	return r.db.Create(&orders).Error
}

// FindAll retrieves all synced purchase orders
func (r *repository) FindAll() ([]entities.PurchaseOrderSync, error) {
	var orders []entities.PurchaseOrderSync
	// Return ordered by latest synchronization
	err := r.db.Order("fecha_oc DESC").Find(&orders).Error
	return orders, err
}
