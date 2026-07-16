package material

import (
	"github.com/google/uuid"
	"github.com/unknowncode44/unknown-go-server/pkg/entities"
	"gorm.io/gorm"
)

// Repository defines persistence operations for Material entities.
// The service layer depends on this interface to remain persistence-agnostic.
type Repository interface {
	Create(material *entities.Material) (*entities.Material, error)
	FindAll() ([]entities.Material, error)
	FindByID(id uuid.UUID) (*entities.Material, error)
	FindByCode(code string) (*entities.Material, error)
	FindByERPCode(erpCode string) ([]entities.Material, error)
	Update(material *entities.Material) (*entities.Material, error)
	Delete(id uuid.UUID) error
}

// repository is the GORM-backed Repository implementation.
type repository struct {
	db *gorm.DB
}

// NewRepo returns a new Repository using the provided GORM DB.
func NewRepo(db *gorm.DB) Repository {
	return &repository{db: db}
}

// Create inserts a new material record into the database.
func (r *repository) Create(material *entities.Material) (*entities.Material, error) {
	if err := r.db.Create(material).Error; err != nil {
		return nil, err
	}
	return material, nil
}

// FindAll retrieves all material records.
func (r *repository) FindAll() ([]entities.Material, error) {
	var materials []entities.Material
	if err := r.db.Find(&materials).Error; err != nil {
		return nil, err
	}
	return materials, nil
}

// FindByID returns a material by its UUID.
func (r *repository) FindByID(id uuid.UUID) (*entities.Material, error) {
	var material entities.Material
	if err := r.db.First(&material, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &material, nil
}

// FindByCode returns a material by its code.
func (r *repository) FindByCode(code string) (*entities.Material, error) {
	var material entities.Material
	if err := r.db.First(&material, "code = ?", code).Error; err != nil {
		return nil, err
	}
	return &material, nil
}

// FindByERPCode returns all materials that share the given ERP code. Some
// Bejerman "bag" codes (e.g. "0 MAT GOP21") group many distinct materials
// under the same ERP code, so this must return the full list.
func (r *repository) FindByERPCode(erpCode string) ([]entities.Material, error) {
	var materials []entities.Material
	if err := r.db.Where("erp_code = ?", erpCode).Find(&materials).Error; err != nil {
		return nil, err
	}
	return materials, nil
}

// Update applies changes to the material record and returns the updated entity.
// Uses an updates map so falsey values (e.g. IsActive=false) persist correctly.
func (r *repository) Update(material *entities.Material) (*entities.Material, error) {
	if err := r.db.
		Model(&entities.Material{}).
		Where("id = ?", material.ID).
		Updates(map[string]interface{}{
			"name":            material.Name,
			"sector":          material.Sector,
			"material_group":  material.Group,
			"unit_of_measure": material.UnitOfMeasure,
			"erp_code":        material.ERPCode,
			"code":            material.Code,
			"is_active":       material.IsActive,
		}).Error; err != nil {
		return nil, err
	}
	return material, nil
}

// Delete performs a logical delete by setting `is_active` to false.
func (r *repository) Delete(id uuid.UUID) error {
	return r.db.Model(&entities.Material{}).
		Where("id = ?", id).
		Update("is_active", false).Error
}
