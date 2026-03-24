package delivery_record

import (
	"github.com/google/uuid"
	"github.com/unknowncode44/unknown-go-server/pkg/entities"
	"gorm.io/gorm"
)

// deliverySum holds the result of the SUM query grouped by type.
type deliverySum struct {
	Type  entities.DeliveryType
	Total float64
}

// Repository defines persistence operations for DeliveryRecord entities.
type Repository interface {
	Create(dr *entities.DeliveryRecord) (*entities.DeliveryRecord, error)
	FindByInventoryID(inventoryID uuid.UUID) ([]entities.DeliveryRecord, error)
	FindAll() ([]entities.DeliveryRecord, error)
	SumByInventoryID(inventoryID uuid.UUID) (inbound float64, outbound float64, err error)
}

type repository struct {
	db *gorm.DB
}

// NewRepo returns a new Repository using the provided GORM DB.
func NewRepo(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(dr *entities.DeliveryRecord) (*entities.DeliveryRecord, error) {
	if err := r.db.Create(dr).Error; err != nil {
		return nil, err
	}
	return dr, nil
}

func (r *repository) FindByInventoryID(inventoryID uuid.UUID) ([]entities.DeliveryRecord, error) {
	var list []entities.DeliveryRecord
	if err := r.db.Where("material_inventory_id = ?", inventoryID).
		Order("delivery_date desc, created_at desc").
		Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *repository) FindAll() ([]entities.DeliveryRecord, error) {
	var list []entities.DeliveryRecord
	if err := r.db.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *repository) SumByInventoryID(inventoryID uuid.UUID) (inbound float64, outbound float64, err error) {
	var results []deliverySum
	err = r.db.Model(&entities.DeliveryRecord{}).
		Select("type, COALESCE(SUM(quantity), 0) as total").
		Where("material_inventory_id = ?", inventoryID).
		Group("type").
		Scan(&results).Error
	if err != nil {
		return 0, 0, err
	}

	for _, r := range results {
		switch r.Type {
		case entities.DeliveryInbound:
			inbound = r.Total
		case entities.DeliveryOutbound:
			outbound = r.Total
		}
	}
	return inbound, outbound, nil
}
