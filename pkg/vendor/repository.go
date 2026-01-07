package vendor

import (
	"github.com/google/uuid"
	"github.com/unknowncode44/unknown-go-server/pkg/entities"
	"gorm.io/gorm"
)

// Repository defines persistence operations for Vendor entities.
// The service layer depends on this interface to remain persistence-agnostic.
type Repository interface {
	Create(vendor *entities.Vendor) (*entities.Vendor, error)
	FindAll() ([]entities.Vendor, error)
	FindByID(id uuid.UUID) (*entities.Vendor, error)
	Update(vendor *entities.Vendor) (*entities.Vendor, error)
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

// Create inserts a new vendor record into the database.
func (r *repository) Create(vendor *entities.Vendor) (*entities.Vendor, error) {
	if err := r.db.Create(vendor).Error; err != nil {
		return nil, err
	}
	return vendor, nil
}

// FindAll retrieves all vendor records.
func (r *repository) FindAll() ([]entities.Vendor, error) {
	var vendors []entities.Vendor
	if err := r.db.Find(&vendors).Error; err != nil {
		return nil, err
	}
	return vendors, nil
}

// FindByID returns a vendor by its UUID.
func (r *repository) FindByID(id uuid.UUID) (*entities.Vendor, error) {
	var vendor entities.Vendor
	if err := r.db.First(&vendor, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &vendor, nil
}

// Update applies changes to the vendor record and returns the updated entity.
// Uses an updates map so falsey values (e.g. IsActive=false) persist correctly.
func (r *repository) Update(vendor *entities.Vendor) (*entities.Vendor, error) {
	if err := r.db.
		Model(&entities.Vendor{}).
		Where("id = ?", vendor.ID).
		Updates(map[string]interface{}{
			"name":      vendor.Name,
			"code":      vendor.Code,
			"tax_id":    vendor.TaxID,
			"is_active": vendor.IsActive,
		}).Error; err != nil {
		return nil, err
	}
	return vendor, nil
}

// Delete performs a logical delete by setting `is_active` to false.
func (r *repository) Delete(id uuid.UUID) error {
	return r.db.Model(&entities.Vendor{}).
		Where("id = ?", id).
		Update("is_active", false).Error
}
