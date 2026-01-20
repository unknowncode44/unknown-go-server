package material_cost

import (
	"github.com/google/uuid"
	"github.com/unknowncode44/unknown-go-server/pkg/entities"
	"gorm.io/gorm"
)

// Repository defines persistence operations for MaterialCost entities.
type Repository interface {
	Create(mc *entities.MaterialCost) (*entities.MaterialCost, error)
	FindAll() ([]entities.MaterialCost, error)
	FindByID(id uuid.UUID) (*entities.MaterialCost, error)
	FindByMaterial(materialID uuid.UUID) ([]entities.MaterialCost, error)
	FindLatestByMaterialAndVendor(materialID uuid.UUID, vendorID uuid.UUID) (*entities.MaterialCost, error)
}

type repository struct {
	db *gorm.DB
}

// NewRepo returns a new GORM-backed Repository.
func NewRepo(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(mc *entities.MaterialCost) (*entities.MaterialCost, error) {
	if err := r.db.Create(mc).Error; err != nil {
		return nil, err
	}
	return mc, nil
}

func (r *repository) FindAll() ([]entities.MaterialCost, error) {
	var list []entities.MaterialCost
	if err := r.db.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *repository) FindByID(id uuid.UUID) (*entities.MaterialCost, error) {
	var mc entities.MaterialCost
	if err := r.db.First(&mc, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &mc, nil
}

func (r *repository) FindByMaterial(materialID uuid.UUID) ([]entities.MaterialCost, error) {
	var list []entities.MaterialCost
	if err := r.db.Where("material_id = ?", materialID).Order("cost_date desc, created_at desc").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *repository) FindLatestByMaterialAndVendor(materialID uuid.UUID, vendorID uuid.UUID) (*entities.MaterialCost, error) {
	var mc entities.MaterialCost
	if err := r.db.Where("material_id = ? AND vendor_id = ?", materialID, vendorID).
		Order("cost_date desc, created_at desc").
		First(&mc).Error; err != nil {
		return nil, err
	}
	return &mc, nil
}
