package vendor_material

import (
	"github.com/google/uuid"
	"github.com/unknowncode44/unknown-go-server/pkg/entities"
	"gorm.io/gorm"
)

// Repository defines persistence operations for VendorMaterial entities.
type Repository interface {
	Create(vm *entities.VendorMaterial) (*entities.VendorMaterial, error)
	FindAll() ([]entities.VendorMaterial, error)
	FindByID(id uuid.UUID) (*entities.VendorMaterial, error)
	Update(vm *entities.VendorMaterial) (*entities.VendorMaterial, error)
	Delete(id uuid.UUID) error
}

type repository struct {
	db *gorm.DB
}

// NewRepo returns a new GORM-backed Repository.
func NewRepo(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(vm *entities.VendorMaterial) (*entities.VendorMaterial, error) {
	if err := r.db.Create(vm).Error; err != nil {
		return nil, err
	}
	return vm, nil
}

func (r *repository) FindAll() ([]entities.VendorMaterial, error) {
	var vms []entities.VendorMaterial
	if err := r.db.Find(&vms).Error; err != nil {
		return nil, err
	}
	return vms, nil
}

func (r *repository) FindByID(id uuid.UUID) (*entities.VendorMaterial, error) {
	var vm entities.VendorMaterial
	if err := r.db.First(&vm, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &vm, nil
}

func (r *repository) Update(vm *entities.VendorMaterial) (*entities.VendorMaterial, error) {
	if err := r.db.
		Model(&entities.VendorMaterial{}).
		Where("id = ?", vm.ID).
		Updates(map[string]interface{}{
			"vendor_code": vm.VendorCode,
			"is_active":   vm.IsActive,
		}).Error; err != nil {
		return nil, err
	}
	return vm, nil
}

func (r *repository) Delete(id uuid.UUID) error {
	return r.db.Model(&entities.VendorMaterial{}).
		Where("id = ?", id).
		Update("is_active", false).Error
}
