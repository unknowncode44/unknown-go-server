package material_inventory

import (
	"github.com/google/uuid"
	"github.com/unknowncode44/unknown-go-server/pkg/entities"
	"gorm.io/gorm"
)

// Repository defines persistence operations for MaterialInventory entities.
type Repository interface {
	Create(inv *entities.MaterialInventory) (*entities.MaterialInventory, error)
	FindAll() ([]entities.MaterialInventory, error)
	FindByID(id uuid.UUID) (*entities.MaterialInventory, error)
	FindByMaterialID(materialID uuid.UUID) (*entities.MaterialInventory, error)
	Update(inv *entities.MaterialInventory) (*entities.MaterialInventory, error)
	Deactivate(id uuid.UUID) error
}

type repository struct {
	db *gorm.DB
}

// NewRepo returns a new Repository using the provided GORM DB.
func NewRepo(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(inv *entities.MaterialInventory) (*entities.MaterialInventory, error) {
	if err := r.db.Create(inv).Error; err != nil {
		return nil, err
	}
	return inv, nil
}

func (r *repository) FindAll() ([]entities.MaterialInventory, error) {
	var list []entities.MaterialInventory
	if err := r.db.Preload("Material").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *repository) FindByID(id uuid.UUID) (*entities.MaterialInventory, error) {
	var inv entities.MaterialInventory
	if err := r.db.Preload("Material").First(&inv, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &inv, nil
}

func (r *repository) FindByMaterialID(materialID uuid.UUID) (*entities.MaterialInventory, error) {
	var inv entities.MaterialInventory
	if err := r.db.Preload("Material").First(&inv, "material_id = ?", materialID).Error; err != nil {
		return nil, err
	}
	return &inv, nil
}

func (r *repository) Update(inv *entities.MaterialInventory) (*entities.MaterialInventory, error) {
	if err := r.db.
		Model(&entities.MaterialInventory{}).
		Where("id = ?", inv.ID).
		Updates(map[string]interface{}{
			"initial_quantity": inv.InitialQuantity,
			"start_date":      inv.StartDate,
			"is_active":       inv.IsActive,
		}).Error; err != nil {
		return nil, err
	}
	return inv, nil
}

func (r *repository) Deactivate(id uuid.UUID) error {
	return r.db.Model(&entities.MaterialInventory{}).
		Where("id = ?", id).
		Update("is_active", false).Error
}
